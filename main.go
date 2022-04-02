package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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
		if len(x) > 0 {
			m := map[string]any{}
			for k, v := range x {
				if v := filter(v); v != nil {
					m[k] = v
				}
			}
			if len(m) > 0 {
				return m
			}
		}
	case []any:
		if len(x) > 0 && *limit > 0 {
			var vv []any
			for _, v := range x[:min(len(x), int(*limit))] {
				if v := filter(v); v != nil {
					vv = append(vv, v)
				}
			}
			if len(vv) > 0 {
				return vv
			}
		}
	}
	return nil
}

func run() (err error) {
	flag.Parse()
	f := os.Stdin
	if flag.NArg() == 1 {
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			return err
		}
		defer f.Close()
	}
	var v any
	if err := json.NewDecoder(f).Decode(&v); err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	return enc.Encode(filter(v))
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
