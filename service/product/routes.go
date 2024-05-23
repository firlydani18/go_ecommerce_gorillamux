package product

import (
	"fmt"
	"go-ecommerce/model"
	"go-ecommerce/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store model.ProductStore
}

func NewHandler(store model.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.getProductHandler).Methods(http.MethodGet)
	router.HandleFunc("/products", h.createProductHandler).Methods(http.MethodPost)
}

func (h *Handler) getProductHandler(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.AddProduct(model.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
