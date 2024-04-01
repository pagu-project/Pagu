package database

import "fmt"

type MigrationError struct {
	Reason string
}

func (e MigrationError) Error() string {
	return fmt.Sprintf("migration error: %s", e.Reason)
}

type WriteError struct {
	Reason string
}

func (e WriteError) Error() string {
	return fmt.Sprintf("write error: %s", e.Reason)
}

type ReadError struct {
	Reason string
}

func (e ReadError) Error() string {
	return fmt.Sprintf("Read error: %s", e.Reason)
}
