/**
 * Copyright 2016 ECS Team, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package command

import (
	"os"
	"strconv"

	"github.com/bradfitz/slice"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cfcli-plugin-utils/plugin/io"
	pluginversion "github.com/ecsteam/cfcli-plugin-utils/plugin/version"
	"github.com/krujos/cfcurl"
)

type orgSpaceInfo struct {
	OrgName   string
	SpaceName string
}

type appLocator struct {
	orgSpaceInfo
	Name        string
	Memory      string
	Disk        string
	Instances   string
}

var version = "0.0.1"

// ListAllAppsPlugin - the main struct
type ListAllAppsPlugin struct {
	UI io.UI
}

// New - create new plugin with stdin and stdout
func New() *ListAllAppsPlugin {
	ui := io.UI{
		Input:  os.Stdin,
		Output: io.Writer,
	}

	return NewPlugin(ui)
}

// NewPlugin - create new plugin with specified io
func NewPlugin(ui io.UI) *ListAllAppsPlugin {
	return &ListAllAppsPlugin{
		UI: ui,
	}
}

// Start - start
func (cmd *ListAllAppsPlugin) Start() {
	plugin.Start(cmd)
}

// GetMetadata - get metadata
func (cmd *ListAllAppsPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "list-all-apps",
		Version: pluginversion.GetVersionType(version),
		Commands: []plugin.Command{
			{
				Name:     "list-all-apps",
				HelpText: "List all apps in the foundation, sorted by Org and Space",
				UsageDetails: plugin.Usage{
					Usage: "cf list-all-apps",
				},
			},
		},
	}
}

// Run -
func (cmd *ListAllAppsPlugin) Run(cli plugin.CliConnection, args []string) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
		}
	}()

	if args[0] != cmd.GetMetadata().Commands[0].Name {
		return
	}

	cmd.UI.Say("List all apps in the foundation...\n")

	apps, err := getApps(cli)

	if err != nil {
		cmd.UI.Failed("Error completing request: %v", err)
		return
	}

	cmd.UI.Ok()
	cmd.UI.Say("")

	if len(apps) == 0 {
		cmd.UI.Say("No apps found")
		return
	}

	table := cmd.UI.Table([]string{"application", "org", "space", "instances", "memory", "disk"})
	for _, app := range apps {
		table.Add(app.Name, app.OrgName, app.SpaceName, app.Instances, app.Memory, app.Disk)
	}

	table.Print()
	return
}

func getApps(cli plugin.CliConnection) (apps []appLocator, err error) {
	apps = make([]appLocator, 0, 5)
	orgSpaceMap := make(map[string]*orgSpaceInfo)

	var json map[string]interface{}

	var nextURL interface{}
	nextURL = "/v2/apps"

	for nextURL != nil {
		json, err = cfcurl.Curl(cli, nextURL.(string))
		if err != nil {
			return
		}

		appsResources := toJSONArray(json["resources"])
		for _, appIntf := range appsResources {
			app := toJSONObject(appIntf)
			appEntity := toJSONObject(app["entity"])
			appName := appEntity["name"].(string)
			appMemory := strconv.FormatFloat(appEntity["memory"].(float64), 'f', -1, 64)
			appDisk := strconv.FormatFloat(appEntity["disk_quota"].(float64), 'f', -1, 64)
			appInstances := strconv.FormatFloat(appEntity["instances"].(float64), 'f', -1, 64)

			appSpaceURL := appEntity["space_url"].(string)

			if orgSpaceMap[appSpaceURL] == nil {
				orgSpaceMap[appSpaceURL], err = getOrgSpaceInfo(cli, appSpaceURL)
				if err != nil {
					return
				}
			}

			info := orgSpaceMap[appSpaceURL]

			appInfo := appLocator{
				orgSpaceInfo: *info,
				Name:         appName,
				Memory:       appMemory,
				Disk:         appDisk,
				Instances:    appInstances,
 			}

			apps = append(apps, appInfo)
		}

		nextURL = json["next_url"]
	}

	slice.Sort(apps, func(i, j int) bool {
		locator1, locator2 := apps[i], apps[j]

		if locator1.OrgName < locator2.OrgName {
			return true
		} else if locator1.OrgName > locator2.OrgName {
			return false
		}

		if locator1.SpaceName < locator2.SpaceName {
			return true
		} else if locator1.SpaceName > locator2.SpaceName {
			return false
		}

		if locator1.Name <= locator2.Name {
			return true
		}

		return false
	})

	return
}

func getOrgSpaceInfo(cli plugin.CliConnection, spaceURL string) (info *orgSpaceInfo, err error) {
	json, err := cfcurl.Curl(cli, spaceURL)
	if err != nil {
		return
	}

	info = new(orgSpaceInfo)
	entity := toJSONObject(json["entity"])
	info.SpaceName = entity["name"].(string)

	json, err = cfcurl.Curl(cli, entity["organization_url"].(string))
	if err != nil {
		info = nil
		return
	}

	entity = toJSONObject(json["entity"])
	info.OrgName = entity["name"].(string)

	return
}

func toJSONArray(obj interface{}) []interface{} {
	return obj.([]interface{})
}

func toJSONObject(obj interface{}) map[string]interface{} {
	return obj.(map[string]interface{})
}
