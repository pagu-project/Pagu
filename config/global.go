package config

var (
	TargetMaskMain      = 1
	TargetMaskTest      = 2
	TargetMaskModerator = 4

	TargetMaskAll = TargetMaskMain | TargetMaskTest | TargetMaskModerator
)

const (
	PriceCacheKey = "PriceCacheKey"
)
