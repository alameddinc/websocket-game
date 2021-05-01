package controller

import (
	"errors"
	"github.com/alameddinc/websocket-game/websocketLib"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func ConnectionTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is already")
	responseSchema := struct {
		Name string `json:"name"`
	}{
		Name: "Alameddin Çelik",
	}
	responseSchemaJson(w, responseSchema)
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		responseErrorJson(w, errors.New("UserId can not be null"))
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := websocketLib.NewClient(conn, WsServer, userId)

	go client.WritePump()
	go client.ReadPump()

	WsServer.RegisterNewUser(client)
}
