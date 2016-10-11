package chart

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	for _, r := range f.Repositories {
		repoList.RepoNames = append(repoList.RepoNames, r.Name)
	}
	return repoList, nil
}

// GetRepositoryCharts get charts in a repository.
func GetRepositoryCharts(repoName string) (*RepositoryChartListSpec, error) {
	chartList := &RepositoryChartListSpec{
		Charts: make([]ChartSpec, 0),
	}
	i, err := repo.LoadIndexFile(cacheIndexFile(repoName))
	if err != nil {
		return chartList, err
	}
	for _, e := range i.Entries {
		for _, c := range e {
			if c == nil {
				continue
			}
			icon := c.Icon
			if icon == "" {
				icon = "https://deis.com/assets/images/svg/helm-logo.svg"
			}
			desc := c.Description
			if len(desc) > 45 {
				desc = desc[0:41] + "..."
			}
			chart := &ChartSpec{
				Name:        c.Name,
				Version:     c.Version,
				FullURL:     c.URLs[0],
				Description: desc,
				Icon:        icon,
			}
			chartList.Charts = append(chartList.Charts, *chart)
		}
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
	cif := cacheIndexFile(name)
	if err := repo.DownloadIndexFile(name, url, cif); err != nil {
		return fmt.Errorf("Looks like %q is not a valid chart repository or cannot be reached: %s", url, err.Error())
	}

	return insertRepoLine(name, url)
}

func removeRepoLine(out io.Writer, name string) error {
	repoFile := repositoriesFile()
	r, err := repo.LoadRepositoriesFile(repoFile)
	if err != nil {
		return err
	}

	if !r.Remove(name) {
		return fmt.Errorf("no repo named %q found", name)
	}
	if err := r.WriteFile(repoFile, 0644); err != nil {
		return err
	}

	if err := removeRepoCache(name); err != nil {
		return err
	}

	fmt.Fprintf(out, "%q has been removed from your repositories", name)

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
	cif := cacheIndexFile(name)
	f, err := repo.LoadRepositoriesFile(repositoriesFile())
	if err != nil {
		return err
	}

	if f.Has(name) {
		return fmt.Errorf("The repository name you provided (%s) already exists. Please specify a different name.", name)
	}
	f.Add(&repo.Entry{
		Name:  name,
		URL:   url,
		Cache: filepath.Base(cif),
	})
	return f.WriteFile(repositoriesFile(), 0644)
}
