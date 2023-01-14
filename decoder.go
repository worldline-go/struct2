package struct2

import "regexp"

var reIgnoreSeperator = regexp.MustCompile(`[-_ ]`)

// Decoder is main struct of struct2, holds config and functions.
type Decoder struct {
	// Tagname to lookup struct's field tag, default is `struct`
	TagName string
	// Hooks function run before decode and enable to change of data.
	Hooks []HookFunc

	// Mapstructure operations

	// WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//   - negative numbers to overflowed uint values (base 10)
	//   - slice of maps to a merged map
	//   - single values are converted to slices if required. Each
	//     element is weakly decoded. For example: "4" can become []int{4}
	//     if the target type is an int slice.
	//
	WeaklyTypedInput bool

	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before decoded values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// Squash will squash embedded structs.  A squash tag may also be
	// added to an individual struct field using a tag.  For example:
	//
	//  type Parent struct {
	//      Child `struct:",squash"`
	//  }
	Squash bool
	// IgnoreUntaggedFields ignores all struct fields without explicit
	// TagName, comparable to `struct:"-"` as default behaviour.
	IgnoreUntaggedFields bool

	// BackupTagName usable if TagName not found.
	BackupTagName string

	// WeaklyDashUnderscore apply underscore/dash conversion to variables
	// on map to struct. variable_name == variable-name.
	WeaklyDashUnderscore bool

	// WeaklyIgnoreSeperator ignore seperator on map to struct. variable_name == variablename
	// values are -, _ and space.
	WeaklyIgnoreSeperator bool
}

func (d *Decoder) SetTagName(t string) *Decoder {
	d.TagName = t
	return d
}

func (d *Decoder) SetHooks(hooks []HookFunc) *Decoder {
	d.Hooks = hooks
	return d
}

func (d *Decoder) tagName() string {
	if d.TagName == "" {
		return "struct"
	}

	return d.TagName
}
