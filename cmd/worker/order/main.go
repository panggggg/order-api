package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/panggggg/order-service/config"
	"github.com/panggggg/order-service/pkg/adapter"
	"github.com/panggggg/order-service/pkg/entity"
	"github.com/panggggg/order-service/pkg/repository"
	"github.com/panggggg/order-service/pkg/usecase"
	"github.com/wisesight/spider-go-utilities/cache"
	"github.com/wisesight/spider-go-utilities/queue"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queueNames := []string{cfg.OrderDLQ}
	errorQueueNames := []string{cfg.OrderErrorQueueName}
	rabbitmqAdapter, err := queue.NewRabbitMQ(cfg.RabbitMQURI, queue.QueueConfig{
		QueueNames:      queueNames,
		ErrorQueueNames: errorQueueNames,
		ExchangeName:    cfg.OrderDLX,
		ExchangeType:    cfg.OrderExchangeType,
		RoutingKey:      cfg.OrderQueueName,
		DeadLetter: []map[string]interface{}{
			{"x-message-ttl": 10000},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	defer rabbitmqAdapter.CleanUp()

	redisAdapter := cache.NewRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisPass, 0, ctx)

	orderApi := adapter.NewOrderAPI(cfg)

	orderRepo := repository.NewOrder(nil, nil, nil, nil, cfg, orderApi)
	orderUsecase := usecase.NewOrder(orderRepo, redisAdapter)

	jobs, err := rabbitmqAdapter.Consume(cfg.OrderQueueName, 1)
	log.Println("Start consume message...")
	if err != nil {
		log.Panic(err)
	}
	for job := range jobs {
		var order entity.Order
		err := json.Unmarshal(job.Body, &order)
		if err != nil {
			fmt.Println(err)
			continue
		}

		isExisted, err := orderUsecase.IsExist(ctx, order)
		if isExisted {
			// dead letter
			log.Println("This order is exist")
			job.Reject(false)
			continue
		}
		if err != nil {
			log.Panic(err)
		}

		err = orderUsecase.Set(ctx, order)
		if err != nil {
			log.Panic(err)
		}

		err = orderUsecase.Save(order)
		if err != nil {
			log.Panic(err)
		}
		job.Ack(false)
	}

}
