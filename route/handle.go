package route

import (
	. "github.com/alameddinc/websocket-game/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var mainRouter *mux.Router

func Handle() {
	go WsServer.Run()
	mainRouter = mux.NewRouter()

	mainRouter.HandleFunc("/", ConnectionTest)
	mainRouter.HandleFunc("/ws/{userId}", ServeWS)

	log.Fatalln(http.ListenAndServe(":8095", mainRouter))
}
