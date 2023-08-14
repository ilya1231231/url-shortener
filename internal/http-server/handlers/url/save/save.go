package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	res "url-shortender/internal/http-server/api/response"
	"url-shortender/internal/lib/sl"
)

type Request struct {
	Url string `json:"url" validate:"required,url"`
	//omitempty - неисполненность
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	res.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	CreateUrl(url string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"
		log.With(
			slog.String("fn", fn),
			// @todo изучить детальнее
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		req := new(Request)
		// @todo изучить детальнее chi/render
		err := render.DecodeJSON(r.Body, req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			//@todo изучить render подробнее
			render.JSON(w, r, res.Error("failed to decode request"))
			//верхняя функция не преривает работу хендлера, поэтому прериваем его самостоятельно выходом из функции
			return
		}

		log.Info("request body decoded", slog.Any("data", req))
	}

}
