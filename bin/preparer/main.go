package main

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/square/p2/pkg/logging"
	"github.com/square/p2/pkg/preparer"
	"gopkg.in/yaml.v2"
)

type PreparerConfig struct {
	NodeName       string `yaml:"node_name"`
	ConsulAddress  string `yaml:"consul_address"`
	HooksDirectory string `yaml:"hooks_directory"`
}

func main() {
	logger := logging.NewLogger(logrus.Fields{
		"app": "preparer",
	})
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logger.NoFields().Println("No CONFIG_PATH variable was given")
		os.Exit(1)
		return
	}
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"inner_err": err,
		}).Println("Could not read the config file")
		os.Exit(1)
		return
	}
	preparerConfig := PreparerConfig{}
	err = yaml.Unmarshal(configBytes, &preparerConfig)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"inner_err": err,
		}).Println("The config file was malformatted")
		os.Exit(1)
		return
	}
	if preparerConfig.NodeName == "" {
		logger.NoFields().Println("`node_name` was not set in the file at CONFIG_PATH")
		os.Exit(1)
		return
	}
	if preparerConfig.ConsulAddress == "" {
		preparerConfig.ConsulAddress = "127.0.0.1:8500"
	}
	if preparerConfig.HooksDirectory == "" {
		preparerConfig.HooksDirectory = "/usr/local/p2hooks.d"
	}
	logFile, err := os.OpenFile("/tmp/platypus", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logFile.Close()

	logger.WithFields(logrus.Fields{
		"starting":  true,
		"node_name": preparerConfig.NodeName,
		"consul":    preparerConfig.ConsulAddress,
		"hooks_dir": preparerConfig.HooksDirectory,
	}).Println("Preparer started successfully") // change to logrus message

	preparer.WatchForPodManifestsForNode(preparerConfig.NodeName, preparerConfig.ConsulAddress, preparerConfig.HooksDirectory, logger)
	logger.WithFields(logrus.Fields{
		"stopping": true,
	}).Println("Terminating")
}