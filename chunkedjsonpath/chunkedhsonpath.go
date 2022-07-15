package chunkedjsonpath

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

var lang = []gval.Language{jsonpath.PlaceholderExtension()}
var evaluator = gval.Full(lang...)

// AppendLang - add new gval.Language to our evaluator (it is has base functions by default)
func AppendLang(l gval.Language) {
	lang = append(lang, l)
	evaluator = gval.Full(lang...)
}

// Chunk - Struct, describes one json map chunk
type Chunk struct {
	Key string    // Key - name of key in whole json
	R   io.Reader // source of
}

// Filter - iterate through source and filter all data chunks as a whole json map
func Filter(ctx context.Context, s chan Chunk, path string) (interface{}, error) {
	filter, err := evaluator.NewEvaluable(path)
	if err != nil {
		return nil, fmt.Errorf("bad jsonpath: %v", err)
	}

	ret := map[string]interface{}{}
	for i := range s {
		var chunk interface{}
		if err := json.NewDecoder(i.R).Decode(&chunk); err != nil {
			return nil, fmt.Errorf("can't decode source chunk '%s': %v", i.Key, err)
		}

		chunk = map[string]interface{}{
			i.Key: chunk,
		}

		res, err := filter(ctx, chunk)
		if err != nil {
			return nil, fmt.Errorf("can't evaluate on '%s': %v", i.Key, err)
		}

		ret[i.Key] = res
	}

	return ret, nil
}
