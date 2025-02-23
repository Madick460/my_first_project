package mcq_test

import (
	"github.com/Rasikrr/my_project/api"
	"net/http"
)

func (c *Controller) getMCQTestsAnswers(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	var req getMCQAnswersRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	questions := req.convertToEntity()
	s, err := c.mcqService.GetAnswers(ctx, apiKey, questions)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, err)
		return
	}
	api.SendData(w, convertAnswers(s), http.StatusOK)
}
