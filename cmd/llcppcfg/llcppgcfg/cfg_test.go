package llcppgcfg_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/goplus/llcppg/cmd/llcppcfg/llcppgcfg"
)

func TestNormalizePackageName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"lua5.4",
			"lua5_4",
		},
		{
			"tree-sitter",
			"tree_sitter",
		},
		{
			"python-3.12-embed",
			"python_3_12_embed",
		},
		{
			"-python-3.12-embed-",
			"python_3_12_embed",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			packageName := llcppgcfg.NormalizePackageName(tc.input)
			if !cmp.Equal(packageName, tc.expected) {
				t.Error(cmp.Diff(packageName, tc.expected))
				t.Fail()
			}
		})
	}
}
