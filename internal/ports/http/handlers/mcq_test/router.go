package mcq_test

import (
	"github.com/Rasikrr/my_project/internal/ports/http/middlewares"
	mcqS "github.com/Rasikrr/my_project/internal/services/mcq"
	"net/http"
)

type Controller struct {
	mcqService mcqS.Service
	m          *middlewares.AuthMiddleware
}

func NewController(mcqService mcqS.Service, m *middlewares.AuthMiddleware) *Controller {
	return &Controller{
		mcqService: mcqService,
		m:          m,
	}
}

func (c *Controller) Init(r *http.ServeMux) {
	r.HandleFunc("POST /get_answers", c.m.Handle(c.getMCQTestsAnswers))
}
