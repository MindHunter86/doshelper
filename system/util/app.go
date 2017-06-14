package util


import "github.com/sirupsen/logrus"

type Application struct {
	Logout *logrus.Logger

	PTR_system System
	PTR_controller Controller
	PTR_model Model
}

type System interface {
	Destroy() error
}
type Controller interface {
	Destroy() error
}
type Model interface {
	Destroy() error
}
