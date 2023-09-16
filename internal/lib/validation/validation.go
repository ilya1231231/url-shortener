package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ChangeErrMsg(errs validator.ValidationErrors) string {
	ruMsgMap := map[string]string{
		"required": "Значение обязательно.",
		"url":      "Значение должно быть ссылкой.",
	}
	sb := &strings.Builder{}
	for _, e := range errs {
		ruMsg := ruMsgMap[e.Tag()]
		//функция Sprintf. Возвращается отформатированную строку
		sb.WriteString(fmt.Sprintf("Ошибка в поле %s. %s", e.Field(), ruMsg))
	}
	return sb.String()
}
