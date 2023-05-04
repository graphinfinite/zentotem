package controllers

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

type UserController struct {
	UserStorage UserStorager
	log         zerolog.Logger
}

type UserStorager interface {
	CreateUser(ctx context.Context, name string, age uint32) (userId uint64, err error)
}

func NewUserController(us UserStorager, logger zerolog.Logger) UserController {
	return UserController{UserStorage: us, log: logger}
}

type AddUserRequest struct {
	Name string `json:"Name"`
	Age  uint32 `json:"Age"`
}
type AddUserResponce struct {
	Id uint `json:"Id"`
}

func (c *UserController) AddUser(w http.ResponseWriter, r *http.Request) {

	var req AddUserRequest
	if err := DecodeJSONBody(w, r, &req); err != nil {
		c.log.Err(err).Msg("Add User error")
		JSON(w, STATUS_ERROR, err.Error())
		return
	}

	userId, err := c.UserStorage.CreateUser(r.Context(), req.Name, req.Age)
	if err != nil {
		c.log.Err(err).Msg("UserStorage.CreateUser")
		JSON(w, STATUS_ERROR, err.Error())
		return
	}

	JSONstruct(w, AddUserResponce{Id: uint(userId)})

}
