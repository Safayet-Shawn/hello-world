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
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	url := req.URL
	path := url.Path
	log.Printf("url %v", path)
	str := strings.Split(path, "/")
	fmt.Println("--------------", str[3])
	if str[3] == "auth" {
		// if req.Method == "GET" {
		host := "http://localhost:8081"
		mainURL := Join(host, url.String()[1:])
		// fmt.Println(mainURL)
		newReq, err := http.NewRequest(req.Method, mainURL, req.Body)
		if err != nil {
			return err
		}
		// fmt.Println(newReq)
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

		return c.String(http.StatusOK, string(b))
	}
	if str[3] == "user" {

		host := "http://localhost:8080"
		mainURL := Join(host, url.String()[1:])
		// fmt.Println(mainURL)
		newReq, err := http.NewRequest(req.Method, mainURL, req.Body)
		if err != nil {
			return err
		}
		// fmt.Println(newReq)
		req1 := c.Request()

		newReq.Header = req1.Header
		newReq.Header.Add("Accept", "application/json")
		newReq.Header.Add("Content-Type", "application/json")
		fmt.Println("header=>", newReq.Header)
		client := &http.Client{}
		rsp, err := client.Do(newReq)
		if err != nil {
			log.Fatalln(err)
		}
		c.Response().WriteHeader(http.StatusOK)
		defer rsp.Body.Close()

		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatalln(err)

		}
		return c.String(http.StatusOK, string(b))

	}
	return nil
}
