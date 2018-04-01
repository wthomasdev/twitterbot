package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func main() {
	log := &logger{logrus.New()}
	hashtag := flag.String("hashtag", "#golang", "a hashtag")
	flag.Parse()

	log.Infof("Will search for %s", *hashtag)

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	api.SetLogger(log)
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{*hashtag},
	})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Errorf("received unexpected value of type %T", v)
			continue
		}
		fmt.Printf("%s\n", t.Text)
	}

}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
