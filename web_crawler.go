package main

import (
	"fmt" // for printing to standard output
	"time" // to measure time
	"github.com/gocolly/colly" // web crawler package for go
	"github.com/gocolly/colly/queue"
)

func main(){
	
	start := time.Now() //Start time
	link_count := 0 // Counter for links visited
	count_checker := 1 // Checker for 1 , 10 ,100 .... pages

	// Initializing Colly Collector Object
	c := colly.NewCollector(
		colly.AllowedDomains("www.en.wikipedia.org","en.wikipedia.org"), //Only english wikipedia pages will be crawled through
	)

	// Creating a queue for links to go through
	q, _ := queue.New(
		16, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 100000}, // Queue size
	)

	// Parsing this element of HTML (content) for links
	c.OnHTML(".mw-parser-output p a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		if link_count >= 1000001{
			return;
		} else {
			q.AddURL(e.Request.AbsoluteURL(link)) //Add link to queue
		}
	})

	c.OnResponse(func(r *colly.Response) {

		link_count = link_count + 1 // a new link is visited

		//Print time stats
		if link_count == count_checker{
			fmt.Println("Crawled through ",count_checker," page(s) in time: ",time.Now().Sub(start))
			count_checker = count_checker * 10
		}
	})

	q.AddURL("https://en.wikipedia.org/wiki/Web_crawler") //Seed Url
	q.Run(c) //Start Crawling

	return
}