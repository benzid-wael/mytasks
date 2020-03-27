package cli

type COLOR string

const (
	// 3-4 bits color
	NC      string = "\033[0m"
	BLACK          = "\033[0;30m"
	RED            = "\033[0;31m"
	GREEN          = "\033[0;32m"
	YELLOW         = "\033[0;33m"
	BLUE           = "\033[0;34m"
	MAGENTA        = "\033[0;35m"
	CYAN           = "\033[0;36m"
	WHITE          = "\033[0;37m"
	// RGB colors
	GREY = "\033[0;38;2;128;128;128m"
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
