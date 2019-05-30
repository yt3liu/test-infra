/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// pullrequest.go provides generic functions related to PullRequest

package ghutil

import (
	"fmt"

	"github.com/google/go-github/github"
)

const (
	// PullRequestOpenState is the state of open PullRequest
	PullRequestOpenState PullRequestState = "open"
	// PullRequestCloseState is the state of closed PullRequest
	PullRequestCloseState PullRequestState = "closed"
	// PullRequestAllState is the state for all, useful when querying PullRequest
	PullRequestAllState PullRequestState = "all"
)

// PullRequestState represents different states of PullRequest
type PullRequestState string

// ListPullRequests lists pull requests within given repo, filters by head user and branch name if
// provided as "user:ref-name", and by base name if provided, i.e. "master"
func (gc *GithubClient) ListPullRequests(org, repo, head, base string) ([]*github.PullRequest, error) {
	PRsListOptions := github.PullRequestListOptions{
		State: string(PullRequestAllState),
		Head:  head,
		Base:  base,
	}

	options := &github.ListOptions{}
	genericList, err := gc.depaginate(
		fmt.Sprintf("listing Pull Requests with head '%s' and base '%s'", head, base),
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.PullRequests.List(ctx, org, repo, &PRsListOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, PR := range page {
					interfaceList = append(interfaceList, PR)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.PullRequest, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.PullRequest)
	}
	return res, err
}

// ListCommits lists commits from a pull request
func (gc *GithubClient) ListCommits(org, repo string, ID int) ([]*github.RepositoryCommit, error) {
	options := &github.ListOptions{}
	genericList, err := gc.depaginate(
		fmt.Sprintf("listing commits in Pull Requests '%d'", ID),
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.PullRequests.ListCommits(ctx, org, repo, ID, nil)
			var interfaceList []interface{}
			if nil == err {
				for _, commit := range page {
					interfaceList = append(interfaceList, commit)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.RepositoryCommit, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.RepositoryCommit)
	}
	return res, err
}

// ListFiles lists files from a pull request
func (gc *GithubClient) ListFiles(org, repo string, ID int) ([]*github.CommitFile, error) {
	options := &github.ListOptions{}
	genericList, err := gc.depaginate(
		fmt.Sprintf("listing files in Pull Requests '%d'", ID),
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.PullRequests.ListFiles(ctx, org, repo, ID, nil)
			var interfaceList []interface{}
			if nil == err {
				for _, f := range page {
					interfaceList = append(interfaceList, f)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.CommitFile, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.CommitFile)
	}
	return res, err
}

// CreatePullRequest creates PullRequest, passing head user and branch name "user:ref-name", and base branch name like "master"
func (gc *GithubClient) CreatePullRequest(org, repo, head, base, title, body string) (*github.PullRequest, error) {
	b := true
	PR := &github.NewPullRequest{
		Title:               &title,
		Body:                &body,
		Head:                &head,
		Base:                &base,
		MaintainerCanModify: &b,
	}

	var res *github.PullRequest
	_, err := gc.retry(
		fmt.Sprintf("creating PullRequest from '%s' to '%s', title: '%s'. body: '%s'", head, base, title, body),
		maxRetryCount,
		func() (*github.Response, error) {
			var resp *github.Response
			var err error
			res, resp, err = gc.Client.PullRequests.Create(ctx, org, repo, PR)
			return resp, err
		},
	)
	return res, err
}