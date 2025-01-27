package logger

import (
	"github.com/gookit/color"
	"log"
)

// Error error output
func Error(err error) {
	log.Printf("%v\n", color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err))
}

// Warning warning message
func Warning(out ...any) {
	log.Println(color.New(color.Yellow).Sprint(out...))
}
func WarningF(format string, out ...any) {
	log.Println(color.New(color.Yellow).Sprintf(format, out...))
}

// Success success message
func Success(out ...any) {
	log.Println(color.New(color.Green).Sprint(out...))
}
func SuccessF(format string, out ...any) {
	log.Println(color.New(color.Green).Sprintf(format, out...))
}

// Info info message
func Info(out ...any) {
	log.Println(color.New(color.Cyan).Sprint(out...))
}
func InfoF(format string, out ...any) {
	log.Println(color.New(color.Cyan).Sprintf(format, out...))
}
