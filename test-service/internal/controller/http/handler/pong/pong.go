package pong

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"test-service/internal/lib/api/response"
)

type Response struct {
	response.Response
	URL string `json:"url,omitempty"`
}

type UseCase interface {
	Ping(userName string) (string, error)
}

func New(log *slog.Logger, useCase UseCase) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers/ping/ping.New"

		log = log.With(
			slog.String("op", op),
			slog.String("pong", middleware.GetReqID(r.Context())))

		result := "pong"

		log.Info("Pinged", slog.String("url", result))

		render.JSON(w, r, Response{
			Response: response.OK(),
			URL:      result,
		})
	}
}
