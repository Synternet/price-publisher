package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/syntropynet/price-publisher/internal/config"
	"github.com/syntropynet/price-publisher/internal/service"
	"github.com/syntropynet/price-publisher/pkg/cmc"

	svcnats "github.com/SyntropyNet/pubsub-go/pubsub"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	logStdout(func() {
		log.Println("Data collection:")
		log.Printf("- Interval: %d", cfg.PublishIntervalSec)
	})

	opts := []nats.Option{}

	flagUserCredsJWT, err := svcnats.CreateAppJwt(cfg.NatsConfig.NKey)
	if err != nil {
		log.Fatalf("failed to create sub JWT: %v", err)
	}

	opts = append(opts, nats.UserJWTAndSeed(flagUserCredsJWT, cfg.NatsConfig.NKey))

	svcnPub := svcnats.MustConnect(
		svcnats.Config{
			URI:  cfg.NatsConfig.Urls,
			Opts: opts,
		})
	logStdout(func() {
		log.Println("NATS server connected.")
	})

	allMsgChan := make(service.AllMessageChannel, 1024)

	sPub := service.NewPublishService(svcnPub, context.Background(), service.NewConfig(cfg.NatsConfig.Subject), allMsgChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Service interrupted. Exiting...")
		cancel()
	}()

	go sPub.Serve(ctx)

	ticker := time.NewTicker(time.Duration(cfg.PublishIntervalSec) * time.Second)
	defer ticker.Stop()

	cmcConfig := config.CmcConfig{
		Ids:    cfg.CmcConfig.Ids,
		ApiKey: cfg.CmcConfig.ApiKey,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			quotes, err := cmc.RetrievePrices(cmcConfig)
			if err != nil {
				log.Println(err)
				continue
			}

			symbolQuotes := make(map[string]cmc.QuoteInfo)
			for id, dataItem := range quotes.Data {
				usdQuote, ok := dataItem.Quote["USD"]
				if !ok {
					log.Printf("USD quote for ID %s not found\n", id)
					continue
				}

				symbolQuotes[dataItem.Symbol] = cmc.QuoteInfo{
					Price:           usdQuote.Price,
					PercentChange24: usdQuote.PercentChange24h,
					LastUpdated:     usdQuote.LastUpdated.Unix(),
				}
			}

			allMsgChan <- symbolQuotes
		}
	}
}

func logStdout(myFn func()) {
	originalOutput := log.Writer()
	log.SetOutput(os.Stdout)
	myFn()
	log.SetOutput(originalOutput)
}
