package database

import "fmt"

type MigrationError struct {
	Reason string
}

func (e MigrationError) Error() string {
	return fmt.Sprintf("migration error: %s", e.Reason)
}
