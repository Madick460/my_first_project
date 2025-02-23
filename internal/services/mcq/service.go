package mcq

import (
	"context"
	"fmt"
	"github.com/Rasikrr/my_project/internal/cache/answers"
	"github.com/Rasikrr/my_project/internal/clients/gpt"
	"github.com/Rasikrr/my_project/internal/domain/entity"
	"github.com/Rasikrr/my_project/internal/util"
	"github.com/samber/lo"
	"log"
	"strings"
)

var (
	questionTemplate = `
	Question ID: %s
	Question: %s
	Answer options: %s
	`
	prompt = `Answer only, please with format question_id | char and answer. 
example: question_id: 1 | answer: a. Car .
Or if question do not have char options - example: question_id: 2 |  answer: answer text
HERE is more examples of answers formats:
1. q420231:18 | e. Grief
2. q420231:16 | b. Cognitive Psychology
3. q420231:19 | True
4. q420231:14 | b. SMART goals`
)

type Service interface {
	GetAnswers(ctx context.Context, apiKey string, questions []*entity.MCQQuestion) ([]*entity.MCQAnswer, error)
}

type service struct {
	gptClient    gpt.Client
	answersCache answers.Cache
}

func NewService(gptClient gpt.Client, answersCache answers.Cache) Service {
	return &service{
		gptClient:    gptClient,
		answersCache: answersCache,
	}
}

func (s *service) GetAnswers(ctx context.Context, apiKey string, questions []*entity.MCQQuestion) ([]*entity.MCQAnswer, error) {
	questionIDs := lo.Map(questions, func(q *entity.MCQQuestion, _ int) string {
		return q.QuestionID
	})

	answersFromCache, err := s.answersCache.GetAnswers(ctx, questionIDs)
	if err != nil {
		log.Printf("answers not found: %v\n", err)
		return nil, err
	}

	for _, a := range answersFromCache {
		log.Printf("id of question from cache: %s | answer: %s\n", a.QuestionID, a.Answer)
	}

	if len(answersFromCache) == len(questions) {
		log.Println("all answers found in cache")
		return answersFromCache, nil
	}

	answeredMap := lo.SliceToMap(answersFromCache, func(a *entity.MCQAnswer) (string, *entity.MCQAnswer) {
		return a.QuestionID, a
	})

	filteredQuestions := make([]*entity.MCQQuestion, 0, len(questions)-len(answersFromCache))
	for _, q := range questions {
		if answeredMap[q.QuestionID] != nil {
			log.Printf("question %s already answered\n", q.QuestionID)
			continue
		}
		if q.ImageQuestion {
			log.Println("question is image question")
			filePath, err := util.DecodeImage(q.ImageBytes)
			if err != nil {
				return nil, err
			}
			imageText, err := util.GetImageText(filePath)
			if err != nil {
				return nil, err
			}
			log.Println("ImageText: ", imageText)
			q.Question = imageText
		}
		filteredQuestions = append(filteredQuestions, q)
	}
	prompt, err := s.createPrompt(ctx, filteredQuestions)
	if err != nil {
		return nil, err
	}

	answersRaw, err := s.gptClient.SendRequest(ctx, apiKey, prompt)
	if err != nil {
		return nil, err
	}

	answersFromGpt, err := s.parseAnswers(ctx, answersRaw)
	if err != nil {
		return nil, err
	}

	if err = s.answersCache.SaveAnswers(ctx, answersFromGpt); err != nil {
		return nil, err
	}

	answersFromGpt = append(answersFromGpt, answersFromCache...)
	return answersFromGpt, nil
}

func (s *service) createPrompt(_ context.Context, questions []*entity.MCQQuestion) (string, error) {
	b := strings.Builder{}
	for _, q := range questions {
		b.WriteString(fmt.Sprintf(questionTemplate, q.QuestionID, q.Question, q.AnswerOptions))
		b.WriteString("\n\n")
	}
	b.WriteString(prompt)
	return b.String(), nil
}

func (s *service) parseAnswers(_ context.Context, answer string) ([]*entity.MCQAnswer, error) {
	answers := strings.Split(answer, "\n")
	out := make([]*entity.MCQAnswer, 0, len(answers))
	for _, a := range answers {
		if a == "" {
			continue
		}
		parts := strings.Split(a, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid answer format: %s", a)
		}
		ans := &entity.MCQAnswer{
			QuestionID: strings.TrimSpace(parts[0]),
			Answer:     strings.TrimSpace(parts[1]),
		}
		out = append(out, ans)
	}
	return out, nil
}
