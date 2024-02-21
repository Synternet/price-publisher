package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/syntropynet/data-layer-sdk/pkg/service"
	"github.com/syntropynet/price-publisher/pkg/cmc"
)

type (
	AllMessageChannel chan map[string]cmc.QuoteInfo
)

type PublishService struct {
	*service.Service
	ctx             context.Context
	nats            *nats.Conn
	allMsgTokenChan AllMessageChannel
}

func New(conn *nats.Conn, ctx context.Context, prefixName string, publisherName string, allMsgChan AllMessageChannel) *PublishService {
	ret := &PublishService{
		Service:         &service.Service{},
		ctx:             ctx,
		allMsgTokenChan: allMsgChan,
	}

	ret.Configure(service.WithContext(ctx), service.WithPrefix(prefixName), service.WithName(publisherName), service.WithPubNats(conn))

	return ret
}

func (s *PublishService) Start() context.Context {
	return s.Service.Start()
}

func (s *PublishService) Serve(ctx context.Context) {
	for msg := range s.allMsgTokenChan {
		if err := s.Publish(msg, "all"); err != nil {
			slog.Error(err.Error())
		}
		for symbol, data := range msg {
			if err := s.Publish(data, fmt.Sprintf("single.%s", symbol)); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
