package handler

import (
	"github.com/FlyKarlik/auth-service/internal/errs"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/middleware"
	"github.com/FlyKarlik/auth-service/pkg/codes"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

const (
	authenticateType = "Authenticate"
	refreshType      = "Refresh"
)

// @Summary Authentication
// @Description get access and refresh access token
// @Tags Authorization
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} goodResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/auth/{id} [post]
func (h *Handler) Authenticate(c *gin.Context) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), h.tracer, authenticateType)
	defer span.Finish()

	userId, err := middleware.GetUserID(c)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Authenticate] middleware.GetUserID failed: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	clientIp, err := middleware.GetClientIP(c)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Authenticate] middleware.GetClientIP failed: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	response, err := h.usecase.Authentication(ctx, userId, clientIp)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Authenticate] h.usecase.Authentication failed: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	h.log.LogInfo(c, codes.StatusOK, "[handler.Authenticate] successfully done")
	successResponse(c, codes.Success, codes.StatusOK, response)
}

// @Summary Refresh authentication tokens
// @Description Refreshes the access and refresh tokens.
// @Tags Authentication
// @Produce  json
// @Param access header string true "Access token"
// @Param refresh header string true "Refresh token"
// @Success 200 {object} goodResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), h.tracer, refreshType)
	defer span.Finish()

	ctx, err := middleware.BindTokens(c, ctx)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Refresh] middleware.BindTokens error: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	clientIp, err := middleware.GetClientIP(c)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Refresh] middleware.GetClientIP failed: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	response, err := h.usecase.Refresh(ctx, clientIp)
	if err != nil {
		h.log.LogError(c, errs.GetCodeFromError(err), "[handler.Refresh] h.usecase.Authentication failed: %s", err)
		errResponse(c, codes.Failure, errs.GetCodeFromError(err), err)
		return
	}

	h.log.LogInfo(c, codes.StatusOK, "[handler.Refresh] successfully done")
	successResponse(c, codes.Success, codes.StatusOK, response)
}
