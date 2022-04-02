package main

import (
	"encoding/json"
	"errors"
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
	Empty = flag.Bool("empty", false, "Leave the empty strings")
	Limit = flag.Uint("limit", 1, "Maximum number of elements in an array")
)

func filter(v any) any {
	switch x := v.(type) {
	default:
		return v
	case string:
		if x != "" || *Empty {
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
		if len(x) > 0 && *Limit > 0 {
			var vv []any
			for _, v := range x[:min(len(x), int(*Limit))] {
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

func run() error {
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("please specify a file name")
	}
	b, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		return err
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
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
