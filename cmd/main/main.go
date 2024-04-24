package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	router "github.com/1704mori/kinko/api"
	"github.com/1704mori/kinko/internal/config"
	slogpretty "github.com/1704mori/kinko/internal/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	authToken := flag.String("auth-token", "", "Authentication token (required)")
	port := flag.Int("port", 8080, "Port number")
	host := flag.String("host", "0.0.0.0", "Host address")
	authUser := flag.String("auth-user", "", "Basic Auth username")
	authPasswd := flag.String("auth-passwd", "", "Basic Auth password")

	flag.Parse()

	if *authToken == "" && os.Getenv("API_TOKEN") == "" {
		panic("flag auth-token is required or set API_TOKEN environment variable")
	}

	if *authToken == "" {
		*authToken = os.Getenv("API_TOKEN")
	}

	if *authUser == "" && os.Getenv("AUTH_USERNAME") == "" || *authPasswd == "" && os.Getenv("AUTH_PASSWD") == "" {
		panic("flag auth-user and auth-passwd are required or set AUTH_USERNAME and AUTH_PASSWD environment variables")
	}

	if *authUser == "" {
		*authUser = os.Getenv("AUTH_USERNAME")
	}

	if *authPasswd == "" {
		*authPasswd = os.Getenv("AUTH_PASSWD")
	}

	config.NewCofig(&config.ConfigParams{
		AuthToken: *authToken,
		Env:       "dev",
		Host:      *host,
		Port:      *port,
	})

	log := setupLog(config.Config.Env)
	log.Info("config", "config", config.Config)
	db, err := sql.Open("sqlite3", "./kinko.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := prepareDB(db); err != nil {
		log.Error("failed to prepare database", "error", err.Error())
		return
	}

	r := router.New(log, db)
	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != *authUser || pass != *authPasswd {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		SvelteKitHandler("/")(w, r)
	})

	srv := &http.Server{
		Addr:    config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("kinko started")

	<-done
	log.Info("stopping kinko")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", "error", err.Error())

		return
	}

	log.Info("kinko stopped")
}

func SvelteKitHandler(path string) http.HandlerFunc {
	filesystem := http.Dir("frontend")

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, path)
		_, err := filesystem.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("file not found at %s, trying %s.html\n", path, path)
			path = fmt.Sprintf("%s.html", path)
			fmt.Printf("new path: %s\n", path)
		}
		r.URL.Path = path
		http.FileServer(filesystem).ServeHTTP(w, r)
	}
}

func setupLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)

		log = slog.New(handler)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func prepareDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS secrets (
		id TEXT PRIMARY KEY,
		secret_name TEXT,
		key TEXT,
		value TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(secret_name, key)
	)`)
	if err != nil {
		return err
	}

	return nil
}
