package router

import (
	"context"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vave-tool/internal/api/handler"
	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const ClientIDKey contextKey = "client_id"

type Router struct {
	productHandler     *handler.ProductHandler
	categoryHandler    *handler.CategoryHandler
	subcategoryHandler *handler.SubcategoryHandler
	authHandler        *handler.AuthHandler
	middleware         *observability.Middleware
	metricsHandler     http.Handler
	clientRepo         domain.ClientRepository
	authService        domain.AuthService
}

func NewRouter(
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	subcategoryHandler *handler.SubcategoryHandler,
	authHandler *handler.AuthHandler,
	middleware *observability.Middleware,
	metricsHandler http.Handler,
	clientRepo domain.ClientRepository,
	authService domain.AuthService,
) *Router {
	return &Router{
		productHandler:     productHandler,
		categoryHandler:    categoryHandler,
		subcategoryHandler: subcategoryHandler,
		authHandler:        authHandler,
		middleware:         middleware,
		metricsHandler:     metricsHandler,
		clientRepo:         clientRepo,
		authService:        authService,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Auth routes (public)
	mux.HandleFunc("POST /auth/google", r.authHandler.GoogleLogin)

	// Product routes
	mux.HandleFunc("GET /api/products", r.productHandler.List)
	mux.HandleFunc("GET /api/products/{id}", r.productHandler.GetByID)
	mux.HandleFunc("POST /internal/products", r.requireAuth(r.productHandler.Create))
	mux.HandleFunc("PUT /internal/products/{id}", r.requireAuth(r.productHandler.Update))
	mux.HandleFunc("DELETE /internal/products/{id}", r.requireAuth(r.productHandler.Delete))

	// Category routes
	mux.HandleFunc("GET /api/categories", r.categoryHandler.List)
	mux.HandleFunc("GET /api/categories/{id}", r.categoryHandler.GetByID)
	mux.HandleFunc("POST /internal/categories", r.requireAuth(r.categoryHandler.Create))
	mux.HandleFunc("PUT /internal/categories/{id}", r.requireAuth(r.categoryHandler.Update))
	mux.HandleFunc("DELETE /internal/categories/{id}", r.requireAuth(r.categoryHandler.Delete))

	// Subcategory routes
	mux.HandleFunc("GET /api/subcategories", r.subcategoryHandler.List)
	mux.HandleFunc("GET /api/subcategories/{id}", r.subcategoryHandler.GetByID)
	mux.HandleFunc("GET /api/categories/{category_id}/subcategories", r.subcategoryHandler.GetByCategoryID)
	mux.HandleFunc("POST /internal/subcategories", r.requireAuth(r.subcategoryHandler.Create))
	mux.HandleFunc("PUT /internal/subcategories/{id}", r.requireAuth(r.subcategoryHandler.Update))
	mux.HandleFunc("DELETE /internal/subcategories/{id}", r.requireAuth(r.subcategoryHandler.Delete))

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/metrics", r.metricsHandler)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return r.middleware.Handler(r.enableCORS(mux))
}

func (r *Router) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := r.authService.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(req.Context(), ClientIDKey, claims.UserID)
			next(w, req.WithContext(ctx))
			return
		}

		r.requireBasicAuth(next)(w, req)
	}
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
