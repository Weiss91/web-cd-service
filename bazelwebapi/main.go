package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		start:       time.Now(),
	}
	go signalScanner()
	go executor(s)

	log.Println("WebAPI running on port ", s.c.serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.c.serverPort), s.routes()))
}

func executor(s *server) {
	for {
		t := s.queue.next()
		if t != nil {
			s.runningTask = t.Id
			t.setState(RUNNING)
			t.Updated = time.Now()
			err := s.executeBazel(t)
			if err != nil {
				t.Err = err.Error()
			}

			t.setState(DONE)
			now := time.Now()
			t.Updated = now
			t.End = now
			s.runningTask = ""

			s.history.add(t)
			s.activeTasks.delete(t.Id)
		}
		time.Sleep(time.Second * 1)
	}
}

func signalScanner(s *server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	_ = <-signals

	s.prepareShutdown()
	os.Exit(0)
}
