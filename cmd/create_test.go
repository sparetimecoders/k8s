package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreate_NonExistingFile(t *testing.T) {
	var writer = new(bytes.Buffer)
	cmd := NewCmdRoot(util.NewMockFactory(), writer)
	cmd.SetArgs([]string{"create", "-f", "non-existing.yaml"})

	_, err := cmd.ExecuteC()

	assert.Nil(t, err, "Error: %v", err)

	assert.Equal(t, "open non-existing.yaml: no such file or directory", writer.String())
}

func TestCreate_ExistingCluster(t *testing.T) {
	var writer = new(bytes.Buffer)
	tempFile, _ := ioutil.TempFile(".", "abc")
	defer os.RemoveAll(tempFile.Name())

	_ = ioutil.WriteFile(tempFile.Name(), []byte(`
name: gotest
dnsZone: example.com
kubernetesVersion: 1.12.2
masters:
  zones:
  - a
cloudLabels: {}
`), os.ModeExclusive)

	factory := util.NewMockFactory()
	factory.ClusterExists = true

	cmd := NewCmdRoot(factory, writer)
	cmd.SetArgs([]string{"create", "-f", tempFile.Name()})

	_, err := cmd.ExecuteC()

	assert.Nil(t, err, "Error: %v", err)

	assert.Equal(t, "Cluster gotest.example.com already exists", writer.String())
}

func TestCreate_NonExistingCluster(t *testing.T) {
	var writer = new(bytes.Buffer)
	tempFile, _ := ioutil.TempFile(".", "abc")
	defer os.RemoveAll(tempFile.Name())

	_ = ioutil.WriteFile(tempFile.Name(), []byte(`
name: gotest
dnsZone: example.com
kubernetesVersion: 1.15.5
masters:
  zones:
  - a
cloudLabels: {}
`), os.ModeExclusive)

	factory := util.NewMockFactory()
	factory.ClusterExists = false

	factory.Handler.Responses <- "Version 1.15.5"
	factory.Handler.Responses <- ""
	factory.Handler.Responses <- ""
	go func() {
		cmd := NewCmdRoot(factory, writer)
		cmd.SetArgs([]string{"create", "-f", tempFile.Name()})

		_, err := cmd.ExecuteC()

		assert.Nil(t, err, "Error: %v", err)

		close(factory.Handler.Cmds)
	}()

	assert.Equal(t, "version", <-factory.Handler.Cmds)
	assert.Equal(t, "create cluster\n--name=gotest.example.com\n--node-count 2\n--zones eu-west-1a,eu-west-1b,eu-west-1c\n--master-zones eu-west-1a\n--node-size t3.medium\n--master-size t3.small\n--topology public\n--ssh-public-key ~/.ssh/id_rsa.pub\n--networking calico\n--encrypt-etcd-storage\n--authorization=RBAC\n--target=direct\n--cloud=aws\n--cloud-labels \n--network-cidr 172.21.0.0/22\n--kubernetes-version=1.15.5\n--dns-zone example.com", <-factory.Handler.Cmds)
	assert.Equal(t, "get ig nodes --name gotest.example.com -o yaml", <-factory.Handler.Cmds)
	assert.Equal(t, "replace ig  --name gotest.example.com -f -", <-factory.Handler.Cmds)
	assert.Equal(t, "get ig master-eu-west-1a --name gotest.example.com -o yaml", <-factory.Handler.Cmds)
	assert.Equal(t, "replace ig  --name gotest.example.com -f -", <-factory.Handler.Cmds)
	assert.Equal(t, "update cluster gotest.example.com --yes", <-factory.Handler.Cmds)
	assert.Equal(t, "validate cluster gotest.example.com", <-factory.Handler.Cmds)
}
