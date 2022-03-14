package struct2

import "strings"

// tagOptions contains a slice of tag options.
type tagOptions []string

// Has returns true if the given option is available in tagOptions.
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")

	return res[0], res[1:]
}
