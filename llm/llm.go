package llm

type LLM interface {
	SendRequest(text string) (string, error)
}
