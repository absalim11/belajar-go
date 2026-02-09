package transaction

import (
	"belajar-go/pkg/response"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if len(req.Items) == 0 {
		response.Error(w, http.StatusBadRequest, "Items cannot be empty")
		return
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid product_id")
			return
		}
		if item.Quantity <= 0 {
			response.Error(w, http.StatusBadRequest, "Quantity must be greater than 0")
			return
		}
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, transaction)
}

func (h *Handler) GetDailySalesReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetDailySalesReport()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, report)
}
