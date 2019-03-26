package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"testing"
)

func TestVersion(t *testing.T) {
	var writer = new(bytes.Buffer)

	cmd := NewCmdRoot(util.NewFactory(), writer)
	cmd.SetArgs([]string{"version"})

	_, err := cmd.ExecuteC()

	assert.Nil(t, err)

	assert.Equal(t, "Version: version, GitCommit: sha, GitBranch: branch, BuildDate: now\n", writer.String())
}
