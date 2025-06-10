package handler

import (
	"fmt"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Auth handler response")
}