package engine

func boosterPrice(allPackages int) int {
	switch {
	case allPackages < 100:
		return 30
	case allPackages < 200:
		return 40
	case allPackages < 300:
		return 50
	default:
		return 100
	}
}
