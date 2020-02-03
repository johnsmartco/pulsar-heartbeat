package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// DefaultConfigFile - default config file
// it can be overwritten by env variable PULSAR_BEAM_CONFIG
const DefaultConfigFile = "../config/pulsar_beam.json"

// PrometheusCfg configures Premetheus set up
type PrometheusCfg struct {
	Port          string `json:"port"`
	ExposeMetrics bool   `json:"exposeMetrics"`
}

// SlackCfg is slack key
type SlackCfg struct {
	AlertURL string `json:"alertUrl"`
}

// OpsGenieCfg is opsGenie key
type OpsGenieCfg struct {
	Key             string `json:"key"`
	IntervalSeconds int    `json:"intervalSeconds"`
}

// SiteCfg configures general website
type SiteCfg struct {
	Headers         map[string]string `json:"headers"`
	URL             string            `json:"url"`
	Name            string            `json:"name"`
	IntervalSeconds int               `json:"intervalSeconds"`
	ResponseSeconds int               `json:"responseSeconds"`
	StatusCode      int               `json:"statusCode"`
	Retries         int               `json:"retries"`
}

// SitesCfg configures a list of website`
type SitesCfg struct {
	Sites []SiteCfg `json:"sites"`
}

// PulsarOpsCfg is for monitor a list of Pulsar cluster
type PulsarOpsCfg struct {
	MasterToken     string   `json:"masterToken"`
	Clusters        []string `json:"clusters"`
	IntervalSeconds int      `json:"intervalSeconds"`
}

// TopicCfg is topic configuration
type TopicCfg struct {
	LatencyBudgetMs int    `json:"latencyBudgetMs"`
	PulsarURL       string `json:"pulsarUrl"`
	TopicName       string `json:"topicName"`
}

// PulsarPerfCfg is configuration to monitor Pulsar pub sub latency
type PulsarPerfCfg struct {
	Token           string     `json:"token"`
	TrustStore      string     `json:"trustStore"`
	IntervalSeconds int        `json:"intervalSeconds"`
	TopicCfgs       []TopicCfg `json:"topicCfgs"`
}

// Configuration - this server's configuration
type Configuration struct {
	PrometheusConfig PrometheusCfg `json:"prometheusConfig"`
	SlackConfig      SlackCfg      `json:"slackConfig"`
	OpsGenieConfig   OpsGenieCfg   `json:"opsGenieConfig"`
	PulsarOpsConfig  PulsarOpsCfg  `json:"pulsarOpsConfig"`
	PulsarPerfConfig PulsarPerfCfg `json:"pulsarPerfConfig"`
	SitesConfig      SitesCfg      `json:"sitesConfig"`
}

// Config - this server's configuration instance
var Config Configuration

// ReadConfigFile reads configuration file.
func ReadConfigFile(configFile string) {

	//filename is the path to the json config file
	file, err := os.Open(configFile)
	if err != nil {
		log.Printf("failed to load configuration file %s", configFile)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		panic(err)
	}

	log.Println(Config)
}

//GetConfig returns a reference to the Configuration
func GetConfig() *Configuration {
	return &Config
}

//
type monitorFunc func()

// RunInterval runs interval
func RunInterval(fn monitorFunc, interval time.Duration) {
	go func() {
		fn()
		for {
			select {
			case <-time.Tick(interval):
				fn()
			}
		}

	}()
}
