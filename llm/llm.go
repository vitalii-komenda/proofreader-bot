package llm

type Role string

const (
	Proofread Role = "proofread"
	Slang     Role = "slang"
	Rephrase  Role = "rephrase"
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
	Proofread: {
		Role: "system",
		Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proposed**: $corrected_text`,
	},
	Slang: {
		Role: "system",
		Content: `You are bot to make slang version of the text. Users will be asking to make the text more slang. Reply with slang version and no explanations.
Format like this:
*Proposed*: $slang_text

Please strictly follow the format. It should have *Proposed* at the beginning of the message and then slang version.`,
	},
	Rephrase: {
		Role: "system",
		Content: `You are slang rephraser. Users will be asking to rephrase the text. Correct them with no explanations.
Format like this:
*Proposed*: $rephrased_text

Please strictly follow the format. It should have *Proposed* at the beginning of the message and then slang version.`,
	},
}
