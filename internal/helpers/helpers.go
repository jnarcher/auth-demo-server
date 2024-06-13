package helpers

import (
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter,  msg string, statusCode int) {
    log.Println("ERROR -", msg)
    http.Error(w, msg, statusCode)
}

