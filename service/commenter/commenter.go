// Package commenter provides functionality to post a comment into a GitHub PR
package commenter

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// AddComment adds a comment to an issue or pull request
func AddComment(token, owner, repo, body string, issueNumber int) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	comment := &github.IssueComment{
		Body: &body,
	}
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, issueNumber, comment)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
