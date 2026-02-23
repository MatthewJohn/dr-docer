package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type GitRepoProvider interface {
	GetDomain() (string, error)
	ConvertUrlToRepoAndPath(url string) (string, string, error)
}

type Repo struct {
	Repository *git.Repository
	Storage    *memory.Storage
	Filesystem billy.Filesystem
}

type GitService struct {
	gitProviders []GitRepoProvider
	repos        map[string]Repo
}

func NewGitService() (*GitService, error) {
	return &GitService{}, nil
}

func (g *GitService) RegisterGitRepoProvider(repo GitRepoProvider) error {
	domain, err := repo.GetDomain()
	if err != nil {
		return err
	}
	for _, gitRepoProvider := range g.gitProviders {
		if gitRepoProviderDomain, err := gitRepoProvider.GetDomain(); err == nil && gitRepoProviderDomain == domain {
			return fmt.Errorf("RegisterGitRepoProvider: provider already registered with domain: %s", domain)
		}
	}
	g.gitProviders = append(g.gitProviders, repo)
	return nil
}

func (g *GitService) GetRepoProviderForUrl(gitUrl string) (*GitRepoProvider, error) {
	urlParts, err := url.Parse(gitUrl)
	if err != nil {
		return nil, err
	}
	domain := urlParts.Host
	for _, gitRepoProvider := range g.gitProviders {
		if gitRepoProviderDomain, err := gitRepoProvider.GetDomain(); err == nil && gitRepoProviderDomain == domain {
			return &gitRepoProvider, nil
		}
	}
	return nil, nil
}

func (g *GitService) generateRepoHash(repoUrl string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(repoUrl))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (g *GitService) CloneRepository(repoUrl string) (*git.Repository, error) {
	repoHash, err := g.generateRepoHash(repoUrl)
	if err != nil {
		return nil, err
	}
	if _, ok := g.repos[repoHash]; ok {
		return g.repos[repoHash].Repository, nil
	}

	storage := memory.NewStorage()
	filesystem := memfs.New()

	clonedRepo, err := git.Clone(storage, filesystem, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		return nil, err
	}

	repo := Repo{
		Filesystem: filesystem,
		Storage:    storage,
		Repository: clonedRepo,
	}
	g.repos[repoHash] = repo
	return repo.Repository, nil

}
