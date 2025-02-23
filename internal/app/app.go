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
		log.Println(err)
	}

	c, err = app.LoadConfig("..")
	if err != nil {
		log.Println(err)
	}

	mux := web.NewMux()
	DB := storage.DatabaseInit(c)
	h := handlers.New(DB)

	go mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	go mux.HandleFunc("GET /api/v1/wallets/{WALLET_UUID}", h.HandlerGetBalance)
	go mux.HandleFunc("POST /api/v1/wallet", h.HandlerUpdateBalance)

	http.ListenAndServe(c.APPPort, mux)
	w := fmt.Sprint()
	fmt.Print(w)
}
