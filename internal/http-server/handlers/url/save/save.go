package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	res "url-shortender/internal/http-server/api/response"
	"url-shortender/internal/lib/random"
	"url-shortender/internal/lib/sl"
	"url-shortender/internal/lib/validation"
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
		reqLog := log.With(
			slog.String("fn", fn),
			// @todo изучить детальнее как использовать request_id
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		req := new(Request)
		err := render.DecodeJSON(r.Body, req)
		if err != nil {
			reqLog.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, res.Error("failed to decode request"))
			return
		}

		reqLog.Info("request body decoded", slog.Any("data", req))

		if err := validator.New().Struct(req); err != nil {
			//@todo понять как работают ошибки и узнать в каком виде все возвращается
			validateErrs := err.(validator.ValidationErrors)
			errMsg := validation.ChangeErrMsg(validateErrs)
			//@todo прокинуть название инпута
			render.JSON(w, r, res.Error(errMsg))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.MakeRandStr(8)
		}

		id, err := urlSaver.CreateUrl(req.Url, alias)
		if errors.Is(err, storage.ErrURLExists) {
			reqLog.Info("url already exists", slog.String("URL", req.Url))
			render.JSON(w, r, res.Error("url already exists"))
			return
		}

		if err != nil {
			//@todo проследить как ошибка поднимается с нижних уровней
			reqLog.Error("filed to create url", sl.Err(err))
			render.JSON(w, r, res.Error("filed to create url"))
			return
		}
		reqLog.Info("added new URL", slog.Int64("ID", id))

		//@todo изучить как рабоатют ключи в интерфесах. Так как ключ был задан Response неявно
		render.JSON(w, r, Response{
			Response: res.OK(),
			Alias:    alias,
		})
	}

}
