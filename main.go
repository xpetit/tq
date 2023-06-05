package main

import (
	"encoding/json"
	"flag"
	"os"

	. "github.com/xpetit/x/v2"
)

var (
	empty = flag.Bool("empty", false, "Leave the empty strings")
	limit = flag.Uint("limit", 1, "Maximum number of elements in an array")
)

func filter(v any) any {
	switch x := v.(type) {
	default:
		return v
	case string:
		if x != "" || *empty {
			return x
		}
	case map[string]any:
		m := map[string]any{}
		for k, v := range x {
			if v := filter(v); v != nil {
				m[k] = v
			}
		}
		if len(m) > 0 {
			return m
		}
	case []any:
		var vv []any
		for _, v := range x {
			if len(vv) == int(*limit) {
				break
			}
			if v := filter(v); v != nil {
				vv = append(vv, v)
			}
		}
		if len(vv) > 0 {
			return vv
		}
	}
	return nil
}

func main() {
	flag.Parse()
	for dec := json.NewDecoder(os.Stdin); dec.More(); {
		var v any
		C(dec.Decode(&v))
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		C(enc.Encode(filter(v)))
	}
}
