package router

import (
	"context"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vave-tool/internal/api/handler"
	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const ClientIDKey contextKey = "client_id"

type Router struct {
	productHandler *handler.ProductHandler
	middleware     *observability.Middleware
	metricsHandler http.Handler
	clientRepo     domain.ClientRepository
}

func NewRouter(productHandler *handler.ProductHandler, middleware *observability.Middleware, metricsHandler http.Handler, clientRepo domain.ClientRepository) *Router {
	return &Router{
		productHandler: productHandler,
		middleware:     middleware,
		metricsHandler: metricsHandler,
		clientRepo:     clientRepo,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/products", r.productHandler.List)
	mux.HandleFunc("GET /api/products/{id}", r.productHandler.GetByID)
	mux.HandleFunc("POST /internal/products", r.requireBasicAuth(r.productHandler.Create))
	mux.HandleFunc("PUT /internal/products/{id}", r.requireBasicAuth(r.productHandler.Update))
	mux.HandleFunc("DELETE /internal/products/{id}", r.requireBasicAuth(r.productHandler.Delete))

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/metrics", r.metricsHandler)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return r.middleware.Handler(r.enableCORS(mux))
}

func (r *Router) requireBasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Internal API"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		client, err := r.clientRepo.GetByUsername(req.Context(), username)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Internal API"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(password)); err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Internal API"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ClientIDKey, client.ID)
		next(w, req.WithContext(ctx))
	}
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
