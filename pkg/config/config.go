package config

import "github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/util/parsers"

type AlertRuleConfig struct {
	Interval string      `yaml:"interval"`
	List     []AlertRule `yaml:"list"`
	Cluster  string
}

type CollectRuleConfig struct {
	Interval string        `yaml:"interval"`
	List     []CollectRule `yaml:"list"`
}

type CollectRule struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}

type AlertRule struct {
	Name    string `yaml:"name"`
	Sha265  string `yaml:"sha265"`
	Version string `yaml:"version"`
}

func NewAlertRuleConfig(in string) (*AlertRuleConfig, error) {
	var config *AlertRuleConfig
	err := parsers.ParserConfigurationByFile(parsers.YAML, in, &config)
	return config, err
}

func NewCollectRuleConfig(in string) (*CollectRuleConfig, error) {
	var config *CollectRuleConfig
	err := parsers.ParserConfigurationByFile(parsers.YAML, in, &config)
	return config, err
}
