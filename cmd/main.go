package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/tocy1/toggl/api/handlers"
	"github.com/tocy1/toggl/config"
	"github.com/tocy1/toggl/db"
)

func main() {
	logger := configureLogging()
	config := config.GetConfig(logger)
	db := db.NewMariaDBDataStore(config.MariaDB, logger)
	APIHandler := handlers.NewApiHandler(config.Server, logger, db)

	mux := mux.NewRouter()

	mux.HandleFunc("/api/v1/deck", basicAuthMiddleware(APIHandler.CreateDeck, &config)).Methods("POST")
	mux.HandleFunc("/api/v1/deck", basicAuthMiddleware(APIHandler.OpenDeck, &config)).Methods("GET")
	mux.HandleFunc("/api/v1/deck/draw", basicAuthMiddleware(APIHandler.DrawCards, &config)).Methods("GET")
	srv := configureServer(config, mux)

	logger.Fatal("server-error", srv.ListenAndServe())
}

func basicAuthMiddleware(next http.HandlerFunc, app *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.Server.AdminUsername))
			expectedPasswordHash := sha256.Sum256([]byte(app.Server.AdminPassword))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func configureLogging() lager.Logger {
	logger := lager.NewLogger("toggl")
	sink := lager.NewPrettySink(os.Stdout, config.LogLevel())
	redactingSink, _ := lager.NewRedactingSink(sink, nil, nil)
	logger.RegisterSink(redactingSink)
	return logger
}

func configureServer(config config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         config.Server.BindAddress,
		WriteTimeout: 1500 * time.Second,
		ReadTimeout:  1500 * time.Second,
	}
}
