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

type SingleMessage struct {
	Symbol string
	Data   cmc.QuoteInfo
}

type (
	AllMessageChannel    chan map[string]cmc.QuoteInfo
	SingleMessageChannel chan SingleMessage
)

type PublishService struct {
	*service.Service
	ctx                context.Context
	nats               *nats.Conn
	allMsgTokenChan    AllMessageChannel
	singleMsgTokenChan SingleMessageChannel
}

func New(conn *nats.Conn, ctx context.Context, prefixName string, publisherName string, allMsgChan AllMessageChannel, singleMsgChan SingleMessageChannel) *PublishService {
	ret := &PublishService{
		Service:            &service.Service{},
		ctx:                ctx,
		allMsgTokenChan:    allMsgChan,
		singleMsgTokenChan: singleMsgChan,
	}

	ret.Configure(service.WithContext(ctx), service.WithPrefix(prefixName), service.WithName(publisherName), service.WithPubNats(conn))

	return ret
}

func (s *PublishService) Start() context.Context {
	return s.Service.Start()
}

func (s *PublishService) Serve(ctx context.Context) {
	go func() {
		for msg := range s.allMsgTokenChan {
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				slog.Error("Failed to serialize all token message", "error", err)
				continue
			}
			if err := s.PublishBuf(msgBytes, "all"); err != nil {
				slog.Error(err.Error())
			}
		}
	}()
	go func() {
		for msg := range s.singleMsgTokenChan {
			dataBytes, err := json.Marshal(msg.Data)
			if err != nil {
				slog.Error("Failed to serialize single token message:", "error", err)
				continue
			}
			if err := s.PublishBuf(dataBytes, fmt.Sprintf("single.%s", msg.Symbol)); err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	<-ctx.Done()
}
