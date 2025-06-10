package handler

import (
	"fmt"
	"net/http"
)

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Booking handler response")
}