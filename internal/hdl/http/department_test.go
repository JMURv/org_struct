package http

import (
	"bytes"
	"context"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JMURv/org-struct/internal/config"
	"github.com/JMURv/org-struct/internal/ctrl"
	"github.com/JMURv/org-struct/internal/dto"
	"github.com/JMURv/org-struct/internal/hdl/http/utils"
	md "github.com/JMURv/org-struct/internal/models"
	"github.com/JMURv/org-struct/tests/mocks"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func ctxWithID(r *http.Request, id string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), config.IDKey, id))
}

func newReq(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer

	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestHandler_CreateDepartment(t *testing.T) {
	const uri = "/departments"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := &Handler{ctrl: mctrl}

	validReq := dto.CreateDepartmentRequest{
		Name: "Backend",
	}

	tests := []struct {
		name       string
		payload    any
		status     int
		expect     func()
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name:    "BadRequest_InvalidJSON",
			payload: "bad-json",
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(w *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(w.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},
		{
			name:    "Success",
			payload: validReq,
			status:  http.StatusCreated,
			expect: func() {
				mctrl.EXPECT().
					CreateDepartment(gomock.Any(), gomock.Any()).
					Return(&md.Department{}, nil)
			},
			assertions: func(w *httptest.ResponseRecorder) {
				var res md.Department
				err := json.NewDecoder(w.Body).Decode(&res)
				assert.NoError(t, err)
			},
		},
		{
			name:    "Internal",
			payload: validReq,
			status:  http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					CreateDepartment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal server error"))
			},
			assertions: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := newReq(t, http.MethodPost, uri, tt.payload)
			w := httptest.NewRecorder()

			h.createDepartment(w, req)

			assert.Equal(t, tt.status, w.Code)
			tt.assertions(w)
		})
	}
}

func TestHandler_CreateEmployee(t *testing.T) {
	const uri = "/departments/1/employees"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := &Handler{ctrl: mctrl}

	validReq := dto.CreateDepartmentEmployeeRequest{
		FullName: "John Doe",
		Position: "Backend Dev",
	}

	tests := []struct {
		name       string
		setupReq   func() *http.Request
		expect     func()
		status     int
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name: "MissingContextID",
			setupReq: func() *http.Request {
				return newReq(t, http.MethodPost, uri, validReq)
			},
			expect: func() {},
			status: http.StatusBadRequest,
			assertions: func(w *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(w.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},

		{
			name: "InvalidContextID",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPost, uri, validReq)
				ctx := context.WithValue(req.Context(), config.IDKey, "abc")
				return req.WithContext(ctx)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		{
			name: "InvalidJSON",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, uri, bytes.NewBufferString("bad-json"))
				req.Header.Set("Content-Type", "application/json")

				ctx := context.WithValue(req.Context(), config.IDKey, "1")
				return req.WithContext(ctx)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		{
			name: "ValidationError",
			setupReq: func() *http.Request {
				bad := dto.CreateDepartmentEmployeeRequest{
					FullName: "",
					Position: "",
				}

				req := newReq(t, http.MethodPost, uri, bad)
				ctx := context.WithValue(req.Context(), config.IDKey, "1")
				return req.WithContext(ctx)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		{
			name: "DepartmentNotFound",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPost, uri, validReq)
				ctx := context.WithValue(req.Context(), config.IDKey, "1")
				return req.WithContext(ctx)
			},
			expect: func() {
				mctrl.EXPECT().
					CreateEmployee(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, ctrl.ErrNotFound)
			},
			status: http.StatusNotFound,
		},

		{
			name: "InternalError",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPost, uri, validReq)
				ctx := context.WithValue(req.Context(), config.IDKey, "1")
				return req.WithContext(ctx)
			},
			expect: func() {
				mctrl.EXPECT().
					CreateEmployee(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			status: http.StatusInternalServerError,
		},

		{
			name: "Success",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPost, uri, validReq)
				ctx := context.WithValue(req.Context(), config.IDKey, "1")
				return req.WithContext(ctx)
			},
			expect: func() {
				mctrl.EXPECT().
					CreateEmployee(gomock.Any(), uint64(1), gomock.Any()).
					Return(&md.Employee{
						ID: 1,
					}, nil)
			},
			status: http.StatusCreated,
			assertions: func(w *httptest.ResponseRecorder) {
				var res md.Employee
				err := json.NewDecoder(w.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), res.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := tt.setupReq()
			w := httptest.NewRecorder()

			h.createEmployee(w, req)

			assert.Equal(t, tt.status, w.Code)

			if tt.assertions != nil {
				tt.assertions(w)
			}
		})
	}
}

