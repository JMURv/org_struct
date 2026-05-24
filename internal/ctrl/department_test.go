package ctrl

import (
	"context"
	"errors"
	"testing"

	"github.com/JMURv/org-struct/internal/dto"
	md "github.com/JMURv/org-struct/internal/models"
	"github.com/JMURv/org-struct/internal/repo"
	"github.com/JMURv/org-struct/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_CreateDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockAppRepo(ctrl)
	c := &Controller{repo: r}

	req := dto.CreateDepartmentRequest{
		Name: "Backend",
	}

	expected := &md.Department{ID: 1, Name: "Backend"}

	t.Run("success", func(t *testing.T) {
		r.EXPECT().
			CreateDepartment(gomock.Any(), req).
			Return(expected, nil)

		res, err := c.CreateDepartment(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestController_GetDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockAppRepo(ctrl)
	c := &Controller{repo: r}

	req := dto.GetDepartmentQuery{Depth: 1}

	expected := &md.Department{ID: 1}

	t.Run("success", func(t *testing.T) {
		r.EXPECT().
			GetDepartment(gomock.Any(), uint64(1), req).
			Return(expected, nil)

		res, err := c.GetDepartment(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("not_found", func(t *testing.T) {
		r.EXPECT().
			GetDepartment(gomock.Any(), uint64(1), req).
			Return(nil, repo.ErrNotFound)

		res, err := c.GetDepartment(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestController_UpdateDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockAppRepo(ctrl)
	c := &Controller{repo: r}

	ptr := func(s string) *string { return &s }
	req := dto.UpdateDepartmentRequest{
		Name: ptr("NewName"),
	}

	expected := &md.Department{ID: 1, Name: "NewName"}

	t.Run("success", func(t *testing.T) {
		r.EXPECT().
			UpdateDepartment(gomock.Any(), uint64(1), req).
			Return(expected, nil)

		res, err := c.UpdateDepartment(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("not_found", func(t *testing.T) {
		r.EXPECT().
			UpdateDepartment(gomock.Any(), uint64(1), req).
			Return(nil, repo.ErrNotFound)

		res, err := c.UpdateDepartment(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("cycle_error", func(t *testing.T) {
		r.EXPECT().
			UpdateDepartment(gomock.Any(), uint64(1), req).
			Return(nil, repo.ErrDepartmentCycle)

		res, err := c.UpdateDepartment(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.Equal(t, ErrDepartmentCycle, err)
	})

	t.Run("internal_error", func(t *testing.T) {
		r.EXPECT().
			UpdateDepartment(gomock.Any(), uint64(1), req).
			Return(nil, errors.New("boom"))

		res, err := c.UpdateDepartment(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.EqualError(t, err, "boom")
	})
}

func TestController_DeleteDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockAppRepo(ctrl)
	c := &Controller{repo: r}

	req := dto.DeleteDepartmentQuery{
		Mode: "cascade",
	}

	t.Run("success", func(t *testing.T) {
		r.EXPECT().
			DeleteDepartment(gomock.Any(), uint64(1), req).
			Return(nil)

		err := c.DeleteDepartment(context.Background(), 1, req)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		r.EXPECT().
			DeleteDepartment(gomock.Any(), uint64(1), req).
			Return(errors.New("boom"))

		err := c.DeleteDepartment(context.Background(), 1, req)

		assert.EqualError(t, err, "boom")
	})
}

func TestController_CreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockAppRepo(ctrl)
	c := &Controller{repo: r}

	req := dto.CreateDepartmentEmployeeRequest{
		FullName: "John",
		Position: "Dev",
	}

	expected := &md.Employee{ID: 1}

	t.Run("success", func(t *testing.T) {
		r.EXPECT().
			CreateEmployee(gomock.Any(), uint64(1), req).
			Return(expected, nil)

		res, err := c.CreateEmployee(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("not_found", func(t *testing.T) {
		r.EXPECT().
			CreateEmployee(gomock.Any(), uint64(1), req).
			Return(nil, repo.ErrNotFound)

		res, err := c.CreateEmployee(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("internal_error", func(t *testing.T) {
		r.EXPECT().
			CreateEmployee(gomock.Any(), uint64(1), req).
			Return(nil, errors.New("boom"))

		res, err := c.CreateEmployee(context.Background(), 1, req)

		assert.Nil(t, res)
		assert.EqualError(t, err, "boom")
	})
}
