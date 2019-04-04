package kops

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"testing"
)

func TestExistingPolicies(t *testing.T) {
	originalConfig := `apiVersion: kops/v1alpha2
kind: Cluster
metadata:
  creationTimestamp: 2019-03-12T10:29:15Z
  name: gotest.k8s.sparetimecoders.com
spec:
  additionalPolicies:
    node: '[{"Action":["route53:ChangeResourceRecordSets"],"Effect":"Allow","Resource":["arn:aws:route53:::hostedzone/*"]},{"Action":
["route53:ListHostedZones","route53:ListResourceRecordSets"],"Effect":"Allow","Resource":["*"]}]'
  api:
    dns: {}
  authorization:
    alwaysAllow: {}
  channel: stable
`

	policies := config.ClusterAutoscaler{}.Policies()
	updatedConfig := updateClusterConfigWithPolicies(originalConfig, policies)
	assert.Equal(t, "apiVersion: kops/v1alpha2\nkind: Cluster\nmetadata:\n  creationTimestamp: 2019-03-12T10:29:15Z\n  name: gotest.k8s.sparetimecoders.com\nspec:\n  additionalPolicies: \n    node: '[{\"Action\":[\"autoscaling:DescribeAutoScalingGroups\",\"autoscaling:DescribeAutoScalingInstances\",\"autoscaling:SetDesiredCapacity\",\"autoscaling:TerminateInstanceInAutoScalingGroup\",\"autoscaling:DescribeTags\"],\"Effect\":\"Allow\",\"Resource\":[\"*\"]}]'\n  api:\n    dns: {}\n  authorization:\n    alwaysAllow: {}\n  channel: stable\n",
		updatedConfig)
}
func TestNoPolicies(t *testing.T) {
	originalConfig := `apiVersion: kops/v1alpha2
kind: Cluster
metadata:
  creationTimestamp: 2019-03-12T10:29:15Z
  name: gotest.k8s.sparetimecoders.com
spec:
  api:
    dns: {}
  authorization:
    alwaysAllow: {}
  channel: stable
`

	policies := config.ClusterAutoscaler{}.Policies()
	updatedConfig := updateClusterConfigWithPolicies(originalConfig, policies)
	assert.Equal(t, "apiVersion: kops/v1alpha2\nkind: Cluster\nmetadata:\n  creationTimestamp: 2019-03-12T10:29:15Z\n  name: gotest.k8s.sparetimecoders.com\nspec:\n  additionalPolicies: \n    node: '[{\"Action\":[\"autoscaling:DescribeAutoScalingGroups\",\"autoscaling:DescribeAutoScalingInstances\",\"autoscaling:SetDesiredCapacity\",\"autoscaling:TerminateInstanceInAutoScalingGroup\",\"autoscaling:DescribeTags\"],\"Effect\":\"Allow\",\"Resource\":[\"*\"]}]'\n    \n  api:\n    dns: {}\n  authorization:\n    alwaysAllow: {}\n  channel: stable\n",
		updatedConfig)
}
