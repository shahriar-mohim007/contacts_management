package httpserver

import (
	"contacts/state"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Serve(s *state.State) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", handleRegisterUser(s))
		r.Post("/users/activate", HandleActivateUser(s))
		r.Post("/token/auth", handleLogin(s))
		r.Post("/token/refresh", handleRefreshToken(s))
	})
	r.Route("/api/v1/contacts", func(r chi.Router) {

		r.Use(AuthMiddleware(s.Cfg.SecretKey))
		r.Get("/", GetAllContactsHandler(s))
		r.Post("/", CreateContactHandler(s))
		r.Get("/{id}", GetContactByIDHandler(s))
		r.Patch("/{id}", PatchContactHandler(s))
		r.Delete("/{id}", DeleteContactHandler(s))

	})

	serverAddress := fmt.Sprintf(":%d", s.Cfg.ApplicationPort)
	fmt.Printf("Server is running at http://localhost%s\n", serverAddress)

	srv := &http.Server{
		Addr:         serverAddress,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Info().Msgf("Shutting down server. Received signal: %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
			return
		}

		var wg sync.WaitGroup
		wg.Wait()

		shutdownError <- nil
	}()

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Server failed to start")
	}

	if err := <-shutdownError; err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	} else {
		log.Info().Msg("Server shutdown gracefully")
	}
}
