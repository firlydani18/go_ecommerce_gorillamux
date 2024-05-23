package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"go-ecommerce/service/cart"
	"go-ecommerce/service/order"
	"go-ecommerce/service/product"
	"go-ecommerce/service/user"
	"log"
	"net/http"
)

type APIServer struct {
	Addr string
	DB   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		Addr: addr,
		DB:   db,
	}
}

func (api *APIServer) Start() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(api.DB)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(api.DB)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(api.DB)
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Server running on", api.Addr)
	return http.ListenAndServe(api.Addr, router)
}
