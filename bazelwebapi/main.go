package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	c, err := loadconfig()
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		activeTasks: newTasks(),
		history:     newTasks(),
		queue:       newQueue(),
		c:           c,
	}

	go executor(s)

	log.Println("WebAPI running on port ", s.c.serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.c.serverPort), s.routes()))
}

func executor(s *server) {
	for {
		t := s.queue.next()
		if t != nil {
			s.runningTask = t.id
			t.state = RUNNING
			t.updated = time.Now()
			err := s.executeBazel(t)
			if err != nil {
				t.err = err.Error()
			}

			t.state = DONE
			now := time.Now()
			t.updated = now
			t.end = now
			s.runningTask = ""

			s.history.add(t)
			s.activeTasks.delete(t.id)
		}
		time.Sleep(time.Second * 1)
	}
}
