package error

import "fmt"

type errorCode string

// обозначение ошибок общего случая
const (
	ErrorNotFound           errorCode = "not found rows"
	ErrorNoConnect          errorCode = "no connection"
	ErrorNoPing             errorCode = "no ping to data base"
	ErrorService            errorCode = "error into service"
	ErrorDataBaseLimitation errorCode = "error into data base limitation"
	ErrorDataBaseIndefinite errorCode = "error into data base indefinite"
	ErrorPublish            errorCode = "error publication"
	ErrorSubscribe          errorCode = "error subscribe"
	ErrorClose              errorCode = "error close"
	ErrorRead               errorCode = "error read"
	ErrorJson               errorCode = "error json"
	ErrorListenAndServe     errorCode = "error listen and serve"
	ErrorHTTP               errorCode = "error with work HTTP"
)

// Error структура собственной ошибки для большей обработки
type Error struct {
	// Вложенная ошибка
	Err error
	// Код ошибки.
	Code errorCode
	// Сообщение об ошибке, которое понятно пользователю.
	Message string
	// Выполняемая операция, путь до ошибки
	Op string
}

// Error функция возвращающая красиво информацию об ошибке
func (e Error) Error() string {
	return fmt.Sprintf("Code: %v; Message: %v; Err: %v; Path: %v", e.Code, e.Message, e.Err, e.Op)
}

// AddOp функция добавления информации о пути
func AddOp(err error, op string) error {
	e := err.(*Error)
	return &Error{
		Err:     e.Err,
		Code:    e.Code,
		Message: e.Message,
		Op:      e.Op + op,
	}
}
