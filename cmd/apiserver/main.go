package main

import (
	"fmt"
	"github.com/cecobask/vmware-coding-challenge/internal/api/middleware"
	"github.com/cecobask/vmware-coding-challenge/internal/api/pagedata"
	log "github.com/cecobask/vmware-coding-challenge/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger := log.NewLogger()
	pageDataHandler := pagedata.NewHandler(logger)
	chiRouter := chi.NewRouter()
	chiRouter.Use(
		middleware.NewRequestLoggerMiddleware(logger).Handle,
		render.SetContentType(render.ContentTypeJSON),
	)
	chiRouter.Mount("/pagedata", pagedata.NewRouter(pageDataHandler))
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: chiRouter,
	}
	logger.Info("starting http server", zap.String("url", server.Addr))
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}
}
