// Package fault предоставляет механизм обработки ошибок с поддержкой интернационализации
package fault

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Arg представляет пару ключ-значение для параметризации ошибок
type Arg struct {
	K string // K Ключ аргумента
	V string // V Значение аргумента
}

// Fault представляет структуру ошибки с поддержкой i18n
type Fault struct {
	err  Code              // err Код ошибки
	args map[string]string // args Аргументы для шаблонизации сообщения об ошибке
}

// ToProto преобразует ошибку в формат gRPC с локализованным сообщением
// Принимает опциональные параметры языка локализации
func (i *Fault) ToProto(lang ...string) error {
	l := i18n.NewLocalizer(Localization, lang...)
	message, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    i.Error(),
		TemplateData: i.Data(),
	})

	if err != nil {
		return status.Error(codes.InvalidArgument, i.Error())
	}
	err = status.Error(codes.InvalidArgument, message)
	return err
}

// Error возвращает строковое представление кода ошибки
func (i *Fault) Error() string {
	return string(i.err)
}

// Data возвращает карту аргументов ошибки
func (i *Fault) Data() map[string]string {
	return i.args
}

// Code представляет тип для кодов ошибок
type Code string

// Err создает новый экземпляр Fault с заданными аргументами
func (c Code) Err(args ...*Arg) *Fault {
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.K] = arg.V
	}
	return &Fault{err: c, args: argsMap}
}

// HandleErr обрабатывает ошибку, преобразуя её в формат gRPC если это возможно
// Если ошибка не является Fault, возвращает UnhandledError
func HandleErr(err error, lang ...string) error {
	f := new(Fault)
	if errors.As(err, &f) {
		return f.ToProto(lang...)
	}
	return UnhandledError.Err().ToProto(lang...)
}

// UnhandledError представляет код ошибки для необработанных исключений
const UnhandledError Code = "UnhandledError"
