package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

const SelectorEndpoint string = "http://localhost:5001/source"

func fmJukeboxNews() {
	// for when we are out of term time
	// and we have jukebox on FM only,
	// we can use the FM source selector
	// to also do news here too

	// wait until 15 seconds to
	currentSeconds := time.Now().Second()
	time.Sleep(time.Duration(45-currentSeconds) * time.Second)

	data := make(url.Values)

	// source 2 is the news feed
	log.Println("selecting FM source 2 (autonews)")
	data["source"] = []string{"2"}
	http.PostForm(SelectorEndpoint, data)

	// wait until the end of the news
	time.Sleep(time.Duration(15+120+5) * time.Second)

	// source 1 is jukebox
	log.Println("selecting FM source 1 (jukebox)")
	data["source"] = []string{"1"}
	http.PostForm(SelectorEndpoint, data)
}
