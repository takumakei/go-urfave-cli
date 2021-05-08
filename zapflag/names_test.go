package zapflag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestShortFor(t *testing.T) {
	cases := []struct {
		In   string
		Want string
	}{
		{"log-type", "lt"},
		{"-log-type", "lt"},
		{"--log-type", "lt"},
		{"--log--type", "lt"},

		{"log-stack-trace", "lst"},
		{"-log-stack-trace", "lst"},
		{"--log-stack-trace", "lst"},
		{"--log--stack--trace", "lst"},
	}

	for i, c := range cases {
		got := shortFor(c.In)
		if diff := cmp.Diff(c.Want, got); len(diff) > 0 {
			t.Errorf("[%d] %q\n-want +got\n%s", i, c.In, diff)
		}
	}
}
