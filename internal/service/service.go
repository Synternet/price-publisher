package service

import (
	"context"
	"fmt"
	"log"

	"github.com/syntropynet/price-publisher/pkg/cmc"

	svcn "github.com/SyntropyNet/pubsub-go/pubsub"
)

type Config struct {
	PubTxSubjectPrefix string
}

type (
	AllMessageChannel chan map[string]cmc.QuoteInfo
)

type PublishService struct {
	ctx             context.Context
	cfg             Config
	nats            *svcn.NatsService
	allMsgTokenChan AllMessageChannel
}

func NewConfig(publisherSubjectPrefix string) Config {
	cfg := Config{
		PubTxSubjectPrefix: publisherSubjectPrefix,
	}
	return cfg
}

func NewPublishService(s *svcn.NatsService, ctx context.Context, cfg Config, allMsgChan AllMessageChannel) *PublishService {
	return &PublishService{
		ctx:             ctx,
		cfg:             cfg,
		nats:            s,
		allMsgTokenChan: allMsgChan,
	}
}

func (s PublishService) Serve(ctx context.Context) {
	for msg := range s.allMsgTokenChan {
		subject := constructAllSubject(s.cfg.PubTxSubjectPrefix)
		if err := s.nats.PublishAsJSON(s.ctx, subject, msg); err != nil {
			log.Printf("ERROR: %s", err.Error())
		}
		for symbol, data := range msg {
			subject := constructSingleSubject(s.cfg.PubTxSubjectPrefix, symbol)
			if err := s.nats.PublishAsJSON(s.ctx, subject, data); err != nil {
				log.Printf("ERROR: %s", err.Error())
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
