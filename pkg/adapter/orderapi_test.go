package adapter_test

import (
	"testing"

	"github.com/panggggg/order-service/pkg/adapter"
	"github.com/panggggg/order-service/pkg/entity"
	"github.com/stretchr/testify/suite"
)

// setup
type OrderAdapterSuite struct {
	suite.Suite

	orderAdapter adapter.OrderAPI
}

func TestOrderAdapterSuite(t *testing.T) {
	suite.Run(t, new(OrderAdapterSuite))
}

func (o *OrderAdapterSuite) TestSaveOrder() {
	order := entity.Order{}
	// url := "localhost:1234/order/123456"

	o.orderAdapter.SaveOrder(order)
}
