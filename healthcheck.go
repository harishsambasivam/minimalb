package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func healthCheck() {
	s := gocron.NewScheduler(time.Local)
	for _, server := range servers {
		s.Every("1s").Do(server.healthCheck)
	}
	s.StartAsync()
}
