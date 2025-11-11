package orderproduction

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rajindersingh041/go-auth-sessions/auth"
	"github.com/rajindersingh041/go-auth-sessions/helper"
	"github.com/rajindersingh041/go-auth-sessions/order"
)

type ProductionHandler struct {
	service ProductionService
	orderService order.OrderService

}

// func NewProductionHandler(service productionService, orderService order.OrderService) *ProductionHandler {
// 	return &ProductionHandler{
// 		service: service,
// 		orderService: orderService,
// 	}
// }

func NewProductionHandler(service ProductionService, orderService order.OrderService) *ProductionHandler {
    return &ProductionHandler{
        service:     service,
        orderService: orderService,
    }
}

func (h *ProductionHandler) RegisterRoutes(mux *http.ServeMux, jwtManager auth.JWTManager) {
	mux.Handle("POST /orderproduction",auth.WithJWTAuth(jwtManager,http.HandlerFunc(h.handleCreateProduction())))
}

func (h *ProductionHandler) handleCreateProduction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(auth.UsernameContextKey).(string)
		if !ok || username == "" {
			helper.RespondError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		ctx := r.Context()

		// parse request body
		var req CreateProductionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.RespondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		defer r.Body.Close()


		// get order ID from request body
		orderObj, err := h.orderService.GetOrderByID(ctx, req.OrderID)
		fmt.Print("Order Obj:", orderObj, "\n")
		fmt.Print("Error:", err,"\n")
		// if err != nil  {
		// 	helper.RespondError(w, http.StatusBadRequest, "valid order ID is required")
		// 	return
		// }
		
		
		production, err := h.service.CreateProduction(ctx, req)
		fmt.Print("Prod Obj:", production, "\n")
		fmt.Print("Error:", err,"\n")
		if err != nil {
			helper.RespondError(w, http.StatusInternalServerError, "failed to create production entry")
			return
		}

		helper.RespondJSON(w,http.StatusCreated,map[string]interface{}{
			"message":"production entry created successfully",
			"production_id": production.ProductionID,

		})

	}
	}
