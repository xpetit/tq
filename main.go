package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
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

func run() (err error) {
	flag.Parse()
	input := io.Reader(os.Stdin)
	if flag.NArg() > 0 {
		var readers []io.Reader
		for _, arg := range flag.Args() {
			f, err := os.Open(arg)
			if err != nil {
				return err
			}
			defer f.Close()
			readers = append(readers, f)
		}
		input = io.MultiReader(readers...)
	}
	for dec := json.NewDecoder(input); dec.More(); {
		var v any
		if err := dec.Decode(&v); err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(filter(v)); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
