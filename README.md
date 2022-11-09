# DuneAPI client
DuneAPI CLI and client library for Go

## Library usage

To add this library to your go project run:

```
go get github.com/duneanalytics/duneapi-client-go
```

First you have to define the API key that will be used to authenticate with
the Dune API. There are two ways to achieve this.

```go
import (
	"github.com/duneanalytics/duneapi-client-go/config"
)

func main() {
	// Read API key from DUNE_API_KEY environment variable
	env, err := config.ParseEnv()
	if err != nil {
		// Handle error
	}

	// Define it from your code
	env = &config.Env{
		APIKey: "Your_API_Key",
		// optionally, you can define a different api domain to connect to
		// Host: "https://duneapi.example.com",
	}
}
```

Next you can instantiate and use a Dune client object:

```go
import (
	"github.com/duneanalytics/duneapi-client-go/dune"
)

func main() {
	client := dune.NewDuneClient(env)
	queryID := 1234
	queryParameters := map[string]any{
		"paramKey": "paramValue",
	}
	rows, err := client.RunQueryGetRows(queryID, queryParameters)
	if err != nil {
		// handle error
	}

	for row := range rows {
		// ...
	}
}
```

The RunQueryGetRows will execute the query, wait for completion and return
only an array of rows, without any metadata. For other ways to use the client,
check out the [package documentation](https://pkg.go.dev/github.com/duneanalytics/duneapi-client-go).

## CLI usage

### Build

```
go build -o dunecli cmd/main.go
```

You can use it from the repo directly or copy to a directory in your `$PATH`

### Usage

The CLI has 2 main modes of operation. Run a query or retrieve information about
an existing execution. In both cases, it will print out raw minified JSON to stdout,
so if you want to prettify it, or select a specific key, you can pipe to [jq](https://stedolan.github.io/jq/).

#### Execute a query

To trigger a query execution and print the results once it's done:

```bash
DUNE_API_KEY=<your_key> ./dunecli -q <query_id>
```

If the query has parameters you want to override, use:

```bash
DUNE_API_KEY=<your_key> ./dunecli -q <query_id> -p '{"<param_key>": "<param_value>"}'
```

For numeric parameters, omit the quotes around the value.

#### Get results for an existing execution

If you already have an execution ID, you can retrieve its results (or state if it
hasn't completed yet) with this:

```bash
DUNE_API_KEY=<your_key> ./dunecli -e <execution_id>
```
