package kongx

import (
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

type TestSample struct {
	Name string `json:"name"`
	Game string `json:"game"`
}

func TestJSONFileMapper(t *testing.T) {
	var cli struct {
		Sample TestSample `type:"jsonfile"`
	}
	opt := kong.NamedMapper("jsonfile", JSONFileMapper)
	parser, err := kong.New(&cli, opt)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--sample", "testdata/sample.json"})
	require.NoError(t, err)

	want := TestSample{Name: "Lee Sedol", Game: "Go"}
	require.Equal(t, want, cli.Sample)
}

func TestJSONFileMapperErr(t *testing.T) {
	var cli struct {
		Sample TestSample `type:"jsonfile"`
	}
	opts := []kong.Option{
		kong.NamedMapper("jsonfile", JSONFileMapper),
		kong.Exit(func(int) { fmt.Println("EXIT") }),
	}
	parser, err := kong.New(&cli, opts...)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--sample", "testdata/MISSING_FILE.json"})
	require.Error(t, err)

	_, err = parser.Parse([]string{"--sample"})
	require.Error(t, err)
}
