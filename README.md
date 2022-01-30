# readtime

[![Go Report Card](https://goreportcard.com/badge/github.com/91go/readtime)](https://goreportcard.com/report/github.com/91go/readtime)
[![codecov](https://codecov.io/gh/91go/readtime/branch/main/graph/badge.svg?token=VLP0KPQ5FS)](https://codecov.io/gh/91go/readtime)


Get article reading time, support multiple languages

```markdown

go get github.com/91go/readtime

```

## usage

```markdown

NewReadTime().Read("").ToMap()

```

or read file or url directly

```markdown

NewReadTime().ReadFile("").ToMap()

NewReadTime().ReadURL("").ToMap()

```

use `SetTranslation()` to convert language

```markdown

NewReadTime().ReadFile("filename").SetTranslation("ca").ToMap()

```

use `ToJSON()` to get json

```markdown

NewReadTime().ReadFile("").ToJSON()

```