package ctrl

type AppRepo interface{ departmentRepo }

type AppCtrl interface{ departmentCtrl }

type Controller struct{ repo AppRepo }

func New(repo AppRepo) *Controller {
	return &Controller{repo: repo}
}
