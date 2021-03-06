package config

import (
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	// DefaultInterval defines default interval in case it is missing in yaml.
	DefaultInterval = 2 * time.Minute
)

type Instance struct {
	Instance     string        `yaml:"instance"`
	Region       string        `yaml:"region"`
	Interval     time.Duration `yaml:"interval"`
	AwsAccessKey string        `yaml:"aws_access_key"`
	AwsSecretKey string        `yaml:"aws_secret_key"`
}

type Config struct {
	Instances []Instance `yaml:"instances"`
}

type Settings struct {
	config Config
	sync.RWMutex
	// AfterLoad is run after every Load request but before releasing Mutex
	AfterLoad func(config Config) error
}

func (s *Settings) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	if err := yaml.Unmarshal(content, &s.config); err != nil {
		return err
	}
	for i := range s.config.Instances {
		if s.config.Instances[i].Interval.Nanoseconds() == 0 {
			s.config.Instances[i].Interval = DefaultInterval
		}
	}

	if s.AfterLoad != nil {
		return s.AfterLoad(s.config)
	}

	return nil
}

func (s *Settings) Config() Config {
	s.RLock()
	defer s.RUnlock()
	return s.config
}
