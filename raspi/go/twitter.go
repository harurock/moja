package main

import "github.com/ChimeraCoder/anaconda"
import "log"
import "strconv"
import "time"
import "net/url"

type TwitterController struct {
	ApiKey            string
	ApiSecret         string
	AccessToken       string
	AccessTokenSecret string
}

type FilterFun func(anaconda.Tweet) bool

func (tc *TwitterController) NewApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(tc.ApiKey)
	anaconda.SetConsumerSecret(tc.ApiSecret)
	return anaconda.NewTwitterApi(tc.AccessToken, tc.AccessTokenSecret)
}

func (tc *TwitterController) GetSearchStream(api *anaconda.TwitterApi, query string) chan anaconda.Tweet {
	c := make(chan anaconda.Tweet)
	go func(chanNotify chan anaconda.Tweet) {
		var since_id int64
		for {
			v := url.Values{}
			v.Set("local", "ja")
			v.Set("count", "20")
			if since_id > 0 {
				ssince_id := strconv.FormatInt(since_id, 10)
				v.Set("since_id", ssince_id)
			}
			resp, err := api.GetSearch(query, v)
			if err == nil {
				for _, status := range resp.Statuses {
					chanNotify <- status
					if status.Id > since_id {
						since_id = status.Id
					}
				}
			} else {
				log.Printf("An error occured while searching. err:%v", err)
			}
			time.Sleep(time.Second * 30)
		}
	}(c)
	return c
}
