package controllers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/rs/zerolog"
)

type CriptoController struct {
	log zerolog.Logger
}

type SingHmacsha512Request struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type SingHmacsha512Response struct {
	EncodedHex string
}

func NewCriptoController(log zerolog.Logger) *CriptoController {
	return &CriptoController{log: log}

}

func (e CriptoController) SingHmacsha512(w http.ResponseWriter, r *http.Request) {
	var req SingHmacsha512Request
	if err := DecodeJSONBody(w, r, &req); err != nil {
		e.log.Err(err).Msg("SingHmacsha512 error ")
		JSON(w, STATUS_ERROR, err.Error())
		return
	}

	//
	h := hmac.New(sha512.New, []byte(req.Key))
	h.Write([]byte(req.Text))
	sha := hex.EncodeToString(h.Sum(nil))
	//

	JSONstruct(w, SingHmacsha512Response{EncodedHex: sha})
}
