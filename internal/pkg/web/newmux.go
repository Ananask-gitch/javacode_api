package web

import (
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
