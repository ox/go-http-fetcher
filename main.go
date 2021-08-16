package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

var watchDir string
var showHelp bool

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show usage information")
	flag.StringVar(&watchDir, "dir", "op", "Directory to watch for operatons")
	flag.Parse()
}

func main() {

	if _, err := os.Stat(watchDir); os.IsNotExist(err) {
		fmt.Printf("'%s' does not exist\n\n", watchDir)
		flag.Usage()
		os.Exit(1)
	}

	if showHelp {
		flag.Usage()
		return
	}

	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	r := regexp.MustCompile(".req$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				contents, err := ioutil.ReadFile(event.Path)
				if err != nil {
					fmt.Printf("Could not read file: %s\n", err)
					return
				}

				requestArgs := strings.SplitN(strings.TrimSpace(string(contents)), " ", 3)
				requestMethod, requestUrl := requestArgs[0], requestArgs[1]

				var requestBody io.Reader
				if len(requestArgs) == 3 {
					requestBody = strings.NewReader(requestArgs[2])
				}

				req, err := http.NewRequest(requestMethod, requestUrl, requestBody)
				if err != nil {
					fmt.Printf("Could not request %s: %s\n", requestUrl, err)
					return
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Printf("Could not request %s: %s\n", requestUrl, err)
					return
				}

				fmt.Printf("%d %s %s\n", resp.StatusCode, requestMethod, requestUrl)

				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Could not read response body: %s\n", err)
					return
				}

				responseFilePath := strings.ReplaceAll(event.Path, ".req", ".res")
				ioutil.WriteFile(responseFilePath, respBody, 0750)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(watchDir); err != nil {
		log.Fatalln(err)
	}
	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Second * 1); err != nil {
		log.Fatalln(err)
	}
}
