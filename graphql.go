package main

import (
	"encoding/json"
	"fmt"
)

// No changes here, please.

var graphQLquery = `query ($limit: Int!) {
  viewer {
    repository(name: "%s") {
      name
      pullRequests(first: $limit, states: OPEN) {
        totalCount
        edges {
          node {
            number
            title
            url
            labels(first: $limit) {
              totalCount
              edges {
                node {
                  name
                }
              }
            }
            reviews(first: $limit, states: [APPROVED, CHANGES_REQUESTED]) {
              totalCount
              edges {
                node {
                  author {
                    login
                  }
                  state
                }
              }
            }
          }
        }
      }
    }
  }
}`

// 100 should be enough forever, come on.
var graphQLvariables = `{
  "limit": 100
}`

// GraphQLRequestBody represents a valid GraphQL request body.
type GraphQLRequestBody struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// GraphQLResponseBody represents the expected GraphQL response body.
type GraphQLResponseBody struct {
	Data struct {
		Viewer struct {
			Repository struct {
				Name         string `json:"name"`
				PullRequests struct {
					TotalCount int `json:"totalCount"`
					Edges      []struct {
						Node struct {
							Number int    `json:"number"`
							Title  string `json:"title"`
							URL    string `json:"url"`
							Labels struct {
								TotalCount int `json:"totalCount"`
								Edges      []struct {
									Node struct {
										Name string `json:"name"`
									} `json:"node"`
								} `json:"edges"`
							} `json:"labels"`
							Reviews struct {
								TotalCount int `json:"totalCount"`
								Edges      []struct {
									Node struct {
										Author struct {
											Login string `json:"login"`
										} `json:"author"`
										State string `json:"state"`
									} `json:"node"`
								} `json:"edges"`
							} `json:"reviews"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"pullRequests"`
			} `json:"repository"`
		} `json:"viewer"`
	} `json:"data"`
}

func buildGraphQLRequestBody(repository string) string {
	r := GraphQLRequestBody{
		Query:     fmt.Sprintf(graphQLquery, repository),
		Variables: graphQLvariables,
	}

	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(body)
}
