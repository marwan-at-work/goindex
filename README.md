# Go Index

A simple client for https://index.golang.org

### Install

go get marwan.io/goindex

### Usage

```golang
var c Client // the zero value is ready to be used
moduleVersions, err := c.Get(context.Background(), time.Time{})
// ...
```

You can use the second argument to specify a "since" argument, or you can call next on the response to get the next paginated list:

```golang
var c Client // the zero value is ready to be used
moduleVersions, err := c.Get(context.Background(), time.Time{})
// check err
nextList, err := moduleVersions.Next(context.Background(), &c)
// ...
```

### Status

Quite simple and early, breaking changes might occur such as adding the "limit" query.
