package route

import (
	"github.com/labstack/echo/v4"
	"github.com/panggggg/order-service/cmd/api/handler"
	"github.com/panggggg/order-service/config"
)

func NewRoute(config config.Config, app *echo.Echo, orderHandler handler.Order) {
	o := app.Group("/order")

	o.PATCH("/:id", orderHandler.Upsert)
	o.POST("/file", orderHandler.UploadCsvFile)
}
