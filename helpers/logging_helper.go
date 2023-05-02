package logr

import "log"

const (
	INFO  = "\033[35m"
	ERROR = "\033[31m"
	RESET = "\033[0m"
	WARN  = "\033[0;33m"
	DEBUG = "\033[034m"
)

func Info(in string) {
	log.Println(INFO + "[INFO]: " + RESET + in)
}

func Error(in string) {
	log.Println(ERROR + "[ERROR]: " + RESET + in)
}

func Warn(in string) {
	log.Println(WARN + "[WARNING]: " + RESET + in)
}

func Debug(in string) {
	log.Println(DEBUG + "[DEBUG]: " + RESET + in)
}
