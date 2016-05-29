package sol

import "testing"

func TestParseTag(t *testing.T) {
	var tag string
	var opts options

	// An empty tag should not error
	tag, opts = parseTag("")
	if tag != "" {
		t.Errorf("Unexpected tag, it should be an empty string: %s", tag)
	}

	// Multiple options
	tag, opts = parseTag(",omitempty,nullable")
	if tag != "" {
		t.Errorf("Unexpected tag, it should be an empty string: %s", tag)
	}
	expected := options{OmitEmpty, "nullable"}
	if !opts.Equals(expected) {
		t.Errorf("Unexpected options: %v != %v", opts, expected)
	}
}
