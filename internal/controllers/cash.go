package controllers

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

type MemoryCashController struct {
	log         zerolog.Logger
	RedisClient RedisClientInterface
}

func NewMemoryCashController(rc RedisClientInterface, log zerolog.Logger) *MemoryCashController {
	return &MemoryCashController{log: log, RedisClient: rc}
}

type RedisClientInterface interface {
	SetParams(ctx context.Context, key string, value uint64) error
	GetParams(ctx context.Context, key string) (uint64, error)
}

type IncrementValueRequest struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}

type IncrementValueResponse struct {
	Value uint64 `json:"value"`
}

func (m *MemoryCashController) IncrementValue(w http.ResponseWriter, r *http.Request) {
	var req IncrementValueRequest
	if err := DecodeJSONBody(w, r, &req); err != nil {
		m.log.Err(err).Msg("IncrementValue")
		JSON(w, STATUS_ERROR, err.Error())
		return
	}
	ctx := r.Context()

	// если нет то инкрeментируем к 0
	value, err := m.RedisClient.GetParams(ctx, req.Key)
	if err != nil {
		m.RedisClient.SetParams(ctx, req.Key, req.Value)
		JSONstruct(w, IncrementValueResponse{Value: req.Value})
	}

	result := value + req.Value
	if err := m.RedisClient.SetParams(ctx, req.Key, result); err != nil {
		m.log.Err(err).Msg("IncrementValue")
		JSON(w, STATUS_ERROR, err.Error())
	}
	JSONstruct(w, IncrementValueResponse{Value: result})
}
