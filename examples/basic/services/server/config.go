package server

type Config struct {
	Value bool
	Text  string
}

func NewConfig() Config {
	return Config{
		Value: true,
		Text:  "test",
	}
}
