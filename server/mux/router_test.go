package mux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type a struct {
	Meta `path:"/a" method:"GET"`
}

func TestMetadata(t *testing.T) {
	path, method := Metadata(a{})
	assert.Equal(t, "/a", path)
	assert.Equal(t, "GET", method)
}
