package creator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate_readGetParts(t *testing.T) {
	assert.Equal(t, 0, len(getParts("")))

	assert.Equal(t, 2, len(getParts(`
test1: 1
---
test2: 2
`)))

	assert.Equal(t, 1, len(getParts(`
test:1
`)))

	assert.Equal(t, 2, len(getParts(`
test:1
---
test:1
---
`)))

}
