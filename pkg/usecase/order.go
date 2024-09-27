package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/panggggg/order-service/pkg/entity"
	"github.com/panggggg/order-service/pkg/repository"
	"github.com/wisesight/spider-go-utilities/cache"
)

type Order interface {
	Upsert(ctx context.Context, orderId string, updateData entity.Order) (bool, error)
	SendToQueue(ctx context.Context, order []string) error

	IsExist(ctx context.Context, order entity.Order) (bool, error)
	Set(ctx context.Context, order entity.Order) error
	Save(order entity.Order) error
}

type order struct {
	orderRepo    repository.Order
	redisAdapter cache.Redis
}

func NewOrder(orderRepo repository.Order, redisAdapter cache.Redis) Order {
	return &order{
		orderRepo,
		redisAdapter,
	}
}

func (r order) Upsert(ctx context.Context, orderId string, updateData entity.Order) (bool, error) {
	// wait group
	var wg sync.WaitGroup
	wg.Add(2)

	var updateError error

	go func() {
		log.Println("Save Status")
		defer wg.Done()
		_, err := r.orderRepo.Set(ctx, updateData)
		if err != nil {
			updateError = err
			return
		}
	}()

	go func() {
		log.Println("Save Order")
		defer wg.Done()
		_, err := r.orderRepo.Upsert(ctx, orderId, updateData)
		if err != nil {
			updateError = err
			return
		}
	}()

	wg.Wait()

	if updateError != nil {
		return false, updateError
	}
	return true, nil
}

func (r order) SendToQueue(ctx context.Context, order []string) error {
	return r.orderRepo.SendToQueue(ctx, order)
}

func (o order) IsExist(ctx context.Context, order entity.Order) (bool, error) {
	orderId := order.OrderId
	_, err := o.redisAdapter.Get(orderId)
	if err != nil {
		return true, err
	}
	// var orderStatus entity.Order
	// err = json.Unmarshal([]byte(*value), &orderStatus)

	// if err != nil {
	// 	return false, err
	// }

	// if value != nil && orderStatus.Status == order.Status {
	// 	return true, nil
	// }
	return false, nil
}

func (o order) Set(ctx context.Context, order entity.Order) error {
	value, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = o.redisAdapter.Set(fmt.Sprintf("%s_%s", order.OrderId, order.Status), string(value), 60*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func (o order) Save(order entity.Order) error {
	return o.orderRepo.SaveWithId(order)
}
