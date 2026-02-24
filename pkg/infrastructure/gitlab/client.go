package gitlab

import (
	"fmt"
	"net/url"
	"strings"

	gitDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/git"
)

type GitlabConfig struct {
	Url   string
	Token string
}

type GitlabClient struct {
	config *GitlabConfig
	domain string
}

func NewGitlabClient(config *GitlabConfig) (*GitlabClient, error) {
	if config == nil {
		return nil, fmt.Errorf("NewGitlabClient: config is nil")
	}
	if config.Url == "" {
		return nil, fmt.Errorf("NewGitlabClient: Gitlab URL is nil")
	}
	return &GitlabClient{
		config: config,
	}, nil
}

func (g *GitlabClient) MakeRequest(url string) {

}

func (g *GitlabClient) GetDomain() (string, error) {
	if g.domain == "" {
		parsedUrl, err := url.Parse(g.config.Url)
		if err != nil {
			return "", err
		}
		g.domain = parsedUrl.Host
	}
	return g.domain, nil
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
	// @TODO handle this case
	if strings.Contains(parsedUrl.Path, "/-/") {
		return "", "", fmt.Errorf("URL appears to be a sub-path of a repo. Not yet supported.")
	}
	return parsedUrl.String(), "", nil
}

var _ gitDomain.GitRepoProvider = &GitlabClient{}
