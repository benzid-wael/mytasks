package cli

import (
	"fmt"
)

type ConsoleRenderer struct{}

func (Logger *ConsoleRenderer) Log(option LoggerOptions, message string) error {
	_, err := fmt.Printf("%v  %v\n", Logger.Colorify(option.Badge, option.Color), message)
	return err
}

func (Logger *ConsoleRenderer) Error(message string) error {
	return Logger.Log(ERROR, message)
}

func (Logger *ConsoleRenderer) Warning(message string) error {
	return Logger.Log(WARNING, message)
}

func (Logger *ConsoleRenderer) Success(message string) error {
	return Logger.Log(SUCCESS, message)
}

func (Logger *ConsoleRenderer) Colorify(message interface{}, color COLOR) string {
	return fmt.Sprintf("%v%v%v", color, message, NC)
}
