package register

import (
	responseBase "auth-service/internal/controllers/http/response"
	"auth-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type response struct {
	responseBase.Response
}

type request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UseCase interface {
	Register(username, password string) (int, error)
}

func New(log *slog.Logger, useCase UseCase) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers/user/register.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var request request

		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			log.Error("Failed to decode request body")

			render.JSON(w, r, responseBase.Error("Failed to decode request body"))

			return
		}

		_, err = useCase.Register(request.Username, request.Password)
		if err != nil {
			log.Info("Username already taken", sl.Error(err))

			render.JSON(w, r, responseBase.Error(err.Error()))

			return
		}

		render.JSON(w, r, response{
			Response: responseBase.OK(),
		})
	}
}
