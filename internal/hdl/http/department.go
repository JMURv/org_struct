package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JMURv/golang-clean-template/internal/config"
	"github.com/JMURv/golang-clean-template/internal/ctrl"
	"github.com/JMURv/golang-clean-template/internal/dto"
	"github.com/JMURv/golang-clean-template/internal/hdl"
	"github.com/JMURv/golang-clean-template/internal/hdl/http/utils"
	"github.com/JMURv/golang-clean-template/internal/hdl/validation"
	_ "github.com/JMURv/golang-clean-template/internal/models"
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
		ctx := context.WithValue(r.Context(), config.IDKey, id)
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

func (h *Handler) createDepartment(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateDepartmentRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	res, err := h.ctrl.CreateDepartment(r.Context(), req)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, res)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.Context().Value(config.IDKey).(string)
	if !ok {
		utils.ErrResponse(w, http.StatusBadRequest, ErrRetrievePathVars)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	var req dto.CreateDepartmentEmployeeRequest
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	res, err := h.ctrl.CreateEmployee(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ctrl.ErrNotFound):
			utils.ErrResponse(w, http.StatusNotFound, err)
			return
		default:
			utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
			return
		}
	}

	utils.SuccessResponse(w, http.StatusCreated, res)
}

func (h *Handler) getDepartment(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.Context().Value(config.IDKey).(string)
	if !ok {
		utils.ErrResponse(w, http.StatusBadRequest, ErrRetrievePathVars)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	depth := 1

	if depthStr := r.URL.Query().Get("depth"); depthStr != "" {
		parsedDepth, err := strconv.Atoi(depthStr)
		if err == nil {
			depth = parsedDepth
		}
	}

	if depth > config.MaxDepth {
		depth = config.MaxDepth
	}

	includeEmployees := true
	if includeEmployeesStr := r.URL.Query().
		Get("include_employees"); includeEmployeesStr == "false" {
		includeEmployees = false
	}

	req := dto.GetDepartmentQuery{
		Depth:            depth,
		IncludeEmployees: new(includeEmployees),
	}

	if err = validation.V.Struct(&req); err != nil {
		zap.L().Debug("failed to validate request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.ctrl.GetDepartment(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ctrl.ErrNotFound):
			utils.ErrResponse(w, http.StatusNotFound, err)
			return
		default:
			utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
			return
		}
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}

func (h *Handler) updateDepartment(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.Context().Value(config.IDKey).(string)
	if !ok {
		utils.ErrResponse(w, http.StatusBadRequest, ErrRetrievePathVars)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	req := dto.UpdateDepartmentRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	res, err := h.ctrl.UpdateDepartment(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ctrl.ErrNotFound):
			utils.ErrResponse(w, http.StatusNotFound, err)
			return
		case errors.Is(err, ctrl.ErrDepartmentCycle):
			utils.ErrResponse(w, http.StatusConflict, err)
			return
		default:
			utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
			return
		}
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}

func (h *Handler) deleteDepartment(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.Context().Value(config.IDKey).(string)
	if !ok {
		utils.ErrResponse(w, http.StatusBadRequest, ErrRetrievePathVars)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	req := dto.DeleteDepartmentQuery{}
	req.Mode = r.URL.Query().Get("mode")

	reasIDStr := r.URL.Query().Get("reassign_to_department_id")

	reasID, err := strconv.ParseInt(reasIDStr, 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	req.ReassignToDepartmentID = int(reasID)
	if err = validation.V.Struct(&req); err != nil {
		zap.L().Debug("failed to validate request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.ctrl.DeleteDepartment(r.Context(), id, req)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
		return
	}

	utils.StatusResponse(w, http.StatusNoContent)
}
