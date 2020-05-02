package main

import (
	"flag"
	"github.com/m1keru/teletype/internal/config"
	"github.com/m1keru/teletype/internal/speech"
	"github.com/m1keru/teletype/internal/telegram"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

func main() {
	configpath := flag.String("config", "config.yaml", "path to config file")
	var cfg config.Config
	if err := cfg.Setup(configpath); err != nil {
		log.Fatalf("%+v", err)
	}

	if cfg.Daemon.LogFile != "" {
		logfile, err := os.Open(cfg.Daemon.LogFile)
		if err != nil {
			log.Fatalf("Unable to read config: %+v", err)
		}
		log.SetOutput(logfile)
	}
	if cfg.Daemon.Debug == true {
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(log.InfoLevel)
	audioPipe := make(chan []byte)
	textPipe := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go telegram.Run(&cfg, &wg, &audioPipe, &textPipe)
	go speech.Run(&cfg, &wg, &audioPipe, &textPipe)
	log.Debug("waiting on WaitGroup")
	wg.Wait()

}
