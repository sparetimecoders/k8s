package kops

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConfig(t *testing.T) {

	ig := `apiVersion: kops/v1alpha2
kind: InstanceGroup
metadata:
  creationTimestamp: 2019-01-23T13:10:02Z
  labels:
    kops.k8s.io/cluster: es-cluster.aws-dspa.volvocars.biz
  name: nodes
spec:
  cloudLabels:
    k8s.io/cluster-autoscaler/enabled: "true"
    k8s.io/cluster-autoscaler/es-cluster.aws-dspa.volvocars.biz: "true"
  image: kope.io/k8s-1.11-debian-stretch-amd64-hvm-ebs-2018-08-17
  machineType: r5.xlarge
  maxPrice: "0.190500"
  maxSize: 6
  minSize: 1
  nodeLabels:
    kops.k8s.io/instancegroup: nodes
  role: Node
  subnets:
  - eu-west-1a
  - eu-west-1b
  - eu-west-1c
`

	parsed, _ := parseInstanceGroup([]byte(ig))

	assert.Equal(t, "nodes", parsed.ig.Metadata.Name)
}

func TestUpdatePrice(t *testing.T) {
	ig := `
spec:
  maxSize: 0
  minSize: 0
  nodeLabels:
    kops.k8s.io/instancegroup: nodes
  role: Node
  subnets:
  - eu-west-1a
`
	parsed, _ := parseInstanceGroup([]byte(ig))
	assert.Equal(t, "1.1000", parsed.MaxPrice(1.1).ig.Spec.MaxPrice)

}
