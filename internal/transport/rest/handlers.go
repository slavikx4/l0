package rest

import (
	"encoding/json"
	"fmt"
	"github.com/slavikx4/l0/internal/service"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
	"html/template"
	"net/http"
)

// Handler структура для обработки запросов от web-клиентов
type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetOrder функция для получения и выдачи информации web-клиенту  об order
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "*Handler.GetOrder -> "

	//если клиент ещё не ввёл orderUID
	//выгружаем ему страницу с полем для заполнения orderUID
	if r.Method == http.MethodGet {
		templ, err := template.ParseFiles("./web/index.html")
		if err != nil {
			err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось распарсить index.html файл как шаблон", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}

		if err := templ.Execute(w, nil); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось выполнить шаблон index.html", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}
		//если клиент ввёл orderUID
	} else if r.Method == http.MethodPost {
		//распарсим запрос
		if err := r.ParseForm(); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось распарсить http request", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}

		//вытащим ведённое значение
		orderUID := r.Form["orderUID"][0]
		fmt.Println(orderUID)

		//передаём в сервис для поиска order
		order, err := h.service.GetOrder(orderUID)
		if err != nil {
			//если такого order нет в базе
			err = er.AddOp(err, op)
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte("Заказа с таким UID не существует")); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}
		//если order нашёлся, то приводим order в json, для красивого форматирования
		data, err := json.MarshalIndent(order, "", " ")
		if err != nil {
			err = &er.Error{Err: err, Code: er.ErrorJson, Message: "не удалось закодировать order в json", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}

		//отправляем информацию об order клиенту
		_, err = w.Write(data)
		if err != nil {
			err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать ответ", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}
	}
}
