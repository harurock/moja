package main

import "os"
import "time"
import "strconv"
import "encoding/csv"

type TweetWriter struct {
	file     *os.File
	writer   *csv.Writer
	prevTime time.Time
}

type TweetReader struct {
	file   *os.File
	reader *csv.Reader
}

func NewTweetWriter(name string) (writer *TweetWriter, err error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	return &TweetWriter{file, csv.NewWriter(file), time.Now()}, nil
}

func NewTweetReader(name string) (reader *TweetReader, err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	return &TweetReader{file, csv.NewReader(file)}, nil
}

func (w *TweetWriter) Write(pid PrefId) (err error) {
	record := make([]string, 0)
	curTime := time.Now()
	sleep := curTime.Sub(w.prevTime)
	ssleep := strconv.FormatInt(int64(sleep), 10)
	spid := strconv.Itoa(int(pid))
	record = append(record, ssleep, spid)
	err = w.writer.Write(record)
	w.writer.Flush()
	w.prevTime = curTime
	return
}

func (r *TweetReader) Read() (pid PrefId, err error) {
	record, err := r.reader.Read()
	if err != nil {
		return
	}
	sleep, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return
	}
	time.Sleep(time.Duration(sleep))
	i, err := strconv.ParseInt(record[1], 10, 32)
	if err != nil {
		return
	}
	pid = PrefId(i)
	return
}

func (w *TweetWriter) Close() (err error) {
	return w.file.Close()
}

func (r *TweetReader) Close() (err error) {
	return r.file.Close()
}
