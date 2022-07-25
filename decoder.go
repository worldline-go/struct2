package struct2

// Decoder is main struct of struct2, holds config and functions.
type Decoder struct {
	// Tagname to lookup struct's field tag.
	TagName string // default is 'struct'
	// Hooks function run before decode and enable to change of data.
	Hooks []HookFunc
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
