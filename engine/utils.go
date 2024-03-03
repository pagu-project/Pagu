package engine

func boosterPrice(allPackages int) int {
	if allPackages < 100 {
		return 30
	} else if allPackages < 200 {
		return 40
	}
	return 50
}
