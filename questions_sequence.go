package main

import (
	"math/rand"
)

type QuestionsSequence struct {
	i         int
	questions []string
}

func newQuestionsSequence(questions []string) *QuestionsSequence {
	qs := &QuestionsSequence{
		i:         0,
		questions: make([]string, len(questions)),
	}
	copy(qs.questions, questions)
	return qs
}

func (qs *QuestionsSequence) Next() string {
	result := qs.questions[qs.i]
	qs.i = (qs.i + 1) % len(qs.questions)
	return result
}

func (qs *QuestionsSequence) Shuffle() {
	rand.Shuffle(len(qs.questions), func(i, j int) {
		qs.questions[i], qs.questions[j] = qs.questions[j], qs.questions[i]
	})
}
