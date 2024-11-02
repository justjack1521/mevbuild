package patch

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	errInvalidVersionFormat = func(str string) error {
		return fmt.Errorf("invalid version format: %s", str)
	}
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) Zero() bool {
	return v == Version{}
}

func NewVersion(str string) (Version, error) {

	var parts = strings.Split(str, ".")
	if len(parts) != 3 {
		return Version{}, errInvalidVersionFormat(str)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, err
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil

}

func (v Version) IsNewMinorVersion() bool {
	return v.Patch == 0
}

func (v Version) GeneratePreviousVersions(steps int) []Version {
	var versions []Version

	for i := 0; i < steps; i++ {
		if v.Patch > 0 {
			v.Patch--
		} else if v.Minor > 0 {
			v.Minor--
			v.Patch = 9
		} else if v.Major > 0 {
			v.Major--
			v.Minor = 9
			v.Patch = 9
		}
		if v.Zero() {
			return versions
		}
		versions = append(versions, v)
	}
	return versions
}

func (v Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Equal(n Version) bool {
	return v.Major == n.Major && v.Minor == n.Minor && v.Patch == n.Patch
}
