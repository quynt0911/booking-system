package handler

import (
	"fmt"
	"net/http"
)

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Notification handler response")
}