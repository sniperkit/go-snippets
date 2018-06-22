package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/distributed-vision/go-resources/version/versiontype"
)

const (
	numbers  string = "0123456789"
	alphas          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum        = alphas + numbers
)

type Version interface {
	Equals(o Version) bool
	Compare(o Version) int
	Type() versiontype.VersionType
	String() string
	Validate() error
}

func New(major uint64, minor ...uint64) Version {

	if len(minor) == 0 {
		return NumericVersion(major)
	}

	var patch uint64

	if len(minor) > 1 {
		patch = minor[1]
	}

	return &SemanticVersion{major, minor[0], patch, nil, nil}
}

func NewPreRelease(major, minor, patch uint64, preReleaseInfo, buildInfo []string) Version {

	preReleaseValues := []PreReleaseValue{}

	if preReleaseInfo != nil {
		for _, prvalstr := range preReleaseInfo {
			prval, err := NewPreReleaseValue(prvalstr)
			if err != nil {
				preReleaseValues = append(preReleaseValues, prval)
			}
		}
	}

	return &SemanticVersion{major, minor, patch, preReleaseValues, buildInfo}
}

type VersionType versiontype.VersionType

type NumericVersion uint64

func (v NumericVersion) Equals(o Version) bool {
	return (v.Compare(o) == 0)
}

func (v NumericVersion) Compare(other Version) int {
	o, ok := other.(NumericVersion)

	if !ok {
		_, ok := other.(*SemanticVersion)
		if ok {
			return +1
		}

		return -1
	}

	if v != o {
		if v > o {
			return 1
		}
		return -1
	}

	return 0
}

func (v NumericVersion) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, uint64(v), 10)
	return string(b)
}

func (v NumericVersion) Bytes() []byte {
	var result []byte

	if v < 0xff {
		result = []byte{byte(v & 0xff)}
	} else if v < 0xffff {
		buf := [2]byte{0, 0}
		result = htons(buf[:], 0, uint16(v&0xffff))
	} else if v < 0xffffffff {
		buf := [4]byte{0, 0, 0, 0}
		result = htonl(buf[:], 0, uint32(v&0xffffffff))
	} else {
		buf := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
		result = htonll(buf[:], 0, uint64(v))
	}

	return result
}

func (v NumericVersion) ByteLength() byte {
	if v < 0xff {
		return 1
	} else if v < 0xffff {
		return 2
	} else if v < 0xffffffff {
		return 4
	}

	return 8
}

func (v NumericVersion) Type() versiontype.VersionType {
	return versiontype.NUMERIC
}

func (v NumericVersion) IsSemantic() bool {
	return false
}

func (v NumericVersion) Validate() error {
	return nil
}

type SemanticVersion struct {
	Major          uint64
	Minor          uint64
	Patch          uint64
	PreReleaseInfo []PreReleaseValue
	BuildInfo      []string
}

// Version to string
func (v *SemanticVersion) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Patch, 10)

	if len(v.PreReleaseInfo) > 0 {
		b = append(b, '-')
		b = append(b, v.PreReleaseInfo[0].String()...)

		for _, pre := range v.PreReleaseInfo[1:] {
			b = append(b, '.')
			b = append(b, pre.String()...)
		}
	}

	if len(v.BuildInfo) > 0 {
		b = append(b, '+')
		b = append(b, v.BuildInfo[0]...)

		for _, build := range v.BuildInfo[1:] {
			b = append(b, '.')
			b = append(b, build...)
		}
	}

	return string(b)
}

func (v *SemanticVersion) Bytes() []byte {
	return []byte(v.String())
}

func (v *SemanticVersion) Type() versiontype.VersionType {
	return versiontype.SEMANTIC
}

// Equals checks if v is equal to o.
func (v *SemanticVersion) Equals(o Version) bool {
	return (v.Compare(o) == 0)
}

