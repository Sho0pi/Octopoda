package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Tunnel struct {
	PublicUrl string		`json:"public_url"`
	Name string `json:"name"`
}

type API struct {
	Tunnels []Tunnel		 `json:"tunnels"`
}

func getURL() (string, error){

	resp, err := http.Get("http://127.0.0.1:4040/api/tunnels")
	if err != nil {
		sleep()
		return getURL()
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err!= nil {
		return "", err
	}

	var api API
	if err := json.Unmarshal([]byte(buff), &api); err != nil {
		fmt.Println(err)
		return "", err
	}

	if len(api.Tunnels) == 0 {
		sleep()
		return getURL()
	}
	// CHECK IF THE HTTP SERVER WAS CREATED
	for _, tunnel := range api.Tunnels {
		if strings.Contains(tunnel.PublicUrl, "https"){
			printURL(tunnel.PublicUrl)
			return tunnel.PublicUrl, nil
		}
	}
	sleep()
	return getURL()
}

func sleep(){
	time.Sleep(500 * time.Millisecond)
}

func printURL(url string) {
	fmt.Printf("%s***%s We are online! %s***%s\n", string(colorYellow), string(colorWhite), string(colorYellow), string(colorReset))
	fmt.Printf("%s[%s*%s]%s Your link is: %s %s\n\n", colorWhite, colorYellow, colorWhite, colorYellow, colorReset, url)
}
