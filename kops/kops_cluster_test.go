package kops

import (
	"fmt"
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	a := `apiVersion: kops/v1alpha2
kind: Cluster
metadata:
  creationTimestamp: 2019-03-12T10:29:15Z
  name: gotest.k8s.sparetimecoders.com
spec:
  api:`

	fmt.Println(strings.Replace(a, "spec:", "spec:\n  additionalPolicies: \n    node: peter\n    master: ", 1))
}
