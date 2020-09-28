// Package kongx provides utility functions to github.com/alecthomas/kong
// such as parsing field value from a JSON file via the CustomMapper
// JSONFileMapper.
package kongx

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/alecthomas/kong"
)

// JSONFileMapper implements kong.MapperValue to decode a JSON file into
// a struct field.
//
//    var cli struct {
//      Profile Profile `type:"jsonfile"`
//    }
//
//    func main() {
//      kong.Parse(&cli, kong.NamedMapper("jsonfile", kongx.JSONFileMapper))
//    }
var JSONFileMapper = kong.MapperFunc(decodeJSONFile)

func decodeJSONFile(ctx *kong.DecodeContext, target reflect.Value) error {
	var fname string
	if err := ctx.Scan.PopValueInto("filename", &fname); err != nil {
		return err
	}
	f, err := os.Open(fname) //nolint:gosec
	if err != nil {
		return err
	}
	defer f.Close() //nolint
	return json.NewDecoder(f).Decode(target.Addr().Interface())
}
