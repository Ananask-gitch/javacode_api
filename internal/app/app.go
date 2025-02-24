package app

import (
	"fmt"
	"log"
	"net/http"

	app "testapp/internal/configs"
	"testapp/internal/pkg/handlers"
	"testapp/internal/pkg/storage"
	"testapp/internal/pkg/web"
)

func Run() {

	c, err := app.LoadConfig(".")
	if err != nil {
		log.Fatalln(err)
	}

	mux := web.NewMux()
	DB := storage.DatabaseInit(c)
	h := handlers.New(DB)

	go mux.HandleFunc("GET /", h.Hello)
	go mux.HandleFunc("GET /api/v1/wallets/{WALLET_UUID}", h.GetBalance)
	go mux.HandleFunc("POST /api/v1/wallet", h.UpdateBalance)

	http.ListenAndServe(c.APPPort, mux)
	w := fmt.Sprint()
	fmt.Print(w)
}
