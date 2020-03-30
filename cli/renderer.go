package cli

type COLOR string

const (
	RESET string = "\033[0m"
	// 3-4 bits color
	BLACK   COLOR = "\033[0;30m"
	RED     COLOR = "\033[0;31m"
	GREEN   COLOR = "\033[0;32m"
	YELLOW  COLOR = "\033[0;33m"
	BLUE    COLOR = "\033[0;34m"
	MAGENTA COLOR = "\033[0;35m"
	CYAN    COLOR = "\033[0;36m"
	WHITE   COLOR = "\033[0;37m"
	// RGB colors
	GREY COLOR = "\033[0;38;2;128;128;128m"
)

type Renderer interface {
	Colorify(message interface{}, color COLOR) string
	Log(option LoggerOptions, message string) error
	Error(message string) error
	Warning(message string) error
	Success(message string) error
}

type LoggerOptions struct {
	Badge string
	Color COLOR
}

var (
	ERROR   = LoggerOptions{Badge: " ✖", Color: RED}
	WARNING = LoggerOptions{Badge: " ⚠", Color: YELLOW}
	SUCCESS = LoggerOptions{Badge: " ✔", Color: GREEN}
)

func NewRenderer() Renderer {
	return &ConsoleRenderer{}
}
