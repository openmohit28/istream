package quiz

// AnswerMin/AnswerMax bound the Likert scale for every question.
const (
	AnswerMin = 1 // strongly disagree
	AnswerMax = 5 // strongly agree
)

// Questions is the fixed bank: four statements per RIASEC dimension.
var Questions = []Question{
	{ID: "r1", Text: "I enjoy working with my hands, tools, or machines.", Dimension: Realistic},
	{ID: "r2", Text: "I would rather fix a physical thing than write a report about it.", Dimension: Realistic},
	{ID: "r3", Text: "Being outdoors or on-site beats sitting at a desk all day.", Dimension: Realistic},
	{ID: "r4", Text: "I like seeing tangible, physical results from my work.", Dimension: Realistic},

	{ID: "i1", Text: "I enjoy digging into data or research to figure out why something happens.", Dimension: Investigative},
	{ID: "i2", Text: "Solving a complex problem is more satisfying than finishing a routine task.", Dimension: Investigative},
	{ID: "i3", Text: "I like learning how systems work under the hood.", Dimension: Investigative},
	{ID: "i4", Text: "I tend to question claims until I see the evidence.", Dimension: Investigative},

	{ID: "a1", Text: "I enjoy creating things - writing, design, music, or visuals.", Dimension: Artistic},
	{ID: "a2", Text: "I prefer open-ended work over tasks with fixed procedures.", Dimension: Artistic},
	{ID: "a3", Text: "I have strong instincts about aesthetics and style.", Dimension: Artistic},
	{ID: "a4", Text: "Expressing ideas in original ways energizes me.", Dimension: Artistic},

	{ID: "s1", Text: "Helping someone learn or grow feels deeply rewarding.", Dimension: Social},
	{ID: "s2", Text: "People often come to me with their problems.", Dimension: Social},
	{ID: "s3", Text: "I would enjoy a job built around caring for or supporting others.", Dimension: Social},
	{ID: "s4", Text: "I read the emotions in a room quickly.", Dimension: Social},

	{ID: "e1", Text: "I like persuading people and closing a deal.", Dimension: Enterprising},
	{ID: "e2", Text: "I am comfortable taking the lead when a group is stuck.", Dimension: Enterprising},
	{ID: "e3", Text: "Building a business or project from scratch excites me.", Dimension: Enterprising},
	{ID: "e4", Text: "I am willing to take calculated risks for a bigger payoff.", Dimension: Enterprising},

	{ID: "c1", Text: "I like organizing information, schedules, or budgets.", Dimension: Conventional},
	{ID: "c2", Text: "Clear rules and well-defined processes help me do my best work.", Dimension: Conventional},
	{ID: "c3", Text: "People trust me to get the details right.", Dimension: Conventional},
	{ID: "c4", Text: "I enjoy bringing order to messy situations.", Dimension: Conventional},
}

// QuestionByID indexes the bank for validation and scoring.
var QuestionByID = func() map[string]Question {
	m := make(map[string]Question, len(Questions))
	for _, q := range Questions {
		m[q.ID] = q
	}
	return m
}()
