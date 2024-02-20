package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/service"
	"github.com/syntropynet/price-publisher/pkg/cmc"
)

type Config struct {
	PubTxSubjectPrefix string
}

type (
	AllMessageChannel chan map[string]cmc.QuoteInfo
)

type PublishService struct {
	*service.Service
	ctx             context.Context
	cfg             Config
	nats            *nats.Conn
	allMsgTokenChan AllMessageChannel
}

func NewConfig(publisherSubjectPrefix string) Config {
	cfg := Config{
		PubTxSubjectPrefix: publisherSubjectPrefix,
	}
	return cfg
}

func New(conn *nats.Conn, ctx context.Context, cfg Config, allMsgChan AllMessageChannel) *PublishService {
	ret := &PublishService{
		Service:         &service.Service{},
		ctx:             ctx,
		cfg:             cfg,
		allMsgTokenChan: allMsgChan,
	}

	ret.Configure(service.WithContext(ctx), service.WithPrefix("syntropy_defi"), service.WithName("price"), service.WithPubNats(conn))

	return ret
}

func (s *PublishService) Start() context.Context {
	return s.Service.Start()
}

func (s *PublishService) Serve(ctx context.Context) {
	for msg := range s.allMsgTokenChan {
		subject := constructAllSubject(s.cfg.PubTxSubjectPrefix)
		if err := s.Publish(msg, subject); err != nil {
			slog.Error(err.Error())
		}
		for symbol, data := range msg {
			subject := constructSingleSubject(s.cfg.PubTxSubjectPrefix, symbol)
			if err := s.Publish(data, subject); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func constructAllSubject(prefix string) string {
	return fmt.Sprintf("%s.all", prefix)
}

func constructSingleSubject(prefix string, symbol string) string {
	return fmt.Sprintf("%s.single.%s", prefix, symbol)
}
