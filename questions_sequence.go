package main

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