// Compare compares Versions v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o
func (v *SemanticVersion) Compare(other Version) int {
	o, ok := other.(*SemanticVersion)

	if !ok {
		return -1
	}

	if v.Major != o.Major {
		if v.Major > o.Major {
			return 1
		}
		return -1
	}
	if v.Minor != o.Minor {
		if v.Minor > o.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != o.Patch {
		if v.Patch > o.Patch {
			return 1
		}
		return -1
	}

	// Quick comparison if a version has no prerelease versions
	if len(v.PreReleaseInfo) == 0 && len(o.PreReleaseInfo) == 0 {
		return 0
	} else if len(v.PreReleaseInfo) == 0 && len(o.PreReleaseInfo) > 0 {
		return 1
	} else if len(v.PreReleaseInfo) > 0 && len(o.PreReleaseInfo) == 0 {
		return -1
	}

	i := 0
	for ; i < len(v.PreReleaseInfo) && i < len(o.PreReleaseInfo); i++ {
		if comp := v.PreReleaseInfo[i].Compare(o.PreReleaseInfo[i]); comp == 0 {
			continue
		} else if comp == 1 {
			return 1
		} else {
			return -1
		}
	}

	// If all pr versions are the equal but one has further prversion, this one greater
	if i == len(v.PreReleaseInfo) && i == len(o.PreReleaseInfo) {
		return 0
	} else if i == len(v.PreReleaseInfo) && i < len(o.PreReleaseInfo) {
		return -1
	} else {
		return 1
	}

}

// Validate validates v and returns error in case
func (v *SemanticVersion) Validate() error {
	// Major, Minor, Patch already validated using uint64

	for _, pre := range v.PreReleaseInfo {
		if !pre.isNumeric { //Numeric prerelease versions already uint64
			if len(pre.strValue) == 0 {
				return fmt.Errorf("Prerelease can not be empty %q", pre.strValue)
			}
			if !containsOnly(pre.strValue, alphanum) {
				return fmt.Errorf("Invalid character(s) found in prerelease %q", pre.strValue)
			}
		}
	}

	for _, build := range v.BuildInfo {
		if len(build) == 0 {
			return fmt.Errorf("Build meta data can not be empty %q", build)
		}
		if !containsOnly(build, alphanum) {
			return fmt.Errorf("Invalid character(s) found in build meta data %q", build)
		}
	}

	return nil
}

// ParseTolerant allows for certain version specifications that do not strictly adhere to semver
// specs to be parsed by this library. It does so by normalizing versions before passing them to
// Parse(). It currently trims spaces, removes a "v" prefix, and adds a 0 patch number to versions
// with only major and minor components specified
func ParseTolerant(s string) (Version, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) < 3 {
		if strings.ContainsAny(parts[len(parts)-1], "+-") {
			return nil, errors.New("Short version cannot contain PreRelease/Build meta data")
		}
		for len(parts) < 3 {
			parts = append(parts, "0")
		}
		s = strings.Join(parts, ".")
	}

	return Parse(s)
}

// Parse parses version string and returns a validated Version or error
func Parse(s string) (Version, error) {
	if len(s) == 0 {
		return nil, errors.New("Version string empty")
	}

	if containsOnly(s, numbers) {
		ver, err := strconv.ParseUint(s, 10, 64)
		return NumericVersion(ver), err
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil, errors.New("No Major.Minor.Patch elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return nil, fmt.Errorf("Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return nil, fmt.Errorf("Major number must not contain leading zeroes %q", parts[0])
	}
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return nil, fmt.Errorf("Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return nil, fmt.Errorf("Minor number must not contain leading zeroes %q", parts[1])
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	v := SemanticVersion{}
	v.Major = major
	v.Minor = minor

	var build, prerelease []string
	patchStr := parts[2]

	if buildIndex := strings.IndexRune(patchStr, '+'); buildIndex != -1 {
		build = strings.Split(patchStr[buildIndex+1:], ".")
		patchStr = patchStr[:buildIndex]
	}

	if preIndex := strings.IndexRune(patchStr, '-'); preIndex != -1 {
		prerelease = strings.Split(patchStr[preIndex+1:], ".")
		patchStr = patchStr[:preIndex]
	}

	if !containsOnly(patchStr, numbers) {
		return nil, fmt.Errorf("Invalid character(s) found in patch number %q", patchStr)
	}
	if hasLeadingZeroes(patchStr) {
		return nil, fmt.Errorf("Patch number must not contain leading zeroes %q", patchStr)
	}
	patch, err := strconv.ParseUint(patchStr, 10, 64)
	if err != nil {
		return nil, err
	}

	v.Patch = patch

	// Prerelease
	for _, prstr := range prerelease {
		parsedPR, err := NewPreReleaseValue(prstr)
		if err != nil {
			return nil, err
		}
		v.PreReleaseInfo = append(v.PreReleaseInfo, parsedPR)
	}

	// Build meta data
	for _, str := range build {
		if len(str) == 0 {
			return nil, errors.New("Build meta data is empty")
		}
		if !containsOnly(str, alphanum) {
			return nil, fmt.Errorf("Invalid character(s) found in build meta data %q", str)
		}
		v.BuildInfo = append(v.BuildInfo, str)
	}

	return &v, nil
}

// MustParse is like Parse but panics if the version cannot be parsed.
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(`semver: Parse(` + s + `): ` + err.Error())
	}
	return v
}

