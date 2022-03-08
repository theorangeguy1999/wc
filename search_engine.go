package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery" //to perform DOM manipulation
	badger "github.com/dgraph-io/badger/v3"
)


func main(){

	var phrase string

	base_url := "https://en.wikipedia.org"
	working_dir, err := os.Getwd()
	if err != nil {
		return
	}

	// Opening the database
	db_name := "badger_db"
	db,err := badger.Open(badger.DefaultOptions(filepath.Join(working_dir, db_name)))

	if err != nil{
		fmt.Println(err)
		return
	}

	//Taking key word from user
	fmt.Print("Enter key word: ")
	fmt.Scanln(&phrase)

	counter := 0

	// Reading the databse
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		var page []byte
		var k []byte

		//Going through the pages
		for it.Rewind(); it.Valid(); it.Next() {
		  item := it.Item()
		  k = item.Key()
		  err := item.Value(func(v []byte) error {
			page = append([]byte{},v...)
			return nil
		  })
		  if err != nil {
			return err
		  }

		  doc,err := goquery.NewDocumentFromReader(strings.NewReader(string(page)))
		  if err != nil{
			  return err
		  }
		  //Looking in the content part of the body of the webpage
		  text := doc.Find(".mw-parser-output p").Text()
		  if(strings.Contains(strings.ToLower(text),strings.ToLower(phrase))){
			counter++
			//Printing the url to screen
			fmt.Printf("%s%s\n",base_url,string(k))
		  }
		}

		return nil
	  })

	  println("Total ",counter," pages hit")

	defer db.Close()

	return
}