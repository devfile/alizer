package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/devfile/alizer/pkg/apis/model"
	recognizer "github.com/devfile/alizer/pkg/apis/recognizer"
	"github.com/devfile/alizer/pkg/schema"
	"gopkg.in/yaml.v2"
)

type StarterProject struct {
	Repo     string
	SubDir   string
	Revision string
	Remote   string
}

type ProjectReplacement struct {
	Devfile         string
	ReplacementRepo string
	Revision        string
	SubDir          string
}

type RegistryCheckJSONItem struct {
	Devfile  string
	Repo     string
	Revision string
	Registry string
	SubDir   string
}

type DevfileRegistry struct {
	Name   string
	Url    string
	Filter model.DevfileFilter
}

// getProjectReplacements: Returns any stacks that have different remotes and they cannot
// be found from the detection algorithm
func getProjectReplacements() []ProjectReplacement {
	return []ProjectReplacement{
		{
			Devfile:         "java-wildfly",
			ReplacementRepo: "https://github.com/wildfly-extras/wildfly-devfile-examples.git",
			SubDir:          "",
			Revision:        "qs",
		},
		{
			Devfile:         "java-wildfly-bootable-jar",
			ReplacementRepo: "https://github.com/wildfly-extras/wildfly-jar-maven-plugin.git",
			SubDir:          "examples/authentication",
			Revision:        "",
		},
		{
			Devfile:         "java-quarkus",
			ReplacementRepo: "https://github.com/devfile-samples/devfile-sample-code-with-quarkus",
			SubDir:          "",
			Revision:        "",
		},
		{
			Devfile:         "java-jboss-eap-xp",
			ReplacementRepo: "https://github.com/jboss-developer/jboss-eap-quickstarts",
			SubDir:          "",
			Revision:        "",
		},
		{
			Devfile:         "java-jboss-eap-xp-bootable-jar",
			ReplacementRepo: "https://github.com/jboss-developer/jboss-eap-quickstarts",
			SubDir:          "",
			Revision:        "",
		},
	}
}

// getRegistries: Fetches all registries we want to run the check
func getRegistries() []DevfileRegistry {
	return []DevfileRegistry{
		{
			Name:   "Community Registry",
			Url:    "https://registry.devfile.io",
			Filter: model.DevfileFilter{},
		},
		{
			Name:   "Product Registry",
			Url:    "https://devfile-registry.redhat.com",
			Filter: model.DevfileFilter{},
		},
	}
}

// getExcludedEntries: Returns a list of excluded stacks from the check
func getExcludedEntries() []string {
	return []string{"java-maven", "java-websphereliberty", "java-websphereliberty-gradle"}
}

func getStarterProjects(url string) ([]StarterProject, error) {
	// This value is set by the user in order to configure the registry
	var starterProjects []StarterProject
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return []StarterProject{}, err
	}
	defer func(){
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing file: %s", err)
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {

		return []StarterProject{}, fmt.Errorf("unable to fetch starter project from registry %s. code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []StarterProject{}, fmt.Errorf("unable to read body from response - error: %s", err)
	}

	// For each devfile fetched from the list, get the information from the detail page
	var devfileYaml schema.DevfileYaml
	err = yaml.Unmarshal(body, &devfileYaml)
	if err != nil {
		return []StarterProject{}, fmt.Errorf("unable to unmarshal json from response - error: %s", err)
	}

	for _, starterProject := range devfileYaml.StarterProjects {
		starterProjects = append(starterProjects, StarterProject{
			Repo:     starterProject.Git.Remotes.Origin,
			SubDir:   starterProject.SubDir,
			Revision: starterProject.Git.CheckoutFrom.Revision,
			Remote:   starterProject.Git.CheckoutFrom.Remote,
		})
	}
	return starterProjects, nil
}

func appendIfMissing(slice []RegistryCheckJSONItem, r RegistryCheckJSONItem) []RegistryCheckJSONItem {
	for _, ele := range slice {
		if ele == r {
			return slice
		}
	}
	return append(slice, r)
}

func isIgnoredEntry(name string) bool {
	excludedEntries := getExcludedEntries()
	for _, excludedEntry := range excludedEntries {
		if excludedEntry == name {
			return true
		}
	}
	return false
}

func main() {
	var registryEntriesList []RegistryCheckJSONItem
	devfileRegistries := getRegistries()
	projectReplacements := getProjectReplacements()

	for _, registry := range devfileRegistries {
		tmpDevfileTypes, err := recognizer.DownloadDevfileTypesFromRegistry(registry.Url, registry.Filter)
		if err != nil {
			continue
		}
		for _, devfileType := range tmpDevfileTypes {
			// Continnue if the entry is inside the excluded list
			if isIgnoredEntry(devfileType.Name) {
				continue
			}
			// For every stack get its detailed devfile.yaml
			starterProjects, err := getStarterProjects(fmt.Sprintf("%s/devfiles/%s", registry.Url, devfileType.Name))
			if err != nil {
				continue
			}
			if len(starterProjects) == 0 {
				continue
			}

			starterProject := starterProjects[0]
			repo := starterProject.Repo
			revision := starterProject.Revision
			subdir := starterProject.SubDir
			if repo == "" {
				for _, projectReplacement := range projectReplacements {
					if devfileType.Name == projectReplacement.Devfile {
						repo = projectReplacement.ReplacementRepo
						revision = projectReplacement.Revision
						subdir = projectReplacement.SubDir
					}
				}
			}
			registryEntryItem := RegistryCheckJSONItem{
				Devfile:  devfileType.Name,
				Repo:     repo,
				Registry: registry.Url,
				Revision: revision,
				SubDir:   subdir,
			}
			registryEntriesList = appendIfMissing(registryEntriesList, registryEntryItem)
		}
		jsonBytes, err := json.Marshal(registryEntriesList)
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	}
}
