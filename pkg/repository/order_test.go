package repository_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/panggggg/order-service/config"
	"github.com/panggggg/order-service/pkg/adapter/mocks"
	"github.com/panggggg/order-service/pkg/entity"
	"github.com/panggggg/order-service/pkg/repository"
	"github.com/panggggg/order-service/pkg/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wisesight/spider-go-utilities/database"
	mongodbAdapterMocks "github.com/wisesight/spider-go-utilities/database/mocks"
	rabbitmqAdapterMocks "github.com/wisesight/spider-go-utilities/queue/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepositorySuite struct {
	suite.Suite

	orderAPIAdapter       *mocks.OrderAPI
	mockMongoAdapter      *mongodbAdapterMocks.MongoDB
	mockRabbitMQAdapter   *rabbitmqAdapterMocks.RabbitMQ
	orderCollection       database.MongoCollection
	orderStatusCollection database.MongoCollection
	orderRepository       repository.Order
	config                config.Config
}

func TestOrderRepository(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}

func (s *OrderRepositorySuite) SetupSuite() {
	s.orderAPIAdapter = &mocks.OrderAPI{}
	s.mockMongoAdapter = &mongodbAdapterMocks.MongoDB{}
	s.mockRabbitMQAdapter = &rabbitmqAdapterMocks.RabbitMQ{}
	s.orderCollection = nil
	s.orderStatusCollection = nil
	s.config = config.Config{
		OrderQueueName: "order:job",
	}
	s.orderRepository = repository.NewOrder(s.mockMongoAdapter, s.orderCollection, s.orderStatusCollection, s.mockRabbitMQAdapter, s.config, s.orderAPIAdapter)
}

func (o *OrderRepositorySuite) TestSaveWithId() {
	o.Run("When save order using order api, then order api didn't error", func() {
		order := entity.Order{
			OrderId: "1234",
			Status:  "shipped",
			Remark:  "ส่งแล้ว",
		}
		o.orderAPIAdapter.On("SaveOrder", order).Return(nil)

		err := o.orderRepository.SaveWithId(order)

		o.Assert().Nil(err)
		o.orderAPIAdapter.AssertCalled(o.T(), "SaveOrder", order)
	})
}

func (o *OrderRepositorySuite) TestUpsert() {
	o.Run("When can upsert order status, then return true", func() {
		ctx := context.Background()
		timer := service.NewTime()
		timer.Freeze()
		defer timer.Unfreeze()
		order := entity.Order{}
		query := bson.M{
			"_id": "order_1234",
		}
		update := bson.M{
			"$set": order,
			"$setOnInsert": bson.M{
				"created_at": timer.Now(),
			},
			"$currentDate": bson.M{
				"updated_at": true,
			},
		}
		o.mockMongoAdapter.On("UpdateOne", ctx, o.orderCollection, query, update, mock.Anything).Return(true, nil)

		res, err := o.orderRepository.Upsert(ctx, "1234", order)

		o.Assert().Nil(err)
		o.Assert().EqualValues(true, res)
		o.mockMongoAdapter.AssertCalled(o.T(), "UpdateOne", ctx, o.orderCollection, query, update, mock.Anything)
	})
}

func (o *OrderRepositorySuite) TestSet() {
	o.Run("When save order success, then return order id", func() {
		ctx := context.Background()

		timer := service.NewTime()
		timer.Freeze()
		defer timer.Unfreeze()
		now := timer.Now()

		formatData := map[string]interface{}{
			"order_id":   "",
			"status":     "",
			"remark":     "",
			"created_at": now,
		}
		order := entity.Order{}
		id := primitive.NewObjectID()
		o.mockMongoAdapter.On("InsertOne", ctx, o.orderCollection, formatData).Return(&id, nil)

		_, err := o.orderRepository.Set(ctx, order)

		o.Assert().Nil(err)
		o.mockMongoAdapter.AssertCalled(o.T(), "InsertOne", ctx, o.orderCollection, formatData)

	})
}

func (o *OrderRepositorySuite) TestSendToQueue() {
	o.Run("When called SendToQueue, then should called publish once", func() {
		ctx := context.Background()
		order := []string{"12345", "shipped", "ส่งแล้ว"}
		formatOrder := entity.OrderStatus{
			OrderId: order[0],
			Status:  order[1],
			Remark:  order[2],
		}
		body, _ := json.Marshal(formatOrder)
		o.mockRabbitMQAdapter.On("Publish", ctx, o.config.OrderQueueName, body).Return(nil)

		o.orderRepository.SendToQueue(ctx, order)

		o.mockRabbitMQAdapter.AssertCalled(o.T(), "Publish", ctx, o.config.OrderQueueName, body)
	})
}
