package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/ikawaha/kagome"
	"github.com/tarm/serial"
	"io/ioutil"
	"log"
)

var (
	deviceName = flag.String("d", "/dev/ttyUSB0", "serial port device(default /dev/ttyUSB0)")
	baudRate   = flag.Int("b", 115200, "serial port speed(default 115200)")
	conf       = flag.String("c", "./setting.json", "config file path(default ./setting.json)")
	query      = flag.String("q", "ぼっち OR ﾎﾞｯﾁ OR ボッチ", "query string")
	mode       = flag.String("m", "default", "mode: default|record|emulate")
	recFile    = flag.String("r", "./rec.csv", "tweet records made by mode:record")
)

func init() {
	flag.Parse()
}

func main() {
	c := &serial.Config{Name: *deviceName, Baud: *baudRate}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	onHandler := func(pid PrefId) {
		buf := fmt.Sprintf("0,%d,7\r", pid)
		s.Write([]byte(buf))
		if err != nil {
			log.Fatal(err)
		}
	}
	offHandler := func(pid PrefId) {
		buf := fmt.Sprintf("0,%d,0\r", pid)
		s.Write([]byte(buf))
		if err != nil {
			log.Fatal(err)
		}
	}
	blinkHandler := func(pid PrefId) {
		buf := fmt.Sprintf("1,%d,4\r", pid)
		s.Write([]byte(buf))
		if err != nil {
			log.Fatal(err)
		}
	}
	switch *mode {
	case "default":
		defaultMain(nil, onHandler, offHandler, blinkHandler)
	case "record":
		writer, err := NewTweetWriter(*recFile)
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
		defaultMain(writer, onHandler, offHandler, blinkHandler)
	case "emulate":
		emulateMain(onHandler, offHandler, blinkHandler)
	}
}

func emulateMain(onLed AddCallback, offLed DelCallback, blinkHandler SpecialCallback) {
	led := NewLed(onLed, offLed, blinkHandler)
	reader, err := NewTweetReader(*recFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	for {
		pid, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		led.Add(pid)
		fmt.Printf("led:%+v \n", led)
	}
}

func defaultMain(writer *TweetWriter, onLed AddCallback, offLed DelCallback, blinkHandler SpecialCallback) {
	buf, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatal(err)
	}
	led := NewLed(onLed, offLed, blinkHandler)
	var tc TwitterController
	err = json.Unmarshal(buf, &tc)
	if err != nil {
		log.Fatal(err)
	}
	api := tc.NewApi()
	ss := tc.GetSearchStream(api, *query)
	for {
		select {
		case t := <-ss:
			pid := whereIsTweeted(&t)
			log.Printf("[sample] where:%d desc:%s loc:%s tweet:%s", pid, t.User.Description, t.User.Location, t.Text)
			if pid != PrefInvalid {
				if writer != nil {
					writer.Write(pid)
				} else {
					led.Add(pid)
					fmt.Printf("led:%+v \n", led)
				}
			}
		}
	}
}

func whereIsTweeted(t *anaconda.Tweet) (pid PrefId) {
	pid = PrefInvalid
	tokenizer := kagome.NewTokenizer()
	tokens := tokenizer.Tokenize(t.User.Description + " " + t.User.Location + " " + t.Text)
	for _, m := range tokens {
		if m.Id == kagome.BosEosId {
			continue
		}
		features := m.Features()
		if features[0] != "名詞" || features[1] != "固有名詞" {
			continue
		}
		pid = PrefDict[m.Surface]
		if pid != PrefInvalid {
			return
		}
		//fmt.Printf("%s features:%s\n", m, m.Features())
	}
	return
}
