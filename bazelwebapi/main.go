package main

import (
	"log"
	"net/http"
)

func main() {
	c, err := loadconfig()
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		statusMap: make(map[string]*status),
		c:         c,
	}

	log.Println("WebAPI running on port 8088")
	log.Fatal(http.ListenAndServe(":8088", s.routes()))
}
