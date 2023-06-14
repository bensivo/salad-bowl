package game

type Result string

const (
	RESULT_GUESSED Result = "GUESSED"
	RESULT_SKIPPED Result = "SKIPPED"
	RESULT_END     Result = "END"
)
