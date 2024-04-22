package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	secret "github.com/1704mori/kinko/api/handlers"
	"github.com/1704mori/kinko/api/middleware"
	"github.com/go-pkgz/routegroup"
)

func New(log *slog.Logger, db *sql.DB) *routegroup.Bundle {
	router := routegroup.New(http.NewServeMux())
	router.Use(middleware.ServeJSON, middleware.Auth)

	secretRouter := router.Mount("/secret")
	secretHandler := secret.NewHandler(db)
	secretRouter.HandleFunc("PUT /*", secretHandler.AddSecret)
	secretRouter.HandleFunc("GET /*", secretHandler.GetSecret)
	secretRouter.HandleFunc("DELETE /*", secretHandler.DeleteSecretKeyAndValue)
	secretRouter.HandleFunc("POST /*", secretHandler.DeleteSecret)

	return router
}
