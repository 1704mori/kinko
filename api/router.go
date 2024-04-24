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
	router.Use(middleware.Cors, middleware.ServeJSON, middleware.Auth)

	v1Router := router.Mount("/api/v1")

	secretRouter := v1Router.Mount("/secret")
	secretHandler := secret.NewHandler(db)
	secretRouter.HandleFunc("PUT /*", secretHandler.AddSecret)
	secretRouter.HandleFunc("GET /*", secretHandler.GetSecret)
	secretRouter.HandleFunc("DELETE /*", secretHandler.DeleteSecretKeyAndValue)
	secretRouter.HandleFunc("POST /*", secretHandler.DeleteSecret)

	v1Router.HandleFunc("GET /secrets", secretHandler.GetAllSecrets)

	return router
}
