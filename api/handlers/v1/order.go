package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/abdullohsattorov/mytaxi-APIGateway/api/handlers/models"
	pb "github.com/abdullohsattorov/mytaxi-APIGateway/genproto"
	l "github.com/abdullohsattorov/mytaxi-APIGateway/pkg/logger"
	"github.com/abdullohsattorov/mytaxi-APIGateway/pkg/utils"
)

// CreateOrder ...
// @Summary CreateOrder
// @Description This API for creating a new order
// @Tags order
// @Accept  json
// @Produce  json
// @Param Client request body models.OrderReq true "clientCreateRequest"
// @Success 200 {object} models.OrderRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/orders/ [post]
func (h *handlerV1) CreateOrder(c *gin.Context) {
	var (
		body        pb.OrderReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	body.Status = "created"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MyTaxiService().CreateOrder(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create order", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrder ...
// @Summary GetOrder
// @Description This API for getting order detail
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.OrderRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/order/{id} [get]
func (h *handlerV1) GetOrder(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MyTaxiService().GetOrder(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get order", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListOrders ...
// @Summary ListOrders
// @Description This API for getting list of orders
// @Tags order
// @Accept  json
// @Produce  json
// @Param clientId path string true "ClientID"
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} models.ListOrders
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/orders/{clientId} [get]
func (h *handlerV1) ListOrders(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	clientId := c.Param("clientId")
	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MyTaxiService().ListOrders(
		ctx, &pb.ListOrdersReq{
			ClientId: clientId,
			Limit:    params.Limit,
			Page:     params.Page,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list orders", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrder ...
// @Summary UpdateOrder
// @Description This API for updating order
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param Order request body models.UpdateOrder true "orderUpdateRequest"
// @Success 200 {object} models.OrderRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/orders/{id} [put]
func (h *handlerV1) UpdateOrder(c *gin.Context) {
	var (
		body        pb.OrderReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	if body.Status != "accepted" && body.Status != "cancelled" && body.Status != "finished" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot change standart status variables", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MyTaxiService().UpdateOrder(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update order", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOrder ...
// @Summary DeleteOrder
// @Description This API for deleting order
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/orders/{id} [delete]
func (h *handlerV1) DeleteOrder(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MyTaxiService().DeleteOrder(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete order", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
