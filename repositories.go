package main

import (
	"log"
	"strconv"

	"github.com/go-resty/resty"
)

// Repository represents one repository containing all open pull requests.
type Repository struct {
	Name         string
	PullRequests []PullRequest
}

// PullRequest represents one pull request with some descriptive data.
type PullRequest struct {
	Number  int
	Title   string
	URL     string
	Labels  []string
	Reviews []string
}

// Wip is a PullRequest's struct method that returns true if the PR is labeled as "WIP".
func (pr *PullRequest) Wip() bool {
	return isStringInSlice("wip", pr.Labels)
}

// Reviewed is a PullRequest's struct method that returns true if there are any reviews,
// false otherwise.
func (pr *PullRequest) Reviewed() bool {
	if len(pr.Reviews) > 0 {
		return true
	}
	return false
}

// Approved is a PullRequest's struct method that returns true if there are enough
// approved reviews and no changes requested, false otherwise.
func (pr *PullRequest) Approved() bool {
	return stringOccurencesInSlice("APPROVED", pr.Reviews) >= minApprovals &&
		!isStringInSlice("CHANGES_REQUESTED", pr.Reviews)
}

func checkRepositories() {
	if len(repositories) == 0 {
		log.Println("No repositories to loop through, you might want to add some.")
		return
	}

	for _, repo := range repositories {
		resp, err := resty.R().
			SetHeader("Authorization", "bearer "+githubAPIToken).
			SetBody(buildGraphQLRequestBody(repo)).
			SetResult(&GraphQLResponseBody{}).
			Post("https://api.github.com/graphql")

		if err != nil {
			log.Println("API connection error:", err)
			return
		}

		repo := buildRepositoryFromResponse(resp.Result().(*GraphQLResponseBody))

		if len(repo.Name) == 0 {
			log.Println("Something went wrong! –– Did you specify a valid Github API Token?")
			return
		}

		pendingApprovalPRs := []PullRequest{}
		pendingMergePRs := []PullRequest{}

		for _, pr := range repo.PullRequests {
			if !pr.Wip() {
				if !pr.Reviewed() || !pr.Approved() {
					pendingApprovalPRs = append(pendingApprovalPRs, pr)
				} else if pr.Approved() {
					pendingMergePRs = append(pendingMergePRs, pr)
				}
			}
		}

		if len(pendingApprovalPRs) == 0 {
			log.Printf("There are no pending approval PRs for repository \"%s\" –– great!", repo.Name)
			return
		}

		for _, pr := range pendingApprovalPRs {
			log.Printf("PR #%s \"%s\" is still waiting for approval! –– %s\n", strconv.Itoa(pr.Number), pr.Title, pr.URL)
		}

		for _, pr := range pendingMergePRs {
			log.Printf("PR #%s \"%s\" is still waiting for merge! –– %s\n", strconv.Itoa(pr.Number), pr.Title, pr.URL)
		}
	}
}

func buildRepositoryFromResponse(response *GraphQLResponseBody) Repository {
	results := response.Data.Viewer.Repository

	pullRequests := make([]PullRequest, results.PullRequests.TotalCount)

	for i := 0; i < results.PullRequests.TotalCount; i++ {
		pr := results.PullRequests.Edges[i].Node

		labels := make([]string, pr.Labels.TotalCount)
		reviews := make([]string, pr.Reviews.TotalCount)

		for j := 0; j < pr.Labels.TotalCount; j++ {
			labels[j] = pr.Labels.Edges[j].Node.Name
		}

		for j := 0; j < pr.Reviews.TotalCount; j++ {
			reviews[j] = pr.Reviews.Edges[j].Node.State
		}

		pullRequests[i] = PullRequest{
			Number:  pr.Number,
			Title:   pr.Title,
			URL:     pr.URL,
			Labels:  labels,
			Reviews: reviews,
		}
	}

	return Repository{
		Name:         results.Name,
		PullRequests: pullRequests,
	}
}
