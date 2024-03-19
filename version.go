package robopac

import "fmt"

type Version struct {
	Meta  string `json:"meta"  xml:"meta"`
	Major uint8  `json:"major" xml:"major"`
	Minor uint8  `json:"minor" xml:"minor"`
	Patch uint8  `json:"patch" xml:"patch"`
}

var version = Version{
	Major: 0,
	Minor: 0,
	Patch: 1,
	Meta:  "beta",
}

func StringVersion() string {
	v := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Meta != "" {
		v = fmt.Sprintf("%s-%s", v, version.Meta)
	}

	return v
}
