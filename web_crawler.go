package main

import (
	"fmt" // for printing to standard output
	"time" // to measure time
	"github.com/gocolly/colly" // web crawler package for go
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/colly/extensions"
	badger "github.com/dgraph-io/badger/v3" //badgerdb API
)

func main(){
	
	start := time.Now() //Start time
	link_count := 0 // Counter for links visited
	count_checker := 1 // Checker for 1 , 10 ,100 .... pages
	crawling_limit := 1000000 //To set total number of pages crawled

	//opening database
	opt := badger.DefaultOptions("\\dbbadger")
	db,err := badger.Open(opt)

	if err!=nil{
		fmt.Println(err)
		return
	}

	//starting a transcation
	txn := db.NewTransaction(true)

	// Initializing Colly Collector Object
	c := colly.NewCollector(
		colly.AllowedDomains("www.en.wikipedia.org","en.wikipedia.org"), //Only english wikipedia pages will be crawled through
	)

	extensions.RandomUserAgent(c)

	// Creating a queue for links to go through
	q, _ := queue.New(
		20, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Queue size
	)

	// Parsing this element of HTML (content) for links
	c.OnHTML(".mw-parser-output p a[href]", func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))
		queue_size, err := q.Size()
		if err != nil{
			fmt.Println(err)
			return
		}
		if link_count >= crawling_limit || queue_size >= 1000{
			return;
		} else {
			q.AddURL(link) //Add link to queue
		}
	})

	c.OnResponse(func(r *colly.Response) {

		link_count = link_count + 1 // a new link is visited
		//Print time stats
		if link_count == count_checker{
			fmt.Println("Crawled through ",count_checker," page(s) in time: ",time.Now().Sub(start))
			count_checker = count_checker * 10
		}

		//add to transaction
		err := txn.Set([]byte(r.Request.URL.RawPath),r.Body)

		//if transaction size becomes too big then commit and start new transaction
		if err == badger.ErrTxnTooBig{
			_ = txn.Commit()
			// fmt.Println("commited to database at link count: ", link_count)
			txn = db.NewTransaction(true)
			_ = txn.Set([]byte(r.Request.URL.RawPath),r.Body)
		}
	})

	q.AddURL("https://en.wikipedia.org/wiki/Web_crawler") //Seed Url
	start = time.Now()
	q.Run(c) //Start Crawling
	txn.Commit() //Commit to database

	defer db.Close() //Close database
	return
}