package main

import (
  "bufio"
  "fmt"
  "github.com/Masterminds/semver"
  "io"
  "log"
  "os/exec"
  "strings"
)

func getKopsVersion() (string, error) {

  cmd := exec.Command("kops", "version")
  out, err := cmd.CombinedOutput()
  if err != nil {
    return "", err
  }
  _ = cmd.Start()
  _ = cmd.Wait()

  version := strings.TrimSpace(strings.TrimLeft(string(out), "Version: "))
  return version, nil
}

func minimumKopsVersionInstalled(requiredKopsVersion string) bool {
  version, err := getKopsVersion()
  if err != nil {
    fmt.Printf("Failed to get kops version %s", err)
    return false
  }

  v, _ := semver.NewVersion(version)
  r, _ := semver.NewVersion(requiredKopsVersion)

  return v.Major() >= r.Major() && v.Minor() >= r.Minor()
}

func CreateCluster(config Cluster) error {
  if ok := minimumKopsVersionInstalled(config.KubernetesVersion); ok == false {
    log.Fatalf("Installed version of kops can't handle requested kubernetes version (%s)", config.KubernetesVersion)
  }
  name := fmt.Sprintf("%s.%s", config.Name, config.DnsZone)
  zones := fmt.Sprintf("%[1]sa,%[1]sb,%[1]sc", config.Region)
  var masterZones []string
  for _, z := range config.MasterZones {
    masterZones = append(masterZones, fmt.Sprintf("%s%s", config.Region, z))
  }
  var cloudLabels []string
  for k, v := range config.CloudLabels {
    cloudLabels = append(cloudLabels, fmt.Sprintf("%s=%s", k, v))
  }
  params := strings.TrimSpace(strings.Replace(fmt.Sprintf(`create cluster
--name=%s
--node-count %d
--zones %s
--master-zones %s
--dns-zone %s
--node-size %s
--master-size %s
--topology public
--ssh-public-key %s
--networking calico
--encrypt-etcd-storage
--authorization=AlwaysAllow
--target=direct
--cloud=aws
--cloud-labels %s
--network-cidr %s
--kubernetes-version=%s
--state=s3://929619993638-aws-dspa.volvocars.biz-kops-storage
`,
    name,
    config.Nodes.Max,
    zones,
    strings.Join(masterZones, ","),
    config.DnsZone,
    config.Nodes.InstanceType,
    config.MasterInstanceType,
    config.SshKeyPath,
    strings.Join(cloudLabels, ","),
    config.NetworkCIDR,
    config.KubernetesVersion,
  ),
    "\n", " ", -1))
  fmt.Println(params)
  cmd := exec.Command("kops", strings.Split(params, " ")...)

  out, _ := cmd.StdoutPipe()
  err, _ := cmd.StderrPipe()

  _ = cmd.Start()

  go func() {
    printOut(out)
  }()
  go func() {
    printOut(err)
  }()

  return cmd.Wait()
}

func printOut(out io.ReadCloser) {
  scanner := bufio.NewScanner(out)
  scanner.Split(bufio.ScanLines)

  for scanner.Scan() {
    fmt.Println(scanner.Text())
  }
}
