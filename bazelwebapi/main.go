package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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
	go waitSignal(s)
	err = s.loadActiveTasks()
	if err != nil {
		log.Fatal("loading active tasks failed with error: ", err.Error())
	}
	y, m, d := time.Now().Date()
	path := filepath.Join(s.c.storageConf.Path, fmt.Sprintf("history_%d-%s-%d", y, m.String(), d))
	ts, err := loadTasks(path)
	if err != nil {
		log.Fatal("loading history failed with error: ", err.Error())
	}
	s.history = ts

	// adds active tasks that was maybe interrupted.
	for _, v := range s.activeTasks.Tasks {
		s.queue.add(v)
	}
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
			s.saveActiveTasks()
			y, m, d := time.Now().Date()
			path := filepath.Join(s.c.storageConf.Path, fmt.Sprintf("history_%d-%s-%d", y, m.String(), d))
			appendTask(path, t)
		}
		time.Sleep(time.Second * 1)
	}
}

func waitSignal(s *server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-signals

	s.prepareShutdown()
	os.Exit(0)
}
