package aws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegionLocation(t *testing.T) {
	assert.Equal(t, regionLocation("eu-west-1"), "EU (Ireland)")
	assert.Equal(t, regionLocation("us-east-1"), "US East (N. Virginia)")
}
