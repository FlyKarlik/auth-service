package handler

import (
	"github.com/FlyKarlik/auth-service/internal/usecase"
	"github.com/FlyKarlik/auth-service/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type Handler struct {
	usecase *usecase.Usecase
	log     *logger.Logger
	tracer  opentracing.Tracer
}

func New(usecase *usecase.Usecase, log *logger.Logger, tracer opentracing.Tracer) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
		tracer:  tracer,
	}
}
