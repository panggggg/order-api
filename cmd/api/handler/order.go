package handler

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/panggggg/order-service/pkg/entity"
	"github.com/panggggg/order-service/pkg/usecase"
)

type Order interface {
	Upsert(c echo.Context) error
	UploadCsvFile(c echo.Context) error
}

type order struct {
	orderUsecase usecase.Order
}

func NewOrder(orderUsecase usecase.Order) Order {
	return &order{
		orderUsecase: orderUsecase,
	}
}

func (o order) Upsert(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderId := c.Param("id")
	var upsertData entity.Order
	if err := c.Bind(&upsertData); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	_, err := o.orderUsecase.Upsert(ctx, orderId, upsertData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, upsertData)
}

func (o order) UploadCsvFile(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := c.FormFile("order")
	if err != nil {
		log.Panic(err)
	}

	value, err := file.Open()
	if err != nil {
		log.Panic(err)
	}
	defer value.Close()

	csvData, err := io.ReadAll(value)
	if err != nil {
		log.Panic(err)
	}

	data := bytes.NewReader(csvData)
	r := csv.NewReader(data)
	orders, err := r.ReadAll()
	if err != nil {
		log.Panic(err)
	}
	for _, v := range orders[1:] {
		fmt.Println(v)
		err = o.orderUsecase.SendToQueue(ctx, v)
		if err != nil {
			log.Panic(err)
		}
	}

	response := map[string]string{
		"status":  "success",
		"message": "File uploaded successfully",
	}

	return c.JSON(http.StatusOK, response)
}
