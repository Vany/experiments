package chunkedjsonpath

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	d := make(chan Chunk, 3)
	d <- Chunk{"n1", strings.NewReader(`{"v":10}`)}
	d <- Chunk{"n2", strings.NewReader(`{"v":20}`)}
	d <- Chunk{"n3", strings.NewReader(`{"v":30}`)}

	ctx := context.Background()

	t.Run("simple", func(t *testing.T) {
		got, err := Filter(ctx, d, "$.n1")
		assert.NoError(t, err)
		t.Log(got)
	})
}
