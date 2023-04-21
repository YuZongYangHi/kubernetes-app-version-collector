package cmd

import (
	"errors"
	"fmt"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/queue"
	"os/exec"
	"strings"
)

const (
	Docker     = "docker"
	Containerd = "containerd"
)

type Cmd struct {
	criCommand []string
}

func (c *Cmd) runCommand(cmd string) (string, error) {
	name := "/bin/bash"
	args := "-c"
	res := exec.Command(name, args, cmd)
	output, err := res.Output()
	if err != nil {
		return "", err
	}
	return strings.Replace(string(output), "\n", "", -1), nil
}

func (c *Cmd) run(name string) ([]*queue.CollectMetrics, error) {
	var result []*queue.CollectMetrics

	versionCmd := fmt.Sprintf(c.criCommand[0], name)
	sha256Cmd := fmt.Sprintf(c.criCommand[1], name)

	s1, err := c.runCommand(versionCmd)
	if err != nil {
		return nil, err
	}
	s2, err := c.runCommand(sha256Cmd)
	if err != nil {
		return nil, err
	}

	result = append(result, &queue.CollectMetrics{
		Name:    name,
		Type:    Tag,
		Version: s1,
	})
	result = append(result, &queue.CollectMetrics{
		Name:    name,
		Type:    Sha256,
		Version: s2,
	})
	return result, nil
}

func (c *Cmd) Run(name, cmd string) ([]*queue.CollectMetrics, error) {

	var result []*queue.CollectMetrics

	if cmd == "" {
		return c.run(name)
	}

	s, err := c.runCommand(cmd)
	if err != nil {
		return nil, err
	}

	result = append(result, &queue.CollectMetrics{
		Name:    name,
		Type:    Tag,
		Version: s,
	})
	return result, nil
}

func NewCmd(runtime string) (*Cmd, error) {
	var cri []string
	switch runtime {
	case Containerd:
		cri = append(cri, ContainerdGetTagCommand)
		cri = append(cri, ContainerdGetSha256Command)
	case Docker:
		cri = append(cri, DockerGetTagCommand)
		cri = append(cri, DockerGetSha254Command)
	default:
		return nil, errors.New("invalid cri engine")
	}
	return &Cmd{cri}, nil
}
