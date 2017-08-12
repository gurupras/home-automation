package web_api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gurupras/go-stoppable-net-listener"
	log "github.com/sirupsen/logrus"
)

type WebAPI struct {
	*mux.Router
	Port int `yaml:"port"`
	snl  *stoppablenetlistener.StoppableNetListener
}

func (w *WebAPI) Initialize() {
	w.Router = mux.NewRouter()
}

func (w *WebAPI) Start() {
	mux := http.NewServeMux()
	mux.Handle("/", w.Router)
	server := http.Server{}
	server.Handler = mux
	snl, err := stoppablenetlistener.New(w.Port)
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.snl = snl
	log.Infof("Starting webserver on port: %v", w.Port)
	server.Serve(snl)
}

func (w *WebAPI) Stop() {
	if w.snl != nil {
		w.snl.Stop()
		w.snl = nil
	}
}
