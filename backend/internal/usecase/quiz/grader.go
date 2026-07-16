package quiz

import "findopedia/internal/entity"

type SubmittedAnswer struct {
	QuestionID    int64  `json:"question_id"`
	ChosenAnswer  string `json:"chosen_answer"`
}

func GradeAnswers(questions []entity.Question, submitted []SubmittedAnswer) (score, correct int) {
	answerMap := make(map[int64]string, len(submitted))
	for _, a := range submitted {
		answerMap[a.QuestionID] = a.ChosenAnswer
	}

	total := len(questions)
	if total == 0 {
		return 0, 0
	}

	for _, q := range questions {
		chosen, ok := answerMap[q.ID]
		if ok && chosen == q.CorrectAnswer {
			correct++
		}
	}

	score = (correct * 100) / total
	return score, correct
}
