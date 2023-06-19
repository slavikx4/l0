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

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "*Handler.GetOrder -> "

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

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось распарсить http request", Op: op}
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte(er.ErrorService)); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}

		orderUID := r.Form["orderUID"][0]
		fmt.Println(orderUID)

		order, err := h.service.GetOrder(orderUID)
		if err != nil {
			err = er.AddOp(err, op)
			logger.Logger.Error.Println(err.Error())
			if _, err = w.Write([]byte("Заказа с таким UID не существует")); err != nil {
				err = &er.Error{Err: err, Code: er.ErrorHTTP, Message: "не удалось записать response", Op: op}
				logger.Logger.Error.Println(err.Error())
			}
			return
		}

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
