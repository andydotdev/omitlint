# omitlint

[![Test Status](https://github.com/andydotdev/omitlint/actions/workflows/ci.yml/badge.svg)](https://andy.dev/omitlint/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/andy.dev/omitlint)](https://goreportcard.com/report/andy.dev/omitlint) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Checks for impossible JSON `omitempty` options, indicating fields that will never be omitted when encoding. 

For example:

```Go
// Options are ActionRequest optional items.
type Options struct {
	Values []string `json:"values,omitempty"
}

// ActionRequest submits an action to the scheduler.
type ActionRequest struct {
	Action  string  `json:"action"`
	Options Options `json:"options,omitempty"
}
```

In this example, `options` will always appear in the encoded JSON, even if `Value` is empty or nil.  This is because `Options` is a concrete struct type and not a pointer. `encoding/json` does not do any sort of field checking for struct types, and a concrete struct type will always have a zero value.

In order to correct, this you need to make it a pointer type:

```Go
// ActionRequest submits an action to the scheduler.
type ActionRequest struct {
	Action   string  `json:"action"`
	Options *Options `json:"options,omitempty"
}
```

This can be easy to overlook, especially when dealing with types that come from elsewhere. This linter will take care of this check for you by outputting messages in the form of:

```
/project/api/action_request.go:10:2: field "Options" is marked "omitempty", but cannot not be omitted.
```

## Installation

Download `omitlint` from the [releases](https://github.com/andydotdev/omitlint/releases) or get the latest version from source with:

```shell
go get andy.dev/omitlint/cmd/omitlint
```

## Usage

### Shell

Check everything:

```shell
omitlint ./...
```
