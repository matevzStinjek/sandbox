package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/internal/config"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/internal/infrastructure/kafka"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/internal/infrastructure/provider"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/internal/infrastructure/redis"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/internal/processor"
	"github.com/matevzStinjek/distributed-trading-system/market-data-ingest/pkg/marketdata"
)

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	// setup context and config
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// setup clients
	cacheClient, err := redis.NewRedisCacheClient(ctx, cfg)
	if err != nil {
		log.Fatalf("error pinging redis cache instance: %v", err)
	}

	pubsubClient, err := redis.NewRedisPubsubClient(ctx, cfg)
	if err != nil {
		log.Fatalf("error pinging redis pubsub instance: %v", err)
	}

	saramaProducer, err := kafka.NewKafkaAsyncProducer(cfg)
	if err != nil {
		log.Fatalf("error creating sarama async producer: %v", err)
	}

	// setup tradeProcessor, channel, and start consuming channel
	tradeProcessor := processor.NewTradeProcessor(cacheClient, pubsubClient, saramaProducer)

	tradeChannel := make(chan marketdata.Trade, cfg.TradeChannelBuff)

	var consumeWg sync.WaitGroup
	consumeWg.Add(1)
	go func() {
		defer consumeWg.Done()
		tradeProcessor.Start(ctx, tradeChannel)
	}()

	// setup stocks client
	client, err := provider.NewMarketClient(ctx)
	if err != nil {
		log.Fatalf("error connecting to stocks client: %v", err)
	}

	client.SubscribeToTrades(func(t stream.Trade) {
		tradeChannel <- marketdata.Trade{
			ID:        t.ID,
			Symbol:    t.Symbol,
			Price:     t.Price,
			Size:      t.Size,
			Timestamp: t.Timestamp,
		}
	}, cfg.Symbols)

	// mock trades
	tradeChannel <- marketdata.Trade{
		ID:        1,
		Symbol:    "MSFT",
		Price:     1,
		Size:      1,
		Timestamp: time.Now(),
	}
	tradeChannel <- marketdata.Trade{
		ID:        2,
		Symbol:    "MSFT",
		Price:     1,
		Size:      1,
		Timestamp: time.Now(),
	}

	<-ctx.Done()
	consumeWg.Done()
}
