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

	mux := http.NewServeMux()
	// should publish an image. Target required
	mux.HandleFunc("/execute/task", s.ExecuteTask)

	log.Println("WebAPI running on port 8088")
	log.Fatal(http.ListenAndServe(":8088", mux))
}
