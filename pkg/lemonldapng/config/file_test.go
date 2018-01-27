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
	"flag"
	"regexp"
	"testing"

	fakefs "github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/filesystem/fake"
)

func TestAddDeleteVhosts(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	fs := fakefs.NewFileSystem()
	config := NewConfig(fs, "/var/lib/lemonldap-ng/conf")
	vhosts := make(map[string]*VHost)
	vhosts = map[string]*VHost{
		"test42.example.org": NewVHost("test42.example.org", DefaultLocationRules, DefaultExportedHeaders),
	}

	config.AddVhosts(vhosts)
	errSave2 := config.Save()
	if errSave2 != nil {
		t.Errorf("%s", errSave2)
	}
	lmConf2, err2 := fs.ReadFile("/var/lib/lemonldap-ng/conf/lmConf-2.js")
	if err2 != nil {
		t.Errorf("%s", err2)
	}
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile("\"cfgAuthor\": \"lemonldap-ng-controller\","),
		regexp.MustCompile("\"cfgNum\": 2,"),
		regexp.MustCompile(`"exportedHeaders": {\s*"test42.example.org": {\s*"Auth-User": "\$uid"\s*}\s*},`),
		regexp.MustCompile(`"locationRules": {\s*"test42.example.org": {\s*"default": "accept"\s*}\s*}\s*}`),
	} {
		if !re.Match(lmConf2) {
			t.Errorf("lmConf-2.js to match %s\n%s", re, lmConf2)
		}
	}

	config.DeleteVhosts(vhosts)
	errSave3 := config.Save()
	if errSave3 != nil {
		t.Errorf("%s", errSave3)
	}
	lmConf3, err3 := fs.ReadFile("/var/lib/lemonldap-ng/conf/lmConf-3.js")
	if err3 != nil {
		t.Errorf("%s", err3)
	}
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile("\"cfgAuthor\": \"lemonldap-ng-controller\","),
		regexp.MustCompile("\"cfgNum\": 3,"),
		regexp.MustCompile(`"exportedHeaders": {},`),
		regexp.MustCompile(`"locationRules": {}`),
	} {
		if !re.Match(lmConf3) {
			t.Errorf("lmConf-3.js to match %s\n%s", re, lmConf3)
		}
	}
}

func TestOverrides(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	fs := fakefs.NewFileSystem()
	config := NewConfig(fs, "/var/lib/lemonldap-ng/conf")
	overrides := make(map[string]interface{})
	overrides = map[string]interface{}{
		"domain": "example.org",
		"exportedHeaders": map[interface{}]interface{}{
			"foo.example.org": map[interface{}]interface{}{
				"CAS-User": "$uid",
			},
		},
		"locationRules": map[interface{}]interface{}{
			"foo.example.org": map[interface{}]interface{}{
				"default": "$uid eq 'dwho'",
			},
		},
		"arrayOptions": []interface{}{
			"value",
		},
	}

	config.SetOverrides(overrides)
	errSave2 := config.Save()
	if errSave2 != nil {
		t.Errorf("%s", errSave2)
	}
	lmConf2, err2 := fs.ReadFile("/var/lib/lemonldap-ng/conf/lmConf-2.js")
	if err2 != nil {
		t.Errorf("%s", err2)
	}
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile("\"cfgAuthor\": \"lemonldap-ng-controller\","),
		regexp.MustCompile("\"cfgNum\": 2,"),
		regexp.MustCompile(`"domain": "example.org",`),
	} {
		if !re.Match(lmConf2) {
			t.Errorf("lmConf-2.js to match %s\n%s", re, lmConf2)
		}
	}
}

func TestNonExistentConfigDir(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	fs := fakefs.NewFileSystem()
	config := NewConfig(fs, "/nonexistent")

	errSave1 := config.Save() // dirty == false
	if errSave1 != nil {
		t.Errorf("%s", errSave1)
	}

	config.dirty = true
	errSave := config.Save()
	if errSave == nil || errSave.Error() != "Unable to read LemonLDAP::NG configuration file /nonexistent/lmConf-1.js: open /nonexistent/lmConf-1.js: No such file or directory" {
		t.Errorf("Expected 'Unable to read LemonLDAP::NG configuration file /nonexistent/lmConf-1.js: open /nonexistent/lmConf-1.js: No such file or directory', got '%q'", errSave)
	}
}

func TestEmptyConfigDir(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	fs := fakefs.NewFileSystem()
	fs.Mkdir("/empty", 0755)
	config := NewConfig(fs, "/empty")
	_, errLoad := config.LoadFirst()
	if errLoad == nil || errLoad.Error() != "Unable to read LemonLDAP::NG configuration file /empty/lmConf-1.js: open /empty/lmConf-1.js: No such file or directory" {
		t.Errorf("Unable to read LemonLDAP::NG configuration file /empty/lmConf-1.js: open /empty/lmConf-1.js: No such file or directory', got %q", errLoad)
	}
	config.dirty = true
	errSave := config.Save()
	if errSave == nil || errSave.Error() != "Unable to read LemonLDAP::NG configuration file /empty/lmConf-1.js: open /empty/lmConf-1.js: No such file or directory" {
		t.Errorf("Expected 'Unable to read LemonLDAP::NG configuration file /empty/lmConf-1.js: open /empty/lmConf-1.js: No such file or directory', got '%q'", errSave)
	}
}

func TestInvalidLocationRules(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	fs := fakefs.NewFileSystem()
	config := NewConfig(fs, "/var/lib/lemonldap-ng/conf")

	errWrite := fs.WriteFile("/var/lib/lemonldap-ng/conf/lmConf-1.js", []byte("{}"), 0755)
	if errWrite != nil {
		t.Errorf("%s", errWrite)
	}

	config.dirty = true
	errSave2 := config.Save()
	if errSave2 == nil || errSave2.Error() != "exportedHeaders should be a map, got <nil>" {
		t.Errorf("Expected 'exportedHeaders should be a map, got <nil>', got %q", errSave2)
	}
}
