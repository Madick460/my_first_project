package answers

import (
	"context"
	common "github.com/Rasikrr/my_project/internal/cache"
	"github.com/Rasikrr/my_project/internal/domain/entity"
	"log"
)

type Cache interface {
	SaveAnswers(ctx context.Context, answers []*entity.MCQAnswer) error
	GetAnswers(ctx context.Context, questionIDs []string) ([]*entity.MCQAnswer, error)
}

type cache struct {
	client common.Cache
}

func NewCache(client common.Cache) Cache {
	return &cache{
		client: client,
	}
}

func (c *cache) SaveAnswers(ctx context.Context, answers []*entity.MCQAnswer) error {
	keyValues := make([]any, len(answers)*2)
	for i, answer := range answers {
		keyValues[i*2] = answer.QuestionID
		keyValues[i*2+1] = answer.Answer
	}
	return c.client.MSet(ctx, keyValues...)
}

func (c *cache) GetAnswers(ctx context.Context, questionIDs []string) ([]*entity.MCQAnswer, error) {
	vals, err := c.client.MGet(ctx, questionIDs...)
	if err != nil {
		log.Println("error in get answers", err)
		return nil, err
	}

	var answers []*entity.MCQAnswer

	for i, v := range vals {
		if v == nil {
			continue
		}
		answerString, ok := v.(string)
		if !ok || answerString == "" {
			log.Println("answer is not string or empty")
			continue
		}
		answer := &entity.MCQAnswer{
			QuestionID: questionIDs[i],
			Answer:     answerString,
		}
		answers = append(answers, answer)

	}
	return answers, nil
}
