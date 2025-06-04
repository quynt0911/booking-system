package handler

import (
	"fmt"
	"net/http"
)

func ExpertHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Expert handler response")
}