func TestHandler_GetDepartment(t *testing.T) {
	const uri = "/departments/1"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := &Handler{ctrl: mctrl}

	tests := []struct {
		name       string
		setupReq   func() *http.Request
		expect     func()
		status     int
		assertions func(*httptest.ResponseRecorder)
	}{
		// -----------------------------------
		{
			name: "MissingContextID",
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, uri, nil)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------------
		{
			name: "InvalidID",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, uri, nil)
				return ctxWithID(req, "abc")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------------
		{
			name: "DefaultQueryValues",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, uri, nil)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(
						gomock.Any(),
						uint64(1),
						gomock.Any(),
					).
					DoAndReturn(func(_ context.Context, _ uint64, q dto.GetDepartmentQuery) (any, error) {
						assert.Equal(t, 1, q.Depth)
						assert.Equal(t, true, *q.IncludeEmployees)
						return &md.Department{}, nil
					})
			},
			status: http.StatusOK,
		},

		// -----------------------------------
		{
			name: "CustomQueryValues",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/departments/1?depth=3&include_employees=false",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(gomock.Any(), uint64(1), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ uint64, q dto.GetDepartmentQuery) (any, error) {
						assert.Equal(t, 3, q.Depth)
						assert.Equal(t, false, *q.IncludeEmployees)
						return &md.Department{}, nil
					})
			},
			status: http.StatusOK,
		},

		// -----------------------------------
		{
			name: "DepthAboveMaxClamped",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/departments/1?depth=999",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(gomock.Any(), uint64(1), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ uint64, q dto.GetDepartmentQuery) (any, error) {
						assert.LessOrEqual(t, q.Depth, config.MaxDepth)
						return &md.Department{}, nil
					})
			},
			status: http.StatusOK,
		},

		// -----------------------------------
		{
			name: "NotFound",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, uri, nil)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, ctrl.ErrNotFound)
			},
			status: http.StatusNotFound,
		},

		// -----------------------------------
		{
			name: "InternalError",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, uri, nil)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			status: http.StatusInternalServerError,
		},

		// -----------------------------------
		{
			name: "Success",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, uri, nil)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					GetDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(&md.Department{
						ID:   1,
						Name: "Backend",
					}, nil)
			},
			status: http.StatusOK,
			assertions: func(w *httptest.ResponseRecorder) {
				var res md.Department
				err := json.NewDecoder(w.Body).Decode(&res)

				assert.NoError(t, err)
				assert.Equal(t, uint64(1), res.ID)
				assert.Equal(t, "Backend", res.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := tt.setupReq()
			w := httptest.NewRecorder()

			h.getDepartment(w, req)

			assert.Equal(t, tt.status, w.Code)

			if tt.assertions != nil {
				tt.assertions(w)
			}
		})
	}
}

func TestHandler_UpdateDepartment(t *testing.T) {
	const uri = "/departments/1"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := &Handler{ctrl: mctrl}

	ptr := func(s string) *string { return &s }
	validReq := dto.UpdateDepartmentRequest{
		Name: ptr("Backend"),
	}

	tests := []struct {
		name       string
		setupReq   func() *http.Request
		expect     func()
		status     int
		assertions func(*httptest.ResponseRecorder)
	}{
		// -----------------------------
		{
			name: "MissingContextID",
			setupReq: func() *http.Request {
				return newReq(t, http.MethodPatch, uri, validReq)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "InvalidID",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPatch, uri, validReq)
				return ctxWithID(req, "abc")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "InvalidJSON",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, uri, bytes.NewBufferString("bad-json"))
				req.Header.Set("Content-Type", "application/json")
				return ctxWithID(req, "1")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "NotFound",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPatch, uri, validReq)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					UpdateDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, ctrl.ErrNotFound)
			},
			status: http.StatusNotFound,
		},

		// -----------------------------
		{
			name: "CycleConflict",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPatch, uri, validReq)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					UpdateDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, ctrl.ErrDepartmentCycle)
			},
			status: http.StatusConflict,
		},

		// -----------------------------
		{
			name: "InternalError",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPatch, uri, validReq)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					UpdateDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			status: http.StatusInternalServerError,
		},

		// -----------------------------
		{
			name: "Success",
			setupReq: func() *http.Request {
				req := newReq(t, http.MethodPatch, uri, validReq)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					UpdateDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(&md.Department{
						ID:   1,
						Name: "Backend",
					}, nil)
			},
			status: http.StatusOK,
			assertions: func(w *httptest.ResponseRecorder) {
				var res md.Department
				err := json.NewDecoder(w.Body).Decode(&res)

				assert.NoError(t, err)
				assert.Equal(t, uint64(1), res.ID)
				assert.Equal(t, "Backend", res.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := tt.setupReq()
			w := httptest.NewRecorder()

			h.updateDepartment(w, req)

			assert.Equal(t, tt.status, w.Code)

			if tt.assertions != nil {
				tt.assertions(w)
			}
		})
	}
}

func TestHandler_DeleteDepartment(t *testing.T) {
	const uri = "/departments/1"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := &Handler{ctrl: mctrl}

	tests := []struct {
		name       string
		setupReq   func() *http.Request
		expect     func()
		status     int
		assertions func(*httptest.ResponseRecorder)
	}{
		// -----------------------------
		{
			name: "MissingContextID",
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodDelete, uri, nil)
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "InvalidContextID",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, uri, nil)
				return ctxWithID(req, "abc")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "InvalidReassignID",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?mode=reassign&reassign_to_department_id=abc",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "ValidationError_ModeMissing",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?reassign_to_department_id=2",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "ValidationError_ReassignMissingWhenModeReassign",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?mode=reassign",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {},
			status: http.StatusBadRequest,
		},

		// -----------------------------
		{
			name: "CtrlInternalError",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?mode=cascade",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					DeleteDepartment(gomock.Any(), uint64(1), gomock.Any()).
					Return(errors.New("boom"))
			},
			status: http.StatusInternalServerError,
		},

		// -----------------------------
		{
			name: "Success_Cascade",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?mode=cascade",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					DeleteDepartment(gomock.Any(), uint64(1), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ uint64, q dto.DeleteDepartmentQuery) error {
						assert.Equal(t, "cascade", q.Mode)
						assert.Equal(t, 0, q.ReassignToDepartmentID)
						return nil
					})
			},
			status: http.StatusNoContent,
		},

		// -----------------------------
		{
			name: "Success_Reassign",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/departments/1?mode=reassign&reassign_to_department_id=2",
					nil,
				)
				return ctxWithID(req, "1")
			},
			expect: func() {
				mctrl.EXPECT().
					DeleteDepartment(gomock.Any(), uint64(1), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ uint64, q dto.DeleteDepartmentQuery) error {
						assert.Equal(t, "reassign", q.Mode)
						assert.Equal(t, 2, q.ReassignToDepartmentID)
						return nil
					})
			},
			status: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := tt.setupReq()
			w := httptest.NewRecorder()

			h.deleteDepartment(w, req)

			assert.Equal(t, tt.status, w.Code)

			if tt.assertions != nil {
				tt.assertions(w)
			}
		})
	}
}
