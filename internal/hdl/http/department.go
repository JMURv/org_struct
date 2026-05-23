package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/JMURv/golang-clean-template/internal/config"
	"github.com/JMURv/golang-clean-template/internal/ctrl"
	"github.com/JMURv/golang-clean-template/internal/dto"
	"github.com/JMURv/golang-clean-template/internal/hdl"
	"github.com/JMURv/golang-clean-template/internal/hdl/http/utils"
	_ "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) registerDepartmentRoutes() {
	// /departments
	h.mux.HandleFunc("/departments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.createDepartment(w, r)

		default:
			http.NotFound(w, r)
		}
	})

	// /departments/*
	h.mux.HandleFunc("/departments/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/departments/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		parts := strings.Split(path, "/")

		id := parts[0]
		ctx := context.WithValue(r.Context(), config.ParamIDKey, id)
		r = r.WithContext(ctx)

		// /departments/{id}
		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				h.getDepartment(w, r)

			case http.MethodPatch:
				h.updateDepartment(w, r)

			case http.MethodDelete:
				h.deleteDepartment(w, r)

			default:
				http.NotFound(w, r)
			}

			return
		}

		// /departments/{id}/employees
		if len(parts) == 2 && parts[1] == "employees" {
			switch r.Method {
			case http.MethodPost:
				h.createEmployee(w, r)

			default:
				http.NotFound(w, r)
			}

			return
		}

		http.NotFound(w, r)
	})
}
