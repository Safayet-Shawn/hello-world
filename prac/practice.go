package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("localhost:8081/api/v1/user/register")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
	//  // create a WaitGroup
	//  wg := new(sync.WaitGroup)

	//  // add two goroutines to `wg` WaitGroup
	//  wg.Add(2)

	//  // create a default route handler
	//  http.HandleFunc( "/", func( res http.ResponseWriter, req *http.Request ) {
	//      fmt.Fprint( res, "Hello: " + req.Host )
	//  } )

	//  // goroutine to launch a server on port 9000
	//  go func() {
	//      log.Fatal( http.ListenAndServe( ":8082", nil ) )
	//      wg.Done() // one goroutine finished
	//  }()
	//  go func() {
	//      log.Fatal( http.ListenAndServe( ":8080", nil ) )
	//      wg.Done() // one goroutine finished
	//  }()

	//  // goroutine to launch a server on port 9001
	// log.Fatal( http.ListenAndServe( ":8081", nil ) )

	//  // wait until WaitGroup is done
	//  wg.Wait()

}
