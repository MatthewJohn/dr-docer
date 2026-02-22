package gitlab

import (
	"fmt"
	"net/url"

	gitDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/git"
	"gopkg.in/src-d/go-git.v4/storage"
)

type GitlabConfig struct {
	Url    string
	Token  string
	Domain string
}

type GitlabClient struct {
	config  *GitlabConfig
	storage *storage.Storer
}

func NewGitlabClient(config *GitlabConfig, storage *storage.Storer) (*GitlabClient, error) {
	if config == nil {
		return nil, fmt.Errorf("NewGitlabClient: config is nil")
	}
	if config.Url == "" {
		return nil, fmt.Errorf("NewGitlabClient: Gitlab URL is nil")
	}
	if storage == nil {
		return nil, fmt.Errorf("NewGitlabClient: storage is nil")
	}
	return &GitlabClient{
		config:  config,
		storage: storage,
	}, nil
}

func (g *GitlabClient) MakeRequest(url string) {

}

func (g *GitlabClient) GetDomain() (string, error) {
	if g.config.Domain == "" {
		parsedUrl, err := url.Parse(g.config.Url)
		if err != nil {
			return "", err
		}
		g.config.Domain = parsedUrl.Host
	}
	return g.config.Domain, nil
}

func (g *GitlabClient) ConvertUrlToRepoAndPath(cloneUrl string) (string, string, error) {
	parsedUrl, err := url.Parse(cloneUrl)
	if err != nil {
		return "", "", err
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return "", "", fmt.Errorf("ConvertUrlToRepoAndPath: Unsupported scheme: %s", parsedUrl.Scheme)
	}
	parsedUrl.User = url.UserPassword("token", g.config.Token)

	// Get fragment from Gitlab URL to determine if this changes the hostname
	return "", "", nil
}

var _ gitDomain.GitRepoProvider = &GitlabClient{}
