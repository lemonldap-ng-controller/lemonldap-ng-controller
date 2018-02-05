/*
Copyright 2018 Mathieu Parent <math.parent@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem"
)

var validConfigurationName = regexp.MustCompile(`^lmConf-(\d+)\.js$`)

// Config defines a LemonLDAP::NG configuration loader
type Config struct {
	sync.RWMutex

	fs        filesystem.FileSystem
	configDir string
	cfgNum    int
	overrides map[string]interface{}
	vhosts    map[string]*VHost
	dirty     bool
}

// NewConfig creates a new LemonLDAP::NG configuration loader
func NewConfig(fs filesystem.FileSystem, configDir string) *Config {
	return &Config{
		fs:        fs,
		configDir: configDir,
		cfgNum:    1,
		overrides: make(map[string]interface{}),
		vhosts:    make(map[string]*VHost),
	}
}

// First returns the first configuration file name and number
func (c *Config) First() (string, int, error) {
	return fmt.Sprintf("lmConf-%d.js", 1), 1, nil
}

// Last returns the current configuration file name and number
func (c *Config) Last() (string, int, error) {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("lmConf-%d.js", c.cfgNum), c.cfgNum, nil
}

// Next returns the following configuration file name and number
func (c *Config) Next() (string, int, error) {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("lmConf-%d.js", c.cfgNum+1), c.cfgNum + 1, nil
}

// Load loads a specific LemonLDAP::NG configuration
func (c *Config) Load(configName string) (map[string]interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	return c.loadNoLock(configName)
}

// Load loads a specific LemonLDAP::NG configuration
func (c *Config) loadNoLock(configName string) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	path := c.configDir + "/" + configName
	content, err := c.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to read LemonLDAP::NG configuration file %s: %s", path, err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse LemonLDAP::NG configuration file %s: %s", path, err)
	}
	c.dirty = true
	return conf, nil
}

// LoadFirst loads the first LemonLDAP::NG configuration
func (c *Config) LoadFirst() (map[string]interface{}, error) {
	firstConfigName, _, _ := c.First()
	return c.Load(firstConfigName)
}

// Save saves the current LemonLDAP::NG configuration as next
func (c *Config) Save() error {
	c.Lock()
	defer c.Unlock()
	if !c.dirty {
		return nil
	}
	nextConfigNum := c.cfgNum + 1
	nextConfigName := fmt.Sprintf("lmConf-%d.js", c.cfgNum+1)
	path := c.configDir + "/" + nextConfigName
	conf, err := c.loadNoLock("lmConf-1.js")
	if err != nil {
		return err
	}
	for overridek, overridev := range c.overrides {
		conf[overridek] = overridev
	}
	conf["cfgAuthor"] = "lemonldap-ng-controller"
	conf["cfgNum"] = nextConfigNum
	allExportedHeaders, ok := conf["exportedHeaders"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("exportedHeaders should be a map, got %T", conf["exportedHeaders"])
	}
	allLocationRules, ok := conf["locationRules"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("locationRules should be a map, got %T", conf["locationRules"])
	}
	for serverName, vhost := range c.vhosts {
		allExportedHeaders[serverName] = vhost.ExportedHeaders
		allLocationRules[serverName] = vhost.LocationRules
	}
	content, err := json.MarshalIndent(conf, "", "   ")
	if err != nil {
		return fmt.Errorf("Unable to encode LemonLDAP::NG configuration file %s: %s", nextConfigName, err)
	}
	err = c.fs.WriteFile(path, content, 0777)
	if err != nil {
		return fmt.Errorf("Unable to write LemonLDAP::NG configuration file %s: %s", path, err)
	}
	c.cfgNum++
	c.ReloadLemonLDAPNG()
	c.dirty = false
	return nil
}

// stringifyKeysMapValue recurses into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedence mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
// and https://github.com/gohugoio/hugo/pull/4138
func stringifyYAMLMapKeys(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i] = stringifyYAMLMapKeys(v)
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[fmt.Sprintf("%v", k)] = stringifyYAMLMapKeys(v)
		}
		return res
	default:
		return in
	}
}

// SetOverrides creates several new LemonLDAP::NG virtual hosts
func (c *Config) SetOverrides(overrides map[string]interface{}) error {
	c.Lock()
	defer c.Unlock()
	m := map[string]interface{}{}
	for k, v := range overrides {
		m[k] = stringifyYAMLMapKeys(v)
	}
	c.overrides = m
	c.dirty = true
	return nil
}

// AddVhosts creates several new LemonLDAP::NG virtual hosts
func (c *Config) AddVhosts(vhosts map[string]*VHost) error {
	c.Lock()
	defer c.Unlock()
	for _, vhost := range vhosts {
		c.vhosts[vhost.ServerName] = &VHost{
			vhost.ServerName,
			vhost.LocationRules,
			vhost.ExportedHeaders,
		}
	}
	c.dirty = true
	return nil
}

// DeleteVhosts deletes several LemonLDAP::NG virtual hosts
func (c *Config) DeleteVhosts(vhosts map[string]*VHost) error {
	c.Lock()
	defer c.Unlock()
	for _, vhost := range vhosts {
		delete(c.vhosts, vhost.ServerName)
	}
	c.dirty = true
	return nil
}
