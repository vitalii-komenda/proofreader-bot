package llm

type Role string

const (
	Proofreader Role = "proofreader"
	Slang       Role = "slang"
)

type LLM interface {
	SendRequest(text string, role Role) (string, error)
	Init() LLM
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

func Init(l LLM) LLM {
	return l.Init()
}

var Roles = map[Role]Message{
	Proofreader: {
		Role: "system",
		Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proofread*: $whole_corrected_text`,
	},
	Slang: {
		Role: "system",
		Content: `You are bot to make slang version of the text. Users will be asking to make the text more slang. Correct them with no explanations.
Format like this:
*Lowkey*: $whole_slang_text`,
	},
}
