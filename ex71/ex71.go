package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SharedObject struct{
	visited map[string]bool
	visLockCh chan int
	fetcher Fetcher
}

func (so *SharedObject)CrawlThread( url string, depth int, threadDone chan int ){
	if depth <= 0 {
		// inform the parent goroutine that you finished processing url
		threadDone <- 1
		return
	}

	// fetch the urls of this url
	body, urls, err := so.fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		// inform the parent goroutine that you finished processing url
		threadDone <- 1
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	// start GO goroutines for each un-visited url fetched
	uCount := 0
	threadDoneP := make(chan int)
	// acquire the lock byb reading from the channel otherwise block
	<-so.visLockCh
	for _, u := range(urls) {
		if !so.visited[u] {
			so.visited[u] = true
			uCount++
			// start a new goroutine to handle the unvisited url
			go so.CrawlThread(u, depth-1, threadDoneP)
		}
	}
	// release the lock by letting others read from the channel
	so.visLockCh <- 1
	// wait for the children goroutines to finish
	for ; uCount>0; uCount-- {
		<-threadDoneP
	}
	// inform the parent goroutine that you finished processing url
	threadDone <- 1
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	so := &SharedObject{
		make(map[string]bool),
		make(chan int, 1),
		fetcher,
	}

	so.visited[url] = true
	so.visLockCh <- 1

	threadDone := make(chan int)
	go so.CrawlThread(url, depth, threadDone)
	<-threadDone
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}
