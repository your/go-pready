# go-pready

Slack bot notifications for Pull Requests' review/approval status.

## Configuration

Change at least `notificationChannel` in [`constants.go`](../master/constants.go).

```
go build
GITHUB_TOKEN=My1337GithubAPIToken \
SLACK_TOKEN=My1338SlackAPIToken ./pready
```

## Slack commands

Invite the bot to the notification channel and type:

* `review watch <Github repo URL>`
* `review unwatch <Github repo URL>`

...more to come. :sparkle:

## Demo

```
...

2017/07/18 01:25:32 :cry: PR #1 "Test PR" is still waiting for approval! :arrow_right: https://github.com/your/_repo/pull/1

...

2017/07/18 01:25:56 :champagne: PR #1 "Test PR" is still waiting for merge! :arrow_right: https://github.com/your/_repo/pull/1

...

2017/07/18 01:25:57 :tada: There are no pending approval PRs for repository `_repo` :sunglasses: great!

...
```

## Why

It is f*cking functional (GraphQL) and fast (Go).

And it was a boring Saturday.

## TODO

  - Slack commands (in progress)
  - Github webhooks(?)
  - Test (oh well, ...)

## License: MIT

Copyright (C) 2017 Giuseppe Lobraico

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.