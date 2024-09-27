package entity

import "time"

type Order struct {
	OrderId string `json:"order_id" bson:"order_id"`
	Status  string `json:"status"`
	Remark  string `json:"remark,omitempty"`
}

type OrderStatus struct {
	OrderId   string     `json:"order_id"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Remark    string     `json:"remark"`
}
