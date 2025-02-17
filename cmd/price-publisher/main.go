package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/price-publisher/internal/config"
	"github.com/synternet/price-publisher/internal/service"
	"github.com/synternet/price-publisher/pkg/cmc"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	slog.Info("Config", "PublishIntervalSec", cfg.PublishIntervalSec)

	nkey, jwt, err := CreateUser(cfg.NatsConfig.NKey)
	if err != nil {
		panic(fmt.Errorf("failed to create JWT: %w", err))
	}

	conn, err := options.MakeNats("Price Publisher", cfg.NatsConfig.Urls, "", *nkey, *jwt, "", "", "")
	if err != nil {
		panic(fmt.Errorf("failed to connect to NATS %s: %w", cfg.NatsConfig.Urls, err))
	}

	setErrorHandlers(conn)

	slog.Info("Connected to NATS", "URLS", cfg.NatsConfig.Urls)

	allMsgChan := make(service.AllMessageChannel, 100)
	singleMsgChan := make(service.SingleMessageChannel, 100)

	sPub := service.New(conn, context.Background(), cfg.NatsConfig.PrefixName, cfg.NatsConfig.PublisherName, allMsgChan, singleMsgChan)
	sPub.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		slog.Info("Service interrupted. Exiting...")
		cancel()
	}()

	go sPub.Serve(ctx)

	ticker := time.NewTicker(time.Duration(cfg.PublishIntervalSec) * time.Second)
	defer ticker.Stop()

	cmcConfig := config.CmcConfig{
		Ids:       cfg.CmcConfig.Ids,
		IdsSingle: cfg.CmcConfig.IdsSingle,
		ApiKey:    cfg.CmcConfig.ApiKey,
	}

	allowedIdsSingle := make(map[string]struct{})
	for _, id := range strings.Split(cmcConfig.IdsSingle, ",") {
		allowedIdsSingle[strings.TrimSpace(id)] = struct{}{}
	}

	defer sPub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			quotes, err := cmc.RetrievePrices(cmcConfig)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			symbolQuotes := make(map[string]cmc.QuoteInfo)
			for id, dataItem := range quotes.Data {
				usdQuote, ok := dataItem.Quote["USD"]
				if !ok {
					slog.Info("USD quote not found", "ID", id)
					continue
				}

				symbolQuotes[dataItem.Symbol] = cmc.QuoteInfo{
					Price:           usdQuote.Price,
					PercentChange24: usdQuote.PercentChange24h,
					LastUpdated:     usdQuote.LastUpdated.Unix(),
				}
				if _, exists := allowedIdsSingle[id]; exists {
					singleMsgChan <- service.SingleMessage{
						Symbol: dataItem.Symbol,
						Data:   symbolQuotes[dataItem.Symbol],
					}
				}
			}

			allMsgChan <- symbolQuotes
		}
	}
}

func setErrorHandlers(conn *nats.Conn) {
	if conn == nil {
		return
	}

	conn.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		slog.Error("NATS error", "error", err)
	})
	conn.SetDisconnectHandler(func(c *nats.Conn) {
		slog.Error("NATS disconnected", "error", c.LastError())
	})
	conn.SetReconnectHandler(func(_ *nats.Conn) {
		slog.Info("NATS reconnected")
	})
}
