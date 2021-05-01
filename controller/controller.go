package controller

import (
	"encoding/json"
	"github.com/alameddinc/websocket-game/websocketLib"
	"net/http"
)

var WsServer = websocketLib.NewWebsocketServer()

func responseSchemaJson(w http.ResponseWriter, body interface{}) {
	json.NewEncoder(w).Encode(body)
}

func responseErrorJson(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}
