package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-pkgz/routegroup"

	secret "github.com/1704mori/kinko/api/handlers"
)

func New(log *slog.Logger, db *sql.DB) *routegroup.Bundle {
	router := routegroup.New(http.NewServeMux())
	// router.Use(middleware.Cors, middleware.ServeJSON, middleware.Auth)

	v1Router := router.Mount("/api/v1")

	secretHandler := secret.NewHandler(db)
	v1Router.HandleFunc("PUT /secret/{secretName}", secretHandler.AddSecret)
	v1Router.HandleFunc("GET /secret/{secretName}", secretHandler.GetSecret)
	v1Router.HandleFunc(
		"DELETE /{secretName}",
		secretHandler.DeleteSecretKeyAndValue,
	)
	v1Router.HandleFunc("POST /{secretName}", secretHandler.DeleteSecret)

	v1Router.HandleFunc("GET /secrets", secretHandler.GetAllSecrets)

	return router
}
