package scraper

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Cache interface {
	Add(value string) bool
	Clear()
}

type Words map[string]int
type Results map[string]Words

type FifoQueue interface {
	Enqueue(value string)
	Dequeue() (string, bool)
	IsEmpty() bool
}

func Run(urlsQueue FifoQueue, numberOfScrapers int, cache Cache) Results {
	// we will limit max number of gorutines to `numberOfScrapers`, this will be usefull if we convert scraper to crawler
	scrapersLimit := make(chan struct{}, numberOfScrapers)
	// we will store finall results here
	results := make(Results)
	resultsChan := make(chan Results)
	resultsLock := sync.Mutex{}
	finishedChan := make(chan struct{})

	// collect results
	go func() {
		for {
			select {
			case r := <-resultsChan:
				resultsLock.Lock()
				for k, v := range r {
					results[k] = v
				}
				resultsLock.Unlock()
			case <-finishedChan:
				// close this gorutine when scrapers are done
				return
			}
		}
	}()

	// this solution is ready to became crawler. I'm assuming that scraper goruting could add any time new url to queue
	// this is why I'm not using WaitGroup
	for {
		url, ok := urlsQueue.Dequeue()
		if ok {
			// TODO: need to add url parser url.Parse() from "net/url" here and skip everything after `#`

			// if url was not in cache
			if cache.Add(url) {
				// add limit token
				scrapersLimit <- struct{}{}

				// TODO: this function need to be separated to individual code block
				go func(url string) {
					log.Println("Running for", url)
					defer func() {
						// remove limit token on finish
						<-scrapersLimit
					}()

					// let's assume 2s timeout
					client := http.Client{
						Timeout: 2 * time.Second,
					}
					resp, err := client.Get(url)
					if err != nil {
						log.Println(url, err)
						return
					}
					defer resp.Body.Close()
					// skip http codes different than 200
					if resp.StatusCode == http.StatusOK {
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							log.Println(url, err)
							return
						}

						words := strings.Fields(string(body))
						frequency := make(Words)
						for _, word := range words {
							frequency[strings.ToLower(word)]++
						}

						resultsChan <- Results{url: frequency}
					}
				}(url)
			}
		} else {
			// if there are no active workers and queue is empty
			if len(scrapersLimit) == 0 && urlsQueue.IsEmpty() {
				break
			}
			// go to sleep and next check if there is something in the queue
			time.Sleep(100 * time.Millisecond)
		}
	}
	close(finishedChan)

	// aquire lock to prevent data race on `results`
	resultsLock.Lock()
	resultsLock.Unlock()

	return results
}
