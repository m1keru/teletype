package main

import (
	"flag"
	"github.com/m1keru/teletype/internal/config"
	"github.com/m1keru/teletype/internal/speech"
	"github.com/m1keru/teletype/internal/telegram"
	"log"
	"os"
	"sync"
)

func main() {
	configpath := flag.String("config", "config.yaml", "path to config file")
	var cfg config.Config
	if err := cfg.Setup(configpath); err != nil {
		println(err)
	}

	if cfg.Daemon.LogFile != "" {
		logfile, err := os.Open(cfg.Daemon.LogFile)
		if err != nil {
			log.Printf("Unable to read config: %+v", err)
		}
		log.SetOutput(logfile)
	}
	audioPipe := make(chan []byte)
	textPipe := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go telegram.Run(&cfg, &wg, &audioPipe, &textPipe)
	go speech.Run(&cfg, &wg, &audioPipe, &textPipe)
	println("waiting on WaitGroup")
	wg.Wait()

}
