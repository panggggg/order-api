package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/panggggg/order-service/pkg/entity"
	repoMock "github.com/panggggg/order-service/pkg/repository/mocks"
	"github.com/panggggg/order-service/pkg/usecase"
	"github.com/stretchr/testify/suite"
	cacheMock "github.com/wisesight/spider-go-utilities/cache/mocks"
)

type OrderUsecaseSuite struct {
	suite.Suite

	// Components
	redisAdapter    *cacheMock.Redis
	orderRepository *repoMock.Order
	orderUsecase    usecase.Order
}

func TestOrderUsecase(t *testing.T) {
	suite.Run(t, new(OrderUsecaseSuite))
}

// Before all
func (s *OrderUsecaseSuite) SetupSuite() {
	s.redisAdapter = &cacheMock.Redis{}
	s.orderRepository = &repoMock.Order{}
	s.orderUsecase = usecase.NewOrder(s.orderRepository, s.redisAdapter)
}

func (o *OrderUsecaseSuite) TestSave() {
	o.Run("When call Save with order, then should called SaveWithId once", func() {
		order := entity.Order{}
		o.orderRepository.On("SaveWithId", order).Return(nil)

		err := o.orderUsecase.Save(order)

		o.Assert().Nil(err)
		o.orderRepository.AssertCalled(o.T(), "SaveWithId", order)
	})
}

func (o *OrderUsecaseSuite) TestSet() {
	o.Run("When call Set with order, then set redis with 1234_shipped", func() {
		ctx := context.Background()
		order := entity.Order{
			OrderId: "1234",
			Status:  "shipped",
		}
		value, _ := json.Marshal(order)
		o.redisAdapter.On("Set", "1234_shipped", string(value), 60*time.Hour).Return(nil)

		err := o.orderUsecase.Set(ctx, order)

		o.Assert().Nil(err)
		o.redisAdapter.AssertCalled(o.T(), "Set", "1234_shipped", string(value), 60*time.Hour)
	})
}

func (o *OrderUsecaseSuite) TestIsExist() {
	o.Run("When order is not exist, then return false", func() {
		ctx := context.Background()
		order := entity.Order{
			OrderId: "1234",
			Status:  "shipped",
		}
		o.redisAdapter.On("Get", "1234_shipped").Return(nil, nil).Once()

		isExist, err := o.orderUsecase.IsExist(ctx, order)

		o.Assert().Nil(err)
		o.Assert().EqualValues(false, isExist)
		o.redisAdapter.AssertExpectations(o.T())
		// o.redisAdapter.AssertCalled(o.T(), "Get", "1234_shipped")
	})

	o.Run("When order is exist, then return true", func() {
		ctx := context.Background()
		order := entity.Order{
			OrderId: "1234",
			Status:  "shipped",
		}
		value, _ := json.Marshal(order)
		strValue := string(value)
		o.redisAdapter.On("Get", "1234_shipped").Return(&strValue, nil).Once()

		isExist, err := o.orderUsecase.IsExist(ctx, order)

		o.Assert().Nil(err)
		o.Assert().EqualValues(true, isExist)
		o.redisAdapter.AssertCalled(o.T(), "Get", "1234_shipped")
	})
}

func (o *OrderUsecaseSuite) TestSendtoQueue() {
	o.Run("When call SendToQueue, then called repo SendToQueue once", func() {
		order := []string{}
		ctx := context.Background()
		o.orderRepository.On("SendToQueue", ctx, order).Return(nil)

		err := o.orderUsecase.SendToQueue(ctx, order)

		o.Assert().Nil(err)
		o.orderRepository.AssertCalled(o.T(), "SendToQueue", ctx, order)
	})
}

func (o *OrderUsecaseSuite) TestUpsert() {
	o.Run("When upsert order success, then return true", func() {
		orderId := "12345"
		updateData := entity.Order{}
		ctx := context.Background()
		o.orderRepository.On("Set", ctx, updateData).Return(nil, nil).Once()
		o.orderRepository.On("Upsert", ctx, orderId, updateData).Return(true, nil).Once()

		upsert, err := o.orderUsecase.Upsert(ctx, orderId, updateData)

		o.Assert().Nil(err)
		o.Assert().EqualValues(true, upsert)
		o.orderRepository.AssertCalled(o.T(), "Set", ctx, updateData)
		o.orderRepository.AssertCalled(o.T(), "Upsert", ctx, orderId, updateData)
	})

	o.Run("When upsert order but cannot set order in redis, then return error", func() {
		orderId := "12345"
		updateData := entity.Order{}
		ctx := context.Background()
		o.orderRepository.On("Set", ctx, updateData).Return(nil, errors.New("Cannot set order status")).Once()
		o.orderRepository.On("Upsert", ctx, orderId, updateData).Return(true, nil)

		upsert, err := o.orderUsecase.Upsert(ctx, orderId, updateData)

		o.Assert().EqualValues(false, upsert)
		o.Assert().EqualError(err, "Cannot set order status")
	})

	// o.Run("When upsert order but cannot update order in mongo, then return error", func() {
	// 	orderId := "12345"
	// 	updateData := entity.Order{}
	// 	id := primitive.NewObjectID()
	// 	ctx := context.Background()
	// 	o.orderRepository.On("Upsert", ctx, orderId, updateData).Return(mock.AnythingOfType("bool"), errors.New("Cannot update order status")).Once()
	// 	o.orderRepository.On("Set", ctx, updateData).Return(&id, nil).Once()

	// 	upsert, err := o.orderUsecase.Upsert(ctx, orderId, updateData)
	// 	fmt.Println("ERR ==> ", err)

	// 	o.Assert().EqualValues(false, upsert)
	// 	o.Assert().EqualError(err, "Cannot update order status")
	// })
}
