package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"io/ioutil"
	"os"
	"testing"
)

func TestDelete_NonExistingFile(t *testing.T) {
	var writer = new(bytes.Buffer)
	cmd := NewCmdRoot(util.NewMockFactory(), writer)
	cmd.SetArgs([]string{"delete", "-f", "non-existing.yaml"})

	_, err := cmd.ExecuteC()

	assert.Nil(t, err, "Error: %v", err)

	assert.Equal(t, "open non-existing.yaml: no such file or directory", writer.String())
}

func TestDelete_NonExistingCluster(t *testing.T) {
	var writer = new(bytes.Buffer)
	tempFile, _ := ioutil.TempFile(".", "abc")
	defer os.RemoveAll(tempFile.Name())

	_ = ioutil.WriteFile(tempFile.Name(), []byte(`
name: gotest
dnsZone: example.com
domain: example.com
kubernetesVersion: 1.12.2
masterZones:
  - a
cloudLabels: {}
`), os.ModeExclusive)

	cmd := NewCmdRoot(util.NewMockFactory(), writer)
	cmd.SetArgs([]string{"delete", "-f", tempFile.Name()})

	_, err := cmd.ExecuteC()

	assert.Nil(t, err, "Error: %v", err)

	assert.Equal(t, "Cluster gotest.example.com does not exist", writer.String())
}

func TestDelete_ExistingCluster(t *testing.T) {
	var writer = new(bytes.Buffer)
	tempFile, _ := ioutil.TempFile(".", "abc")
	defer os.RemoveAll(tempFile.Name())

	_ = ioutil.WriteFile(tempFile.Name(), []byte(`
name: gotest
dnsZone: example.com
domain: example.com
kubernetesVersion: 1.12.2
masterZones:
  - a
cloudLabels: {}
`), os.ModeExclusive)

	factory := util.NewMockFactory()
	factory.ClusterExists = true

	go func() {
		cmd := NewCmdRoot(factory, writer)
		cmd.SetArgs([]string{"delete", "-f", tempFile.Name()})

		_, err := cmd.ExecuteC()

		assert.Nil(t, err, "Error: %v", err)

		close(factory.Handler.Cmds)
	}()

	assert.Equal(t, "delete cluster\n--name=gotest.example.com\n--yes\n", <-factory.Handler.Cmds)
}
