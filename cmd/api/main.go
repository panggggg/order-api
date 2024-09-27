package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/panggggg/order-service/cmd/api/handler"
	"github.com/panggggg/order-service/cmd/api/route"
	"github.com/panggggg/order-service/config"
	"github.com/panggggg/order-service/pkg/repository"
	"github.com/panggggg/order-service/pkg/usecase"
	"github.com/wisesight/spider-go-utilities/database"
	"github.com/wisesight/spider-go-utilities/queue"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbClient, err := database.NewMongoDBConnection(ctx, cfg.MongoDBURI)
	if err != nil {
		panic(err)
	}

	mongodbAdapter := database.NewMongoDB(mongodbClient)

	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	orderCollection := mongodbClient.Database(cfg.DatabaseName).Collection(cfg.OrderCollection)
	orderStatusCollection := mongodbClient.Database(cfg.DatabaseName).Collection(cfg.OrderStatusCollection)

	queueNames := []string{cfg.OrderQueueName}
	rabbitmqAdapter, err := queue.NewRabbitMQ(cfg.RabbitMQURI, queue.QueueConfig{
		QueueNames:   queueNames,
		ExchangeName: cfg.OrderExchangeName,
		ExchangeType: cfg.OrderExchangeType,
		DeadLetter: []map[string]interface{}{
			{"x-dead-letter-exchange": cfg.OrderDLX},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	defer rabbitmqAdapter.CleanUp()

	orderRepo := repository.NewOrder(mongodbAdapter, orderCollection, orderStatusCollection, rabbitmqAdapter, cfg, nil)
	orderUsecase := usecase.NewOrder(orderRepo, nil)

	app := echo.New()

	route.NewRoute(cfg, app, handler.NewOrder(orderUsecase))
	err = app.Start(":1234")
	if err != nil {
		go func() {
			if err := app.Start(":1234"); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		defer cancel()
		if err := app.Shutdown(ctx); err != nil {
			app.Logger.Fatal(err)
		}
	}

}
