package engine

import "fmt"

func CheckArgs(requiredArgs int, args []string) error {
	if len(args) != requiredArgs {
		return fmt.Errorf("incorrect number of arguments, expected %d but got %d", requiredArgs, len(args))
	}
	return nil
}

func boosterPrice(packagesCount int) int {
	if packagesCount < 100 {
		return 30
	} else if packagesCount < 200 {
		return 40
	}
	return 50
}
