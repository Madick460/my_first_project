package mcq_test

import (
	"github.com/Rasikrr/my_project/internal/domain/entity"
	"github.com/samber/lo"
	"strings"
)

//go:generate easyjson -all models.go

type getMCQAnswersRequest struct {
	Questions []mcqQuestion `json:"questions"`
}

type answer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

func convertAnswers(answers []*entity.MCQAnswer) []*answer {
	out := make([]*answer, 0, len(answers))
	for _, a := range answers {
		ans := &answer{
			QuestionID: a.QuestionID,
		}
		ansStr := strings.TrimSpace(a.Answer)
		if strings.Contains(ansStr, ".") {
			ansStr = strings.TrimSpace(strings.Split(ansStr, ".")[0])
		}
		ans.Answer = ansStr
		out = append(out, ans)
	}
	return out
}

type mcqQuestion struct {
	QuestionID    string `json:"question_id"`
	QuestionText  string `json:"question_text"`
	AnswerOptions string `json:"answer_options"`
	ImageQuestion bool   `json:"image_question"`
	ImageBytes    []byte `json:"image_bytes"`
}

func (r *getMCQAnswersRequest) convertToEntity() []*entity.MCQQuestion {
	return lo.Map(r.Questions, func(q mcqQuestion, _ int) *entity.MCQQuestion {
		return q.convertToEntity()
	})
}

func (m *mcqQuestion) convertToEntity() *entity.MCQQuestion {
	return &entity.MCQQuestion{
		QuestionID:    m.QuestionID,
		Question:      m.QuestionText,
		AnswerOptions: m.AnswerOptions,
		ImageQuestion: m.ImageQuestion,
		ImageBytes:    m.ImageBytes,
	}
}
