package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Any("/*", handleAll)
	if err := e.Start(":8082"); err != nil {
		panic(err)
	}
}

// Join function
func Join(str ...string) string {
	var s []string
	for _, i := range str {
		s = append(s, i)
	}
	fmt.Println(strings.Join(s, "/"))
	return strings.Join(s, "/")
}
func handleAll(c echo.Context) error {
	req := c.Request()
	url := req.URL
	path := url.Path
	log.Printf("url %v", path)
	str := strings.Split(path, "/")
	fmt.Println("--------------", str[3])
	if str[3] == "auth" {
		// if req.Method == "GET" {
		host := "http://localhost:8081"
		mainURL := Join(host, url.String()[1:])
		newReq, err := http.NewRequest(req.Method, mainURL, req.Body)
		if err != nil {
			return err
		}
		fmt.Println(newReq)
		client := &http.Client{}
		rsp, err := client.Do(newReq)
		if err != nil {
			log.Fatalln(err)
		}
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// log.Println(string(b))
		fmt.Printf("%s\n", b)

	}
	if str[3] == "user" {
		host := "http://localhost:8080"
		mainURL := Join(host, url.String()[1:])
		newReq, err := http.NewRequest(req.Method, mainURL, req.Body)
		if err != nil {
			return err
		}
		fmt.Println(newReq)
		client := &http.Client{}
		rsp, err := client.Do(newReq)
		if err != nil {
			log.Fatalln(err)
		}
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// log.Println(string(b))
		fmt.Printf("%s\n", b)
	}
	return nil
}
