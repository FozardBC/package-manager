package verManager

import (
	"fmt"
	"strconv"
	"strings"
)

type Version [2]int

func ParseVersion(v string) (Version, error) {
	parts := strings.Split(v, ".")
	if len(parts) > 3 {
		parts = parts[:3]
	}
	var ver Version
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return ver, fmt.Errorf("invalid version part: %s", p)
		}
		ver[i] = n
	}
	return ver, nil
}

func (v1 Version) GreaterEqual(v2 Version) bool {
	if v1[0] != v2[0] {
		return v1[0] >= v2[0]
	}

	return v1[1] >= v2[1]

}

func (v1 Version) LessEqual(v2 Version) bool {
	if v1[0] != v2[0] {
		return v1[0] <= v2[0]
	}

	return v1[1] <= v2[1]
}
