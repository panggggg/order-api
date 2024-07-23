package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/panggggg/order-service/config"
	"github.com/panggggg/order-service/pkg/entity"
)

type OrderAPI interface {
	SaveOrder(order entity.Order) error
}

type orderAPI struct {
	config config.Config
}

func NewOrderAPI(config config.Config) *orderAPI {
	return &orderAPI{config: config}
}

func (o orderAPI) SaveOrder(order entity.Order) error {
	url := fmt.Sprintf("%s/order/%s", o.config.OrderApiURL, order.OrderId)

	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
