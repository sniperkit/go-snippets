package lib

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/go-github/github"
	pb "github.com/keelerh/omniscience/protos"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	reposToListPerPage = 25 // Maximum number of Github repos to return per page.
	service            = "github"
)

type GithubService struct {
	client       *github.Client
	organization string
}

func NewGithub(token string, organization string) *GithubService {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GithubService{
		client:       client,
		organization: organization,
	}
}

func (g *GithubService) Fetch(modifiedSince time.Time) ([]*pb.Document, error) {
	var allDocuments []*pb.Document

	modifiedSinceProto, err := ptypes.TimestampProto(modifiedSince)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse modified since timestamp as proto")
	}

	ctx := context.Background()
	repos, err := g.getReposInOrganization(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range repos {
		owner := r.GetOwner().GetLogin()
		readme, _, err := g.client.Repositories.GetReadme(ctx, owner, r.GetName(), nil)
		if err != nil {
			return nil, err
		}
		content, err := readme.GetContent()
		if err != nil {
			return nil, err
		}
		// Individual files do not have an ID but as we are currently only retrieving
		// the preferred README for a repository we can instead use the repo ID.
		id := strconv.FormatInt(r.GetID(), 10)
		allDocuments = append(allDocuments, &pb.Document{
			Id:           &pb.DocumentId{Id: id},
			Title:        r.GetName() + "/" + readme.GetName(),
			Description:  "",
			Service:      service,
			Content:      content,
			Url:          readme.GetHTMLURL(),
			LastModified: modifiedSinceProto,
		})

	}

	return allDocuments, nil
}

func (g *GithubService) getReposInOrganization(ctx context.Context) ([]*github.Repository, error) {
	var repos []*github.Repository
	listOpts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: reposToListPerPage,
		},
	}

	pageIdx := 0
	for {
		listOpts.ListOptions.PerPage = pageIdx
		rs, _, err := g.client.Repositories.ListByOrg(ctx, g.organization, listOpts)
		if err != nil {
			return nil, err
		}
		repos = append(repos, rs...)
		if len(rs) == reposToListPerPage {
			pageIdx += reposToListPerPage
		} else {
			break
		}
	}

	return repos, nil
}
