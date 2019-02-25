package kops

import (
  "fmt"
  "log"
  "time"
)

type Cluster struct {
  name string
  kops kops
  _    struct{}
}

func GetCluster(name string, stateStore string) Cluster {
  return Cluster{name: name, kops: New(stateStore)}
}

func (c Cluster) CreateClusterResources() error {
  log.Printf("Creating cloud resources for %v", c.name)
  return c.kops.RunCmd(fmt.Sprintf("update cluster %v --yes", c.name), nil)
}

func (c Cluster) WaitForValidState(maxWaitSeconds int) bool {
  log.Printf("Validating cluster, will wait max %v seconds\n", maxWaitSeconds)
  fmt.Printf("Validating cluster, will wait max %v seconds\n", maxWaitSeconds)
  endTime := time.Now().Add(time.Second * time.Duration(maxWaitSeconds))
  done := false
  out := ""
  for time.Now().Before(endTime) {
    out, done = c.checkValidState()
    if done {
      log.Println("Cluster up and running")
      return true
    } else {
      time.Sleep(5 * time.Second)
      fmt.Printf(".")
    }
  }
  log.Printf("Failed to validate cluster in time, %v\n", out)
  return false
}

func (c Cluster) checkValidState() (string, bool) {
  out, err := c.kops.QueryCmd(fmt.Sprintf("validate cluster %v", c.name), nil)
  if err == nil {
    return "", true
  }
  return string(out), false
}
