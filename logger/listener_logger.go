package logger

import (
	"log"
	"strings"
	"users_service/config"
)

type listenerLogger struct {
	name string
}

// ListenerLogger is the logger for the listener
var ListenerLogger = listenerLogger{name: config.ListenerLoggerName}

// buildSuccessMessage creates a string that contains a success message
func (l listenerLogger) buildSuccessMessage(message string) string {
	return l.name + " " + message
}

// buildErrorMessage creates a string that contains an error message
func (l listenerLogger) buildErrorMessage(message string, err error) string {
	return strings.Join([]string{l.name, message, err.Error()}, " ")
}

// FailedToListen logs an error message when the grpc_server fails to listen
func (l listenerLogger) FailedToListen(err error) {
	message := l.buildErrorMessage("Failed to listen", err)
	log.Fatalf(message)
}

// ServerStarted logs a success message when the grpc_server starts
func (l listenerLogger) ServerStarted(port string) {
	message := l.buildSuccessMessage("Server started on: " + port)
	log.Println(message)
}

// FailedToServe logs an error message when the grpc_server fails to serve
func (l listenerLogger) FailedToServe(err error) {
	message := l.buildErrorMessage("Failed to serve", err)
	log.Fatalf(message)
}
