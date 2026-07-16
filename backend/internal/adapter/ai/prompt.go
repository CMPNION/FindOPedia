package ai

import (
	"bytes"
	"text/template"
)

var quizPromptTmpl = template.Must(template.New("quiz").Parse(`You are a quiz master. Generate exactly 10 quiz questions to test detailed comprehension of this Wikipedia article.

Rules:
- 7 multiple_choice questions (4 options each: A, B, C, D)
- 3 true_false questions (options: True, False)
- Questions must be answerable solely from the article text
- Vary difficulty from easy to hard
- correct_answer must exactly match one option's value field

Article Title: {{.Title}}

Article Content:
{{.Content}}

Respond ONLY with a valid JSON array, no markdown, no explanation. Schema:
[
  {
    "index": 0,
    "type": "multiple_choice",
    "question_text": "...",
    "options": [
      {"label": "A", "value": "actual answer text"},
      {"label": "B", "value": "..."},
      {"label": "C", "value": "..."},
      {"label": "D", "value": "..."}
    ],
    "correct_answer": "actual answer text"
  },
  {
    "index": 1,
    "type": "true_false",
    "question_text": "...",
    "options": [
      {"label": "True", "value": "True"},
      {"label": "False", "value": "False"}
    ],
    "correct_answer": "True"
  }
]`))

func buildPrompt(title, content string) string {
	var buf bytes.Buffer
	_ = quizPromptTmpl.Execute(&buf, struct{ Title, Content string }{title, content})
	return buf.String()
}
