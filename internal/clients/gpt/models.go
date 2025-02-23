package gpt

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

func NewRequest(prompt string) Request {
	return Request{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful coresearcher.",
			},
			{
				Role:    "user",
				Content: prompt + "Answer only, please with format question_id | answer.",
			},
		},
	}
}
