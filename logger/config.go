package logger

type Output string

const (
	OutputDiscard = "discard"
	OutputJson    = "json"
	OutputText    = "text"
)

type Config struct {
	Output Output `json:"output"`
}

func NewConfig() Config {
	return Config{
		Output: OutputText,
	}
}
