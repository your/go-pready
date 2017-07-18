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
	Reviews []Review
}

type Review struct {
	AuthorLogin string
	State       string
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
	reviewers := pr.Reviewers()
	states := []string{}

	for _, reviewer := range reviewers {
		for i := len(pr.Reviews) - 1; i >= 0; i-- {
			if pr.Reviews[i].AuthorLogin == reviewer {
				states = append(states, pr.Reviews[i].State)
				break
			}
		}
	}

	return stringOccurencesInSlice("APPROVED", states) >= minApprovals &&
		!isStringInSlice("CHANGES_REQUESTED", states)
}

// Reviewers is a PullRequest's struct method that returns all the reviewers, without duplicates.
func (pr *PullRequest) Reviewers() []string {
	reviewers := []string{}

	for _, review := range pr.Reviews {
		reviewers = append(reviewers, review.AuthorLogin)
	}

	return uniqueSlice(reviewers)
}

func checkRepositories() {
	if len(repositories) == 0 {
		log.Println("No repositories to loop through, you might want to add some.")
		return
	}

	for _, repo := range repositories {
		resp, err := resty.R().
			SetHeader("Authorization", "bearer "+os.Getenv("GITHUB_TOKEN")).
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

		pendingReviewPRs := []PullRequest{}
		pendingApprovalPRs := []PullRequest{}
		pendingMergePRs := []PullRequest{}

		for _, pr := range repo.PullRequests {
			if !pr.Wip() {
				if !pr.Reviewed() {
					pendingReviewPRs = append(pendingReviewPRs, pr)
					continue
				}
				if pr.Approved() {
					pendingMergePRs = append(pendingMergePRs, pr)
				} else {
					pendingApprovalPRs = append(pendingApprovalPRs, pr)
				}
			}
		}

		for _, pr := range pendingReviewPRs {
			log.Printf("PR #%s \"%s\" is still waiting for review! –– %s\n", strconv.Itoa(pr.Number), pr.Title, pr.URL)
		}

		for _, pr := range pendingMergePRs {
			log.Printf("PR #%s \"%s\" is still waiting for merge! –– %s\n", strconv.Itoa(pr.Number), pr.Title, pr.URL)
		}

		if len(pendingApprovalPRs) == 0 {
			log.Printf("There are no pending approval PRs for repository \"%s\" –– great!", repo.Name)
			return
		}

		for _, pr := range pendingApprovalPRs {
			log.Printf("PR #%s \"%s\" is still waiting for approval! –– %s\n", strconv.Itoa(pr.Number), pr.Title, pr.URL)
		}
	}
}

func buildRepositoryFromResponse(response *GraphQLResponseBody) Repository {
	results := response.Data.Viewer.Repository

	pullRequests := make([]PullRequest, results.PullRequests.TotalCount)

	for i := 0; i < results.PullRequests.TotalCount; i++ {
		pr := results.PullRequests.Edges[i].Node

		labels := make([]string, pr.Labels.TotalCount)
		reviews := make([]Review, pr.Reviews.TotalCount)

		for j := 0; j < pr.Labels.TotalCount; j++ {
			labels[j] = pr.Labels.Edges[j].Node.Name
		}

		for j := 0; j < pr.Reviews.TotalCount; j++ {
			review := pr.Reviews.Edges[j].Node
			reviews[j] = Review{
				AuthorLogin: review.Author.Login,
				State:       review.State,
			}
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
