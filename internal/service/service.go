package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"github.com/synternet/price-publisher/pkg/cmc"
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

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			slog.Error("Failed to serialize message:", err)
			continue
		}

		if err := s.PublishBuf(msgBytes, "all"); err != nil {
			slog.Error(err.Error())
		}
		for symbol, data := range msg {
			dataBytes, err := json.Marshal(data)
			if err != nil {
				slog.Error("Failed to serialize message:", err)
				continue
			}
			if err := s.PublishBuf(dataBytes, fmt.Sprintf("single.%s", symbol)); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
