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
	"strconv"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/filesystem"
)

var validConfigurationName = regexp.MustCompile(`^lmConf-(\d+)\.js$`)

// Config defines a LemonLDAP::NG configuration loader
type Config struct {
	fs        filesystem.FileSystem
	configDir string
	overrides map[string]interface{}
	vhosts    map[string]*VHost
	dirty     bool
}

// NewConfig creates a new LemonLDAP::NG configuration loader
func NewConfig(fs filesystem.FileSystem, configDir string) *Config {
	return &Config{
		fs:        fs,
		configDir: configDir,
		overrides: make(map[string]interface{}),
		vhosts:    make(map[string]*VHost),
	}
}

// first returns the first configuration file name and number
func (c *Config) first() (string, int, error) {
	return "lmConf-1.js", 1, nil
}

// last returns the current configuration file name and number
func (c *Config) last() (string, int, error) {
	f, err := c.fs.Open(c.configDir)
	if err != nil {
		return "", 0, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return "", 0, err
	}
	lastConfigNum := 0
	ret := ""
	for _, conf := range list {
		match := validConfigurationName.FindStringSubmatch(conf.Name())
		if len(match) > 0 {
			configNum, _ := strconv.Atoi(match[1])
			if configNum > lastConfigNum {
				lastConfigNum = configNum
				ret = conf.Name()
			}
		}
	}
	if lastConfigNum == 0 {
		return "", 0, fmt.Errorf("No LemonLDAP::NG configuration file found in %s", c.configDir)
	}
	return ret, lastConfigNum, nil
}

// next returns the following configuration file name and number
func (c *Config) next() (string, int, error) {
	_, lastConfigNum, err := c.last()
	if err != nil {
		return "", 0, err
	}
	nextConfigNum := lastConfigNum + 1
	nextConfigName := fmt.Sprintf("lmConf-%d.js", nextConfigNum)
	return nextConfigName, nextConfigNum, nil
}

// Load loads a specific LemonLDAP::NG configuration
func (c *Config) Load(configName string) (map[string]interface{}, error) {
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
	firstConfigName, _, _ := c.first()
	return c.Load(firstConfigName)
}

// Save saves the current LemonLDAP::NG configuration as next
func (c *Config) Save() error {
	if !c.dirty {
		return nil
	}
	nextConfigName, nextConfigNum, err := c.next()
	if err != nil {
		return err
	}
	path := c.configDir + "/" + nextConfigName
	conf, err := c.LoadFirst()
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
	c.dirty = false
	return nil
}

// SetOverrides creates several new LemonLDAP::NG virtual hosts
func (c *Config) SetOverrides(overrides map[string]interface{}) error {
	c.overrides = overrides
	c.dirty = true
	return nil
}

// AddVhosts creates several new LemonLDAP::NG virtual hosts
func (c *Config) AddVhosts(vhosts map[string]*VHost) error {
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
	for _, vhost := range vhosts {
		delete(c.vhosts, vhost.ServerName)
	}
	c.dirty = true
	return nil
}
