package pagedata

import (
	"errors"
	"github.com/cecobask/vmware-coding-challenge/pkg/entity"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	controller *controller
	logger     *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		controller: &controller{logger: logger},
		logger:     logger,
	}
}

func (h *Handler) GetPageData(w http.ResponseWriter, r *http.Request) {
	req, err := entity.NewGetPageDataRequest(r)
	if err != nil {
		h.logger.Error("error validating http request", zap.Error(err))
		var responseError *entity.ResponseError
		if errors.As(err, &responseError) {
			if err = render.Render(w, r, responseError); err != nil {
				h.logger.Error("error rendering response error", zap.Error(err))
				return
			}
		}
		render.Status(r, http.StatusInternalServerError)
		return
	}
	pageDataResponse, err := h.controller.GetPageData(req)
	if err != nil {
		h.logger.Error("error fetching page data", zap.Error(err))
		var responseError *entity.ResponseError
		if errors.As(err, &responseError) {
			if err = render.Render(w, r, responseError); err != nil {
				h.logger.Error("error rendering response error", zap.Error(err))
				return
			}
		}
		render.Status(r, http.StatusInternalServerError)
		return
	}
	if err = pageDataResponse.Render(w, r); err != nil {
		h.logger.Error("error rendering page data response", zap.Error(err))
	}
}
