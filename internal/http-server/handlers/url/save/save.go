package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"math/rand"
	"net/http"
	res "url-shortender/internal/http-server/api/response"
	"url-shortender/internal/lib/sl"
	"url-shortender/internal/storage"
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
			//верхняя функция не прерывает работу хендлера, поэтому прериваем его самостоятельно выходом из функции
			return
		}

		log.Info("request body decoded", slog.Any("data", req))

		if err := validator.New().Struct(req); err != nil {
			//@todo понять как работают ошибки и узнать в каком виде все возвращается
			validateErrs := err.(validator.ValidationErrors)
			//@todo сделать user friendly сообщения об ошибке валидации и протестировать. Пока захардкодил
			log.Info("array of errors", slog.Any("errors", validateErrs))
			render.JSON(w, r, res.Error("Validation Error"))
			return
		}

		alias := req.Alias
		if alias == "" {
			//@todo сделать генератор рандомной строки куда мы передаем кол-во символов строки. Пакет random.
			alias = string(rune(rand.Int()))
		}

		id, err := urlSaver.CreateUrl(req.Url, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("URL", req.Url))
			render.JSON(w, r, res.Error("url already exists"))
			return
		}

		if err != nil {
			//@todo проследить как ошибка поднимается с нижних уровней
			log.Error("filed to create url", sl.Err(err))
			render.JSON(w, r, res.Error("filed to create url"))
			return
		}
		log.Info("added new URL", slog.Int64("ID", id))

		//@todo изучить как рабоатют ключи в интерфесах. Так как ключ был задан Response неявно
		render.JSON(w, r, Response{
			Response: res.OK(),
			Alias:    alias,
		})
	}

}
