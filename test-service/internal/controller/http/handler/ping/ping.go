package ping

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"test-service/internal/lib/api/response"
	"test-service/internal/lib/logger/sl"
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
			slog.String("ping", middleware.GetReqID(r.Context())))

		fmt.Println(r)

		username := r.Header.Get("X-Username")

		log.Info("token", username)

		result, err := useCase.Ping(username)
		if err != nil {
			log.Info("Error", sl.Error(err))

			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		log.Info("Pinged", slog.String("url", result))

		render.JSON(w, r, Response{
			Response: response.OK(),
			URL:      result,
		})
	}
}
