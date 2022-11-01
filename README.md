# DuneAPI client
DuneAPI CLI and client library for Go

## CLI

### Build

```
go build -o dunecli cmd/main.go
```

You can use it from the repo directly or copy to a directory in your `$PATH`

### Usage

#### Execute a query

To trigger a query execution and print the results once it's done:

```bash
DUNE_API_KEY=<your_key> ./dunecli -q <query_id>
```

If the query has parameters you want to override, use:

```bash
DUNE_API_KEY=<your_key> ./dunecli -q <query_id> -p <param_key>:<param_vlue>
```
