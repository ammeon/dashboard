package chart

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/repo"
)

// RepositorySpec is a specification for a repository.
type RepositorySpec struct {
	// Name of the chart.
	RepoName string `json:"repoName"`

	// Name of the release.
	RepoUrl string `json:"repoUrl"`
}

// RepositoryListSpec is a specification for a repository.
type RepositoryListSpec struct {
	// List of repository names.
	RepoNames []string `json:"repoNames"`
}

// RepositoryListSpec is a specification for a repository.
type RepositoryChartListSpec struct {
	// List of charts.
	Charts []ChartSpec `json:"charts"`
}

// Chartspec representation view of a chart.
type ChartSpec struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	FullURL     string `json:"fullURL"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// AddRepository adds a repository.
func AddRepository(spec *RepositorySpec) error {
	return addRepo(spec.RepoName, spec.RepoUrl)
}

// GetRepositoryList get a list of repository.
func GetRepositoryList() (*RepositoryListSpec, error) {
	ensureHome()
	repoList := &RepositoryListSpec{
		RepoNames: make([]string, 0),
	}
	f, err := repo.LoadRepositoriesFile(repositoriesFile())
	if err != nil {
		return repoList, err
	}
	for repoName, _ := range f.Repositories {
		repoList.RepoNames = append(repoList.RepoNames, repoName)
	}
	return repoList, nil
}

// GetRepositoryCharts get charts in a repository.
func GetRepositoryCharts(repoName string) (*RepositoryChartListSpec, error) {
	chartList := &RepositoryChartListSpec{
		Charts: make([]ChartSpec, 0),
	}
	r, err := repo.LoadIndexFile(cacheIndexFile(repoName))
	if err != nil {
		return chartList, err
	}
	for _, c := range r.Entries {
		if c.Chartfile == nil {
			continue
		}
		icon := c.Chartfile.Icon
		if icon == "" {
			icon = "https://deis.com/assets/images/svg/helm-logo.svg"
		}
		chart := &ChartSpec{
			Name:        c.Chartfile.Name,
			Version:     c.Chartfile.Version,
			FullURL:     c.URL,
			Description: c.Chartfile.Description,
			Icon:        icon,
		}
		chartList.Charts = append(chartList.Charts, *chart)
	}
	return chartList, nil
}

func index(dir, url string) error {
	chartRepo, err := repo.LoadChartRepository(dir, url)
	if err != nil {
		return err
	}
	return chartRepo.Index()
}

func addRepo(name, url string) error {
	if err := repo.DownloadIndexFile(name, url, cacheIndexFile(name)); err != nil {
		return errors.New("Looks like " + url + " is not a valid chart repository or cannot be reached: " + err.Error())
	}

	return insertRepoLine(name, url)
}

func removeRepoLine(name string) error {
	r, err := repo.LoadRepositoriesFile(repositoriesFile())
	if err != nil {
		return err
	}

	_, ok := r.Repositories[name]
	if ok {
		delete(r.Repositories, name)
		b, err := yaml.Marshal(&r.Repositories)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(repositoriesFile(), b, 0666); err != nil {
			return err
		}
		if err := removeRepoCache(name); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("The repository, %s, does not exist in your repositories list", name)
	}

	return nil
}

func removeRepoCache(name string) error {
	if _, err := os.Stat(cacheIndexFile(name)); err == nil {
		err = os.Remove(cacheIndexFile(name))
		if err != nil {
			return err
		}
	}
	return nil
}

func insertRepoLine(name, url string) error {
	f, err := repo.LoadRepositoriesFile(repositoriesFile())
	if err != nil {
		return err
	}
	_, ok := f.Repositories[name]
	if ok {
		return fmt.Errorf("The repository name you provided (%s) already exists. Please specify a different name.", name)
	}

	if f.Repositories == nil {
		f.Repositories = make(map[string]string)
	}

	f.Repositories[name] = url

	b, _ := yaml.Marshal(&f.Repositories)
	return ioutil.WriteFile(repositoriesFile(), b, 0666)
}
