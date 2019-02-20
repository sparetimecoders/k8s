package kops

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/semver"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"io"
	"log"
	"os/exec"
	"strings"
)

type kops struct {
	stateStore string
	cmd        string
	debug      bool
}

func New(stateStore string) kops {
	k := kops{stateStore, "kops", true}
	return k
}

func (k kops) QueryCmd(paramString string, stdInData []byte) ([]byte, error) {
	cmd := exec.Command(k.cmd, k.buildParams(paramString)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out,nil
}

func (k kops) RunCmd(paramString string, stdInData []byte) error {
	cmd := exec.Command(k.cmd, k.buildParams(paramString)...)

	if stdInData != nil {
		cmd.Stdin = bytes.NewBuffer(stdInData)
	}

	out, _ := cmd.StdoutPipe()
	err, _ := cmd.StderrPipe()

	_ = cmd.Start()

	go func() {
		if k.debug {
			k.printOut(out)
		}
	}()
	go func() {
		k.printOut(err)
	}()

	return cmd.Wait()
}

func (k kops) CreateCluster(config config.ClusterConfig) (Cluster, error) {
	if ok := k.minimumKopsVersionInstalled(config.KubernetesVersion); ok == false {
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
	params := fmt.Sprintf(`create cluster
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
	)

	e := k.RunCmd(params, nil)
	if e != nil {
		return Cluster{}, e
	}
	return Cluster{name, k}, nil
}

func (k kops) buildParams(paramString string) []string {
	stateStore := []string{"--state", k.stateStore}
	return append(stateStore, strings.Split(strings.TrimSpace(
		strings.Replace(paramString, "\n", " ", -1)), " ")...)
}

func (k kops) getKopsVersion() (string, error) {
	cmd := exec.Command(k.cmd, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	_ = cmd.Start()
	_ = cmd.Wait()

	version := strings.TrimSpace(strings.TrimLeft(string(out), "Version: "))
	return version, nil
}

func (k kops) minimumKopsVersionInstalled(requiredKopsVersion string) bool {
	version, err := k.getKopsVersion()
	if err != nil {
		log.Printf("Failed to get kops version %s", err)
		return false
	}

	v, _ := semver.NewVersion(version)
	r, _ := semver.NewVersion(requiredKopsVersion)

	return v.Major() >= r.Major() && v.Minor() >= r.Minor()
}

func (k kops) printOut(out io.ReadCloser) {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
