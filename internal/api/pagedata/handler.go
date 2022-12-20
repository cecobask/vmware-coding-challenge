package pagedata

import (
	"errors"
	"github.com/cecobask/vmware-coding-challenge/pkg/entity"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	controller Controller
	logger     *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		controller: NewController(logger),
		logger:     logger,
	}
}

func (h *Handler) GetPageData(w http.ResponseWriter, r *http.Request) {
	req, err := entity.NewGetPageDataRequest(r)
	if err != nil {
		h.logger.Error("error validating http request", zap.Error(err))
		processError(w, r, err)
		return
	}
	pageDataResponse, err := h.controller.GetPageData(req)
	if err != nil {
		h.logger.Error("error fetching page data", zap.Error(err))
		processError(w, r, err)
		return
	}
	_ = pageDataResponse.Render(w, r)
}

func processError(w http.ResponseWriter, r *http.Request, err error) {
	var responseError *entity.ResponseError
	if errors.As(err, &responseError) {
		_ = responseError.Render(w, r)
		return
	}
	render.Status(r, http.StatusInternalServerError)
}
