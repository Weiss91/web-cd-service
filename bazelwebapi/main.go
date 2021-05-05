package main

import (
	"fmt"
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

	log.Println("WebAPI running on port ", s.c.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.c.ServerPort), s.routes()))
}
