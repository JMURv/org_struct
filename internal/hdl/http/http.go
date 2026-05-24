package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	// _ "github.com/JMURv/org-struct/api/rest/v1"
	"github.com/JMURv/org-struct/internal/ctrl"
	mid "github.com/JMURv/org-struct/internal/hdl/http/middleware"
	"github.com/JMURv/org-struct/internal/hdl/http/utils"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	mux  *http.ServeMux
	srv  *http.Server
	ctrl ctrl.AppCtrl
}

func New(ctrl ctrl.AppCtrl) *Handler {
	mux := http.NewServeMux()

	hdl := &Handler{mux: mux, ctrl: ctrl}
	hdl.registerDepartmentRoutes()
	hdl.mux.Handle("/swagger/", httpSwagger.WrapHandler)
	hdl.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		utils.SuccessResponse(w, http.StatusOK, "OK")
	})

	var handler http.Handler = mux
	handler = mid.OT(handler)
	handler = mid.Prometheus(handler)
	handler = mid.StripSlashes(handler)
	handler = mid.Logger(zap.L())(handler)
	hdl.mux = http.NewServeMux()
	hdl.mux.Handle("/", handler)
	return hdl
}

func (h *Handler) Start(port int) {
	h.srv = &http.Server{
		Handler:      h.mux,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	zap.L().Info(
		"Starting HTTP server",
		zap.String("addr", h.srv.Addr),
	)

	err := h.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.L().Error("Server error", zap.Error(err))
	}
}

func (h *Handler) Close(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		done <- h.srv.Shutdown(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
