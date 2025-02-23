package entity

//go:generate easyjson -all MCQ.go

type MCQQuestion struct {
	QuestionID    string `json:"question_id"`
	Question      string `json:"question"`
	AnswerOptions string `json:"answer_options"`
	ImageQuestion bool   `json:"image_question"`
	ImageBytes    []byte `json:"image_bytes"`
}

type MCQAnswer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}
