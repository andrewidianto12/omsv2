package main

import (
	"errors"
	"net/http"

	"github.com/andrewidianto/common"
	pb "github.com/andrewidianto/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	//gateway
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customer/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJson(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.
		CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJson(w, http.StatusOK, o)
}
func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}
	for _, i := range items {
		if i.ID == "" {
			return errors.New("Item ID is Required!")
		}

		if i.Quantity <= 0 {
			return errors.New("Item must have a valid quantity")
		}
	}

	return nil
}
