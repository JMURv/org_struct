package db

import (
	"context"
	"testing"
	"time"

	"github.com/JMURv/org-struct/internal/dto"
	md "github.com/JMURv/org-struct/internal/models"
	"github.com/JMURv/org-struct/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pressly/goose/v3"
	gpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17.4-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(gpg.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, _ := db.DB()

	err = goose.SetDialect("postgres")
	if err != nil {
		t.Fatal(err)
	}

	err = goose.Up(sqlDB, "../../../migrations")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func newTestRepo(t *testing.T) *Repository {
	db := newTestDB(t)

	return &Repository{
		conn: db,
	}
}

func seedDepartmentTree(t *testing.T, r *Repository) (root, child, grandchild md.Department) {
	t.Helper()

	root = md.Department{Name: "Root"}
	if err := r.conn.Create(&root).Error; err != nil {
		t.Fatal(err)
	}

	child = md.Department{Name: "Child", ParentID: &root.ID}
	if err := r.conn.Create(&child).Error; err != nil {
		t.Fatal(err)
	}

	grandchild = md.Department{Name: "GrandChild", ParentID: &child.ID}
	if err := r.conn.Create(&grandchild).Error; err != nil {
		t.Fatal(err)
	}

	return
}

func TestRepository_CreateDepartment(t *testing.T) {
	r := newTestRepo(t)

	req := dto.CreateDepartmentRequest{
		Name: " Backend ",
	}

	res, err := r.CreateDepartment(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Name != "Backend" {
		t.Fatalf("expected trimmed name, got %s", res.Name)
	}
}

func TestRepository_CreateEmployee(t *testing.T) {
	r := newTestRepo(t)

	dept := md.Department{Name: "IT"}
	if err := r.conn.Create(&dept).Error; err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		req := dto.CreateDepartmentEmployeeRequest{
			FullName: " John ",
			Position: " Dev ",
		}

		res, err := r.CreateEmployee(context.Background(), dept.ID, req)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "John", res.FullName)
		assert.Equal(t, "Dev", res.Position)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := r.CreateEmployee(context.Background(), 999999, dto.CreateDepartmentEmployeeRequest{})
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}

func TestRepository_GetDepartment(t *testing.T) {
	r := newTestRepo(t)

	root, _, _ := seedDepartmentTree(t, r)

	t.Run("not_found", func(t *testing.T) {
		_, err := r.GetDepartment(context.Background(), 9999, dto.GetDepartmentQuery{})
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})

	t.Run("depth_0_only_root", func(t *testing.T) {
		res, err := r.GetDepartment(context.Background(), root.ID, dto.GetDepartmentQuery{
			Depth:            0,
			IncludeEmployees: new(true),
		})

		assert.NoError(t, err)
		assert.Equal(t, root.ID, res.ID)
		assert.Len(t, res.Children, 0)
	})

	t.Run("depth_2_recursive", func(t *testing.T) {
		b := true

		res, err := r.GetDepartment(context.Background(), root.ID, dto.GetDepartmentQuery{
			Depth:            2,
			IncludeEmployees: &b,
		})

		assert.NoError(t, err)
		assert.Len(t, res.Children, 1)
		assert.Len(t, res.Children[0].Children, 1)
	})
}

func TestRepository_UpdateDepartment(t *testing.T) {
	r := newTestRepo(t)

	a := md.Department{Name: "A"}
	b := md.Department{Name: "B"}
	c := md.Department{Name: "C"}

	r.conn.Create(&a)
	r.conn.Create(&b)
	r.conn.Create(&c)

	b.ParentID = &a.ID
	c.ParentID = &b.ID
	r.conn.Save(&b)
	r.conn.Save(&c)

	t.Run("rename", func(t *testing.T) {
		newName := "A1"

		res, err := r.UpdateDepartment(context.Background(), a.ID, dto.UpdateDepartmentRequest{
			Name: &newName,
		})

		assert.NoError(t, err)
		assert.Equal(t, "A1", res.Name)
	})

	t.Run("self_cycle", func(t *testing.T) {
		_, err := r.UpdateDepartment(context.Background(), a.ID, dto.UpdateDepartmentRequest{
			ParentID: &a.ID,
		})

		assert.ErrorIs(t, err, repo.ErrDepartmentCycle)
	})

	t.Run("deep_cycle_detected", func(t *testing.T) {
		_, err := r.UpdateDepartment(context.Background(), a.ID, dto.UpdateDepartmentRequest{
			ParentID: &c.ID,
		})

		assert.ErrorIs(t, err, repo.ErrDepartmentCycle)
	})
}

func TestRepository_DeleteDepartment(t *testing.T) {
	r := newTestRepo(t)

	root, _, _ := seedDepartmentTree(t, r)

	t.Run("cascade", func(t *testing.T) {
		err := r.DeleteDepartment(context.Background(), root.ID, dto.DeleteDepartmentQuery{
			Mode: "cascade",
		})

		assert.NoError(t, err)

		var count int64
		r.conn.Model(&md.Department{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	r = newTestRepo(t)
	d1 := md.Department{Name: "D1"}
	d2 := md.Department{Name: "D2"}
	r.conn.Create(&d1)
	r.conn.Create(&d2)

	emp := md.Employee{
		DepartmentID: d1.ID,
		FullName:     "John",
		Position:     "Dev",
	}
	r.conn.Create(&emp)

	t.Run("reassign", func(t *testing.T) {
		err := r.DeleteDepartment(context.Background(), d1.ID, dto.DeleteDepartmentQuery{
			Mode:                   "reassign",
			ReassignToDepartmentID: int(d2.ID),
		})

		assert.NoError(t, err)

		var e md.Employee
		r.conn.First(&e, emp.ID)

		assert.Equal(t, d2.ID, e.DepartmentID)
	})
}
