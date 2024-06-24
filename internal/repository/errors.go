package repository

import "fmt"

type ConnectionError struct {
	Message string
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("connection error: %s", e.Message)
}

type MigrationError struct {
	Message string
}

func (e MigrationError) Error() string {
	return fmt.Sprintf("migration error: %s", e.Message)
}

type WriteError struct {
	Message string
}

func (e WriteError) Error() string {
	return fmt.Sprintf("write error: %s", e.Message)
}

type ReadError struct {
	Message string
}

func (e ReadError) Error() string {
	return fmt.Sprintf("Read error: %s", e.Message)
}
