package main

import (
	"io"
	"net/http"
	"time"
)

func main() {
	c := http.Client{Timeout: 300 * time.Millisecond}

	resp, err := c.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
}
