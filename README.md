# go-pready

Use Github GraphQL API to return observed repositories' Pull Requests info (*approval* status, mainly).

## Internals

A nice `PullRequest` struct with the followings is provided:

```
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

func (pr *PullRequest) Wip() bool {
  ...
}

func (pr *PullRequest) Reviewed() bool {
  ...
}

func (pr *PullRequest) Approved() bool {
  ...
}
```

## Configuration

Change `repositories` in [`constants.go`](../master/constants.go).

```
go build
GITHUB_TOKEN=My1337GithubAPIToken ./pready
```

## Why

It is f*cking functional (GraphQL) and fast (Go).

And it was a boring Saturday.

## TODO

  - Connect (internally?) to a Slack bot (didn't find any nice implementation in go yet... disappointed)
  - Test (oh well, ...)

## License: MIT

Copyright (C) 2017 Giuseppe Lobraico

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.