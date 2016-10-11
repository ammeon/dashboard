package chart

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/repo"
)

var (
	repositoryDir          string = "repository"
	repositoriesFilePath   string = "repositories.yaml"
	cachePath              string = "cache"
	localRepoPath          string = "local"
	localRepoIndexFilePath string = "index.yaml"
)

func homePath() string {
	return "/.helm"
}

func repositoryDirectory() string {
	return homePath() + "/" + repositoryDir
}

func cacheDirectory(paths ...string) string {
	fragments := append([]string{repositoryDirectory(), cachePath}, paths...)
	return filepath.Join(fragments...)
}

func cacheIndexFile(repoName string) string {
	return cacheDirectory(repoName + "-index.yaml")
}

func localRepoDirectory(paths ...string) string {
	fragments := append([]string{repositoryDirectory(), localRepoPath}, paths...)
	return filepath.Join(fragments...)
}

func repositoriesFile() string {
	return filepath.Join(repositoryDirectory(), repositoriesFilePath)
}

var (
	stableRepository    = "kubernetes-charts"
	stableRepositoryURL = "http://storage.googleapis.com/kubernetes-charts"
	customRepository    = "ammeon-charts"
	customRepositoryURL = "http://172.19.29.166:8879/charts"
)

// ensureHome checks to see if $HELM_HOME exists
//
// If $HELM_HOME does not exist, this function will create it.
func ensureHome() error {
	configDirectories := []string{homePath(), repositoryDirectory(), cacheDirectory(), localRepoDirectory()}
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			fmt.Printf("Creating %s \n", p)
			if err := os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("Could not create %s: %s", p, err)
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("%s must be a directory", p)
		}
	}

	repoFile := repositoriesFile()
	if fi, err := os.Stat(repoFile); err != nil {
		fmt.Printf("Creating %s \n", repoFile)
		r := repo.NewRepoFile()
		r.Add(&repo.Entry{
			Name:  stableRepository,
			URL:   stableRepositoryURL,
			Cache: "stable-index.yaml",
		})
		if err := r.WriteFile(repoFile, 0644); err != nil {
			return err
		}
		cif := cacheIndexFile(stableRepository)
		if err := repo.DownloadIndexFile(stableRepository, stableRepositoryURL, cif); err != nil {
			fmt.Printf("WARNING: Failed to download %s: %s (run 'helm repo update')\n", stableRepository, err)
		}

		// TODO: Remove this and add custom chart repos, from an add repo dialog
		if err := addRepo(customRepository, customRepositoryURL); err != nil {
			return err
		}

	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", repoFile)
	}
	if r, err := repo.LoadRepositoriesFile(repoFile); err == repo.ErrRepoOutOfDate {
		fmt.Println("Updating repository file format...")
		if err := r.WriteFile(repoFile, 0644); err != nil {
			return err
		}
	}

	fmt.Printf("$HELM_HOME has been configured at %s.\n", homePath())
	return nil
}
