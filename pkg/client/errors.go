package client

import "fmt"

type NotFoundError struct {
	Search  string
	Address string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found with %s address", e.Search, e.Address)
}

type NetworkInfoError struct {
	Reason string
}

func (e NetworkInfoError) Error() string {
	return e.Reason
}
