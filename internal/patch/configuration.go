package patch

import (
	"errors"
	"fmt"
	_ "github.com/prometheus-community/pro-bing"
	probing "github.com/prometheus-community/pro-bing"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path/filepath"
)

var (
	errFailedBuildConfiguration = func(err error) error {
		return fmt.Errorf("failed to build configuration: %w", err)
	}
	errConfigFilePathNotSet  = errors.New("config file path not set")
	errFailedParseConfigFile = func(err error) error {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
)

type Configuration struct {
	AppName string              `yaml:"app_name"`
	Host    string              `yaml:"host"`
	Source  string              `yaml:"source"`
	Suffix  string              `yaml:"suffix"`
	Differ  DifferConfiguration `yaml:"differ"`
}

type DifferConfiguration struct {
	VerboseLevel int `yaml:"verbose_level"`
	Timeout      int `yaml:"timeout"`
}

func (c Configuration) SourceInputPath() string {
	return filepath.Join(c.Source, c.AppName, "in")
}

func (c Configuration) SourceOutputPath() string {
	return filepath.Join(c.Source, c.AppName, "out")
}

func NewConfiguration(profile string) (Configuration, error) {
	var result Configuration

	var path = os.Getenv("MEVPATCH_CONFIG_PATH")
	if path == "" {
		return Configuration{}, errFailedBuildConfiguration(errConfigFilePathNotSet)
	}

	file, err := os.ReadFile(filepath.Join(path, fmt.Sprintf("%s.yaml", profile)))
	if err != nil {
		return Configuration{}, errFailedBuildConfiguration(err)
	}

	if err := yaml.Unmarshal(file, &result); err != nil {
		return Configuration{}, errFailedBuildConfiguration(errFailedParseConfigFile(err))
	}

	return result, nil

}

var (
	errConfigurationTestFailed = func(err error) error {
		return fmt.Errorf("configuration test failed: %w", err)
	}
	errFailedToVerifySource = func(err error) error {
		return fmt.Errorf("failed to verify source: %w", err)
	}
	errFailedToParseHostURL = func(url string, err error) error {
		return fmt.Errorf("failed to parse host url %s: %w", url, err)
	}
	errFailedToPingHostURL = func(url string, err error) error {
		return fmt.Errorf("failed to ping host url %s: %w", url, err)
	}
)

func (c Configuration) Test() error {

	_, err := os.Stat(c.Source)
	if err != nil {
		return errConfigurationTestFailed(errFailedToVerifySource(err))
	}

	host, err := url.Parse(c.Host)
	if err != nil {
		return errConfigurationTestFailed(errFailedToParseHostURL(c.Host, err))
	}

	var actual = host.Hostname()

	pinger, err := probing.NewPinger(actual)
	pinger.SetPrivileged(true)
	if err != nil {
		return errConfigurationTestFailed(errFailedToPingHostURL(actual, err))
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return errConfigurationTestFailed(errFailedToPingHostURL(actual, err))
	}

	return nil

}