// PreReleaseInfo represents PreRelease Version Info
type PreReleaseValue struct {
	strValue  string
	numValue  uint64
	isNumeric bool
}

// NewPRVersion creates a new valid prerelease version
func NewPreReleaseValue(s string) (PreReleaseValue, error) {
	if len(s) == 0 {
		return PreReleaseValue{}, errors.New("Prerelease is empty")
	}
	v := PreReleaseValue{}
	if containsOnly(s, numbers) {
		if hasLeadingZeroes(s) {
			return PreReleaseValue{}, fmt.Errorf("Numeric PreRelease version must not contain leading zeroes %q", s)
		}
		num, err := strconv.ParseUint(s, 10, 64)

		// Might never be hit, but just in case
		if err != nil {
			return PreReleaseValue{}, err
		}
		v.numValue = num
		v.isNumeric = true
	} else if containsOnly(s, alphanum) {
		v.strValue = s
		v.isNumeric = false
	} else {
		return PreReleaseValue{}, fmt.Errorf("Invalid character(s) found in prerelease %q", s)
	}
	return v, nil
}

// IsNumeric checks if prerelease-version is numeric
func (v PreReleaseValue) IsNumeric() bool {
	return v.isNumeric
}

// Compare compares two PreRelease Versions v and o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o
func (v PreReleaseValue) Compare(o PreReleaseValue) int {

	if v.isNumeric && !o.isNumeric {
		return -1
	} else if !v.isNumeric && o.isNumeric {
		return 1
	} else if v.isNumeric && o.isNumeric {
		if v.numValue == o.numValue {
			return 0
		} else if v.numValue > o.numValue {
			return 1
		} else {
			return -1
		}
	} else { // both are Alphas
		if v.strValue == o.strValue {
			return 0
		} else if v.strValue > o.strValue {
			return 1
		} else {
			return -1
		}
	}
}

// PreRelease version to string
func (v PreReleaseValue) String() string {
	if v.isNumeric {
		return strconv.FormatUint(v.numValue, 10)
	}
	return v.strValue
}

func NewBuildValue(s string) (string, error) {

	if len(s) == 0 {
		return "", errors.New("Buildversion is empty")
	}
	if !containsOnly(s, alphanum) {
		return "", fmt.Errorf("Invalid character(s) found in build meta data %q", s)
	}
	return s, nil
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}

func htonll(buffer []byte, index int, value uint64) []byte {
	buffer[index] = byte(0xff & (value >> 56))
	buffer[index+1] = byte(0xff & (value >> 48))
	buffer[index+2] = byte(0xff & (value >> 40))
	buffer[index+3] = byte(0xff & (value) >> 32)
	buffer[index+4] = byte(0xff & (value >> 24))
	buffer[index+5] = byte(0xff & (value >> 16))
	buffer[index+6] = byte(0xff & (value >> 8))
	buffer[index+7] = byte(0xff & (value))
	return buffer
}

func htonl(buffer []byte, index int, value uint32) []byte {
	buffer[index] = byte(0xff & (value >> 24))
	buffer[index+1] = byte(0xff & (value >> 16))
	buffer[index+2] = byte(0xff & (value >> 8))
	buffer[index+3] = byte(0xff & (value))
	return buffer
}

func htons(buffer []byte, index int, value uint16) []byte {
	buffer[index] = byte(0xff & (value >> 8))
	buffer[index+1] = byte(0xff & (value))
	return buffer
}
