package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Config provides a structure for running a markdown wiki site
type Config struct {
	SitePath        string `json:"sitePath,omitempty"`
	TmpDir          string `json:"tmpDir,omitempty"`
	ApplicationName string `json:"applicationName,omitempty"`
	ConfigFilePath  string `json:"configFilePath,omitempty"`
	GitPath         string `json:"gitPath,omitempty"`
	CanEdit         bool   `json:"canEdit"`
}

func (c *Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panic("Fatal error serialziing a Config object", err)
	}
	return string(b)
}

// Merge overwrites config values from the new config object
func (c *Config) Merge(newConfig *Config) {
	if newConfig.SitePath != "" {
		c.SitePath = newConfig.SitePath
	}

	if newConfig.TmpDir != "" {
		c.TmpDir = newConfig.TmpDir
	}

	if newConfig.ApplicationName != "" {
		c.ApplicationName = newConfig.ApplicationName
	}

	if newConfig.ConfigFilePath != "" {
		c.ConfigFilePath = newConfig.ConfigFilePath
	}

	if newConfig.GitPath != "" {
		c.GitPath = newConfig.GitPath
	}
	c.CanEdit = newConfig.CanEdit
}

type option func(m *Config)

func FromConfig(newConfig *Config) option {
	return func(c *Config) {
		c.Merge(newConfig)
	}
}

func SitePath(s string) option {
	return func(c *Config) {
		c.SitePath = s
	}
}

func TmpDir(s string) option {
	return func(c *Config) {
		c.TmpDir = s
	}
}

func ApplicationName(s string) option {
	return func(c *Config) {
		c.ApplicationName = s
	}
}

func setupConfig(opts ...option) *Config {
	c := &Config{}
	c.TmpDir = "/tmp"
	c.SitePath = "."
	c.ApplicationName = "Your first yanta site"
	c.ConfigFilePath = "./yanta.json"
	c.GitPath = c.SitePath
	c.CanEdit = true
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DecodeConfig deserializes a byte array to a site
func DecodeConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := unmarshalConfig(b)
	if err != nil {
		return nil, err
	}
	c.ConfigFilePath = path

	return c, nil
}

func unmarshalConfig(b []byte) (*Config, error) {
	c := new(Config)
	err := json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ServeConfig is an an http handler for serving the config struct
func ServeConfig(s *Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, s.Config)
	}
}

// UpdateConfig is a handler for updating the config
func UpdateConfig(s *Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Config.CanEdit == false {
			log.Println("not allowed to edit this site")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading body", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cfg, err := unmarshalConfig(b)
		if err != nil {
			log.Println("error unmarshaling", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(cfg.ConfigFilePath, b, 0644)
		if err != nil {
			log.Println("could not write config file", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Config.Merge(cfg)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(b))
	}
}
