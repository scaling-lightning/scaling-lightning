package apierrors

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendBadRequestFromErr(w http.ResponseWriter, err error, msg string) {
	log.Error().Err(err).Msg(msg)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("%v: %v", msg, err.Error())))
}

func SendServerErrorFromErr(w http.ResponseWriter, err error, msg string) {
	log.Error().Err(err).Msg(msg)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("%v: %v", msg, err.Error())))
}
