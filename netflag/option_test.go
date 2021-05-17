package netflag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddress(t *testing.T) {
	type Set struct {
		Value    string
		Required bool
	}
	cases := []struct {
		In   []Option
		Want Set
	}{
		{nil, Set{"", true}},
		{[]Option{Address("1.1.1.1")}, Set{"1.1.1.1", false}},
	}
	for i, c := range cases {
		var cfg config
		cfg.apply(c.In)
		got := Set{
			Value:    cfg.address,
			Required: len(cfg.address) == 0,
		}
		if diff := cmp.Diff(c.Want, got); len(diff) > 0 {
			t.Errorf("[%d] -want +got\n%s", i, diff)
		}
	}
}

func TestNetwork(t *testing.T) {
	type Set struct {
		Predetermined bool
		Required      bool
		Usage         string
		Value         string
	}
	cases := []struct {
		In   []Option
		Want Set
	}{
		{nil, Set{true, false, "", "tcp"}},
		{[]Option{Network("tcp")}, Set{true, false, "", "tcp"}},
		{[]Option{Network("unix")}, Set{true, false, "", "unix"}},
		{[]Option{Network("*")}, Set{false, true, "", ""}},
		{[]Option{Network("udp", "*")}, Set{false, false, "", "udp"}},
		{[]Option{Network("unix", "tcp")}, Set{false, false, " `[unix|tcp]`", "unix"}},
	}
	for i, c := range cases {
		var cfg config
		cfg.apply(c.In)
		got := Set{
			Predetermined: cfg.networkPredetermined(),
			Required:      cfg.networkRequired(),
			Usage:         cfg.networkUsage(),
			Value:         cfg.networkValue(),
		}
		if diff := cmp.Diff(c.Want, got); len(diff) > 0 {
			t.Errorf("[%d] -want +got\n%s", i, diff)
		}
	}
}
