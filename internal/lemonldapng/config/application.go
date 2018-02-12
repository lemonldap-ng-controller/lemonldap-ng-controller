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

// Application defines a LemonLDAP::NG application
type Application struct {
	Category    string
	Name        string
	Description string
	Logo        string
	Display     string
	URI         string
}

// NewApplication creates a new LemonLDAP::NG application from annotations
func NewApplication(vhost *VHost, annotations map[string]string, prefix string) *Application {
	if vhost == nil {
		return nil
	}
	category, ok := annotations[prefix+"/application-category"]
	if !ok {
		return nil
	}
	name, ok := annotations[prefix+"/application-name"]
	if !ok {
		return nil
	}
	description, ok := annotations[prefix+"/application-description"]
	if !ok {
		description = name
	}
	logo, ok := annotations[prefix+"/application-logo"]
	if !ok {
		logo = "gear.png"
	}
	display, ok := annotations[prefix+"/application-display"]
	if !ok {
		display = "auto"
	}
	uri, ok := annotations[prefix+"/application-uri"]
	if !ok {
		uri = "https://" + vhost.ServerName + "/" // FIXME http/https
	}
	return &Application{
		Category:    category,
		Name:        name,
		Description: description,
		Logo:        logo,
		Display:     display,
		URI:         uri,
	}
}

// Path returns the application path in the menu
func (a *Application) Path() string {
	return a.Category + "/" + a.Name
}
