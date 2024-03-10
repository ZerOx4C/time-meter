package webapi

import (
	"encoding/json"
	"net/http"
	"time-meter/logic"
)

type WebApi interface {
	http.Handler

	PostedTasks() []logic.Task
	OnHandled(handler HandledHandler)
}

type RequestType int

const (
	PostSchedule RequestType = iota + 1
)

type HandledHandler func(t RequestType)

type webApi struct {
	postedTasks    []logic.Task
	handledHandler HandledHandler
}

func New() WebApi {
	ret := new(webApi)
	return ret
}

func (wa *webApi) PostedTasks() []logic.Task {
	ret := []logic.Task{}
	ret = append(ret, wa.postedTasks...)
	return ret
}

func (wa *webApi) OnHandled(handler HandledHandler) {
	wa.handledHandler = handler
}

func (wa *webApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.URL.Path {
	case "/schedule":
		switch r.Method {
		case http.MethodPost:
			err = wa.handlePostSchedule(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(err.Error())
		return
	}
}

func (wa *webApi) handlePostSchedule(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&wa.postedTasks); err != nil {
		return err
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		return err
	}

	if wa.handledHandler != nil {
		wa.handledHandler(PostSchedule)
	}

	return nil
}
