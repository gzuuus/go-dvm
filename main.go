package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sebdeveloper6952/go-dvm/lightning/lnbits"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/sebdeveloper6952/go-dvm/engine"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.TraceLevel)

	lnSvc, err := lnbits.New(
		os.Getenv("LNBITS_API_URL"),
		os.Getenv("LNBITS_KEY"),
	)
	if err != nil {
		log.Fatal(err)
	}

	dvm, err := NewMalwareDvm(
		os.Getenv("DVM_SK"),
	)
	if err != nil {
		log.Fatal(err)
	}

	e, err := engine.NewEngine()
	if err != nil {
		log.Fatal(err)
	}
	e.SetLnService(lnSvc)

	e.RegisterDVM(dvm)

	if err := e.Run(
		ctx,
		[]string{
			"wss://nostr-pub.wellorder.net",
			"wss://nos.lol",
		},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("running...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case sig := <-sigChan:
			if sig == os.Interrupt {
				cancelCtx()
				log.Println("bye")
				return
			}
		}
	}
}
