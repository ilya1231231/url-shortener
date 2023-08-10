package sl

import "golang.org/x/exp/slog"

// данная функция нужна для обработки "врапнутых" или "дополненных контекстом" ошибок. slog не умеет такие обрабатывать
// аналог log.Fatalf("failed to init storage %s", err.Error())
// но верхний вариант не используем, так как зачем нам 2 логера
// подробнее - Практики->обработка ошибок

// ниже просто возвращаем атрибут с ключом Error
// выглядеть будет примерно так:
// Error(это кллюч)=(далее Value)"storage.sqlite.New: sql: unknown driver \"sqlitffffe3\" (forgotten import?)"
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "Error",
		Value: slog.StringValue(err.Error()),
	}
}
