package error

type errorCode string

const (
	ErrorNotFound errorCode = "not found"
	ErrorService  errorCode = "error service"

	Arrow = "->"
)

type Error struct {
	// Вложенная ошибка
	Err error
	// Код ошибки.
	Code errorCode
	// Сообщение об ошибке, которое понятно пользователю.
	Message string
	// Выполняемая операция
	Op string
}
