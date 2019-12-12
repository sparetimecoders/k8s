package kops

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/semver"
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
	"io"
	"log"
	"os/exec"
	"strings"
)

type CmdHandler interface {
	QueryCmd(paramString string, stdInData []byte) ([]byte, error)
	RunCmd(paramString string, stdInData []byte) error
}

type osCmdHandler struct {
	stateStore string
	cmd        string
	debug      bool
}

type Kops interface {
	CreateCluster(config config.ClusterConfig) (Cluster, error)
	DeleteCluster(config config.ClusterConfig) error
	UpdateCluster() error
	ReplaceCluster(config string) error
	ValidateCluster() (string, bool)
	GetConfig() (string, error)
	ReplaceInstanceGroup(name string, data []byte) error
	GetInstanceGroup(name string) ([]byte, error)
	Version() (string, error)
	MinimumKopsVersionInstalled(requiredKopsVersion string) bool
}

type kops struct {
	Handler     CmdHandler
	ClusterName string
	_           struct{}
}

func New(name string, stateStore string) Kops {
	if !strings.HasPrefix(stateStore, "s3://") {
		stateStore = fmt.Sprintf("s3://%v", stateStore)
	}
	k := kops{ClusterName: name, Handler: osCmdHandler{stateStore, "kops", true}}
	return k
}

func (k osCmdHandler) QueryCmd(paramString string, stdInData []byte) ([]byte, error) {
	cmd := exec.Command(k.cmd, k.buildParams(paramString)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func (k osCmdHandler) RunCmd(paramString string, stdInData []byte) error {
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

func (k kops) ReplaceCluster(config string) error {
	log.Printf("Replacing cluster %v", k.ClusterName)
	return k.Handler.RunCmd(fmt.Sprintf("replace cluster %v -f -", k.ClusterName), []byte(config))
}

func (k kops) UpdateCluster() error {
	log.Printf("Updating cluster %v", k.ClusterName)
	return k.Handler.RunCmd(fmt.Sprintf("update cluster %v --yes", k.ClusterName), nil)
}

func (k kops) ValidateCluster() (string, bool) {
	out, err := k.Handler.QueryCmd(fmt.Sprintf("validate cluster %v", k.ClusterName), nil)
	if err == nil {
		return "", true
	}
	return string(out), false
}

func (k kops) GetInstanceGroup(name string) ([]byte, error) {
	return k.Handler.QueryCmd(fmt.Sprintf("get ig %v --name %v -o yaml", name, k.ClusterName), nil)
}

func (k kops) ReplaceInstanceGroup(name string, data []byte) error {
	return k.Handler.RunCmd(fmt.Sprintf("replace ig %v --name %v -f -", name, k.ClusterName), data)
}

func (k kops) GetConfig() (string, error) {
	params := fmt.Sprintf("get cluster %v -o yaml", k.ClusterName)
	out, err := k.Handler.QueryCmd(params, nil)
	return string(out), err
}

func (k kops) CreateCluster(clusterConfig config.ClusterConfig) (Cluster, error) {
	if ok := k.MinimumKopsVersionInstalled(clusterConfig.KubernetesVersion); ok == false {
		log.Fatalf("Installed version of kops can't handle requested kubernetes version (%s)", clusterConfig.KubernetesVersion)
	}
	name := clusterConfig.ClusterName()
	var zones []string
	for _, z := range clusterConfig.Nodes.Zones {
		zones = append(zones, fmt.Sprintf("%s%s", clusterConfig.Region, z))
	}
	var masterZones []string
	for _, z := range clusterConfig.Masters.Zones {
		masterZones = append(masterZones, fmt.Sprintf("%s%s", clusterConfig.Region, z))
	}
	var cloudLabels []string
	for k, v := range clusterConfig.CloudLabels {
		cloudLabels = append(cloudLabels, fmt.Sprintf("%s=%s", k, v))
	}
	params := fmt.Sprintf(`create cluster
--name=%s
--node-count %d
--zones %s
--master-zones %s
--node-size %s
--master-size %s
--topology public
--ssh-public-key %s
--networking calico
--encrypt-etcd-storage
--authorization=RBAC
--target=direct
--cloud=aws
--cloud-labels %s
--network-cidr %s
--kubernetes-version=%s
`,
		name,
		clusterConfig.Nodes.Max,
		strings.Join(zones, ","),
		strings.Join(masterZones, ","),
		clusterConfig.Nodes.InstanceType,
		clusterConfig.Masters.InstanceType,
		clusterConfig.SshKeyPath,
		strings.Join(cloudLabels, ","),
		clusterConfig.NetworkCIDR,
		clusterConfig.KubernetesVersion,
	)

	if clusterConfig.DnsZone != config.LocalCluster {
		params += fmt.Sprintf("--dns-zone %v", clusterConfig.DnsZone)
	}
	if clusterConfig.Vpc != "" {
		params += fmt.Sprintf("--vpc %v", clusterConfig.Vpc)
	}

	e := k.Handler.RunCmd(params, nil)
	if e != nil {
		return Cluster{}, e
	}
	return Cluster{kops: k}, nil
}

func (k kops) DeleteCluster(config config.ClusterConfig) error {
	name := config.ClusterName()

	params := fmt.Sprintf(`delete cluster
--name=%s
--yes
`,
		name,
	)

	e := k.Handler.RunCmd(params, nil)
	return e
}

func (k osCmdHandler) buildParams(paramString string) []string {
	stateStore := []string{"--state", k.stateStore}
	return append(stateStore, strings.Split(strings.TrimSpace(
		strings.Replace(paramString, "\n", " ", -1)), " ")...)
}

func (k kops) Version() (string, error) {
	out, err := k.Handler.QueryCmd("version", nil)
	if err != nil {
		return "", err
	}
	s := strings.Split(string(out), " ")[1]
	version := strings.TrimSpace(strings.TrimLeft(s, "Version: "))
	return version, nil
}

func (k kops) MinimumKopsVersionInstalled(requiredKopsVersion string) bool {
	version, err := k.Version()
	if err != nil {
		log.Printf("Failed to get kops version %v", err)
		return false
	}

	v, _ := semver.NewVersion(version)
	r, _ := semver.NewVersion(requiredKopsVersion)

	return v.Major() >= r.Major() && v.Minor() >= r.Minor()
}

func (k osCmdHandler) printOut(out io.ReadCloser) {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
