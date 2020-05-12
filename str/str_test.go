package str

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCameSplit(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want []string
	}{
		"empty":            {in: "", want: []string{}},
		"lowercase":        {in: "lowercase", want: []string{"lowercase"}},
		"Class":            {in: "Class", want: []string{"Class"}},
		"MyClass":          {in: "MyClass", want: []string{"My", "Class"}},
		"MyC":              {in: "MyC", want: []string{"My", "C"}},
		"HTML":             {in: "HTML", want: []string{"HTML"}},
		"PDFLoader":        {in: "PDFLoader", want: []string{"PDF", "Loader"}},
		"AString":          {in: "AString", want: []string{"A", "String"}},
		"SimpleXMLParser":  {in: "SimpleXMLParser", want: []string{"Simple", "XML", "Parser"}},
		"vimRPCPlugin":     {in: "vimRPCPlugin", want: []string{"vim", "RPC", "Plugin"}},
		"GL11Version":      {in: "GL11Version", want: []string{"GL", "11", "Version"}},
		"99Bottles":        {in: "99Bottles", want: []string{"99", "Bottles"}},
		"May5":             {in: "May5", want: []string{"May", "5"}},
		"BFG9000":          {in: "BFG9000", want: []string{"BFG", "9000"}},
		"BöseÜberraschung": {in: "BöseÜberraschung", want: []string{"Böse", "Überraschung"}},
		"Two  spaces":      {in: "Two  spaces", want: []string{"Two", "  ", "spaces"}},
		"BadUTF8":          {in: "BadUTF8\xe2\xe2\xa1", want: []string{"BadUTF8\xe2\xe2\xa1"}},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := CamelSplit(tc.in)
			require.Equal(t, tc.want, got)
		})
	}
}
