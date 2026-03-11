package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vave-tool/backend/internal/api/handler"
)

type Router struct {
	productHandler *handler.ProductHandler
}

func NewRouter(productHandler *handler.ProductHandler) *Router {
	return &Router{
		productHandler: productHandler,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/products", r.productHandler.List)
	mux.HandleFunc("/api/products/get", r.productHandler.GetByID)

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return r.enableCORS(mux)
}

func (r *Router) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, req)
	})
}
