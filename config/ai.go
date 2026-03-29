package config

const (
	key         = "KEY"
	model       = "MODEL"
	baseURL     = "BASE_URL"
	proxy       = "PROXY"
	temperature = "TEMPERATURE"
	maxTokens   = "MAX_TOKENS"
)

const (
	DefaultKey         = "local"
	DefaultModel       = "gemma3n:e4b"
	DefaultBaseURL     = "http://127.0.0.1:11434/v1"
	DefaultProxy       = ""
	DefaultTemperature = 0.2
	DefaultMaxTokens   = 1000
)

type AiConfig struct {
	key         string
	model       string
	baseURL     string
	proxy       string
	temperature float64
	maxTokens   int
}

func (c AiConfig) GetKey() string {
	return c.key
}

func (c AiConfig) GetModel() string {
	return c.model
}

func (c AiConfig) GetBaseURL() string {
	return c.baseURL
}

func (c AiConfig) GetProxy() string {
	return c.proxy
}

func (c AiConfig) GetTemperature() float64 {
	return c.temperature
}

func (c AiConfig) GetMaxTokens() int {
	return c.maxTokens
}
