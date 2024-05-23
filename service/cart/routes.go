package cart

import (
	"fmt"
	"go-ecommerce/app/middlewares"
	"go-ecommerce/model"
	"go-ecommerce/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	orderStore   model.OrderStore
	productStore model.ProductStore
	userStore    model.UserStore
}

func NewHandler(orderStore model.OrderStore, productStore model.ProductStore, userStore model.UserStore) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", middlewares.WithJWTAuth(h.checkoutHandler, h.userStore)).Methods("POST")
}

func (h *Handler) checkoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := middlewares.GetUserIdFromCtx(r.Context())

	var cart model.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// get products
	productIds, err := getCartItemsIds(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIds(productIds)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderId,
	})
}
