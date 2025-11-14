# DuneAPI client
DuneAPI CLI and client library for Go

## Library usage

To add this library to your go project run:

```
go get github.com/duneanalytics/duneapi-client-go
```

First you have to define the configuration that will be used to authenticate
with the Dune API. There are three ways to achieve this.


```go
import (
	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/dune"
)

func main() {
	// Use one of the following options
	// Read config from DUNE_API_KEY and DUNE_API_HOST environment variables
	env, err := config.FromEnvVars()
	if err != nil {
		// handle error
	}

	// Define it from your code
	env = config.FromAPIKey("Your_API_Key")

	// Define manually
	env = &config.Env{
		APIKey: "Your_API_Key",
		// you can define a different domain to connect to, for example for a mocked API
		Host: "https://api.example.com",
	}

	// Next, instantiate and use a Dune client object
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

### Dataset Discovery APIs

The client provides methods to discover and explore datasets available on Dune:

```go
// List all datasets with optional filtering
datasets, err := client.ListDatasets(
	50,           // limit
	0,            // offset
	"dune",       // owner_handle (optional, use "" to skip)
	"spell",      // dataset_type (optional, use "" to skip)
)
if err != nil {
	// handle error
}

for _, dataset := range datasets.Datasets {
	fmt.Printf("Dataset: %s (%s.%s)\n", dataset.Slug, dataset.Namespace, dataset.TableName)
	fmt.Printf("  Owner: %s\n", dataset.Owner.Handle)
	fmt.Printf("  Columns: %d\n", len(dataset.Columns))
}

// Get detailed information about a specific dataset
dataset, err := client.GetDataset("dex.trades")
if err != nil {
	// handle error
}

fmt.Printf("Dataset: %s\n", dataset.Name)
fmt.Printf("Description: %s\n", dataset.Description)
for _, col := range dataset.Columns {
	fmt.Printf("  - %s (%s, nullable: %v)\n", col.Name, col.Type, col.Nullable)
}
```

### Table Management APIs

The client provides comprehensive methods for managing uploaded tables:

```go
// List all uploaded tables
tables, err := client.ListUploads(50, 0)
if err != nil {
	// handle error
}

for _, table := range tables.Tables {
	fmt.Printf("Table: %s (size: %s bytes)\n", table.FullName, table.TableSizeBytes)
}

// Create a new table with defined schema
createResp, err := client.CreateUpload(models.UploadsCreateRequest{
	Namespace:   "my_user",
	TableName:   "interest_rates",
	Description: "10 year daily interest rates",
	IsPrivate:   false,
	Schema: []models.UploadsColumn{
		{
			Name:     "date",
			Type:     "timestamp",
			Nullable: false,
		},
		{
			Name:     "rate",
			Type:     "double",
			Nullable: true,
		},
	},
})
if err != nil {
	// handle error
}
fmt.Printf("Created table: %s\n", createResp.FullName)

// Upload CSV data to create a new table
csvResp, err := client.UploadCSV(models.UploadsCSVRequest{
	TableName:   "my_table",
	Data:        "col1,col2\nval1,val2\nval3,val4",
	Description: "My test table",
	IsPrivate:   false,
})
if err != nil {
	// handle error
}
fmt.Printf("Uploaded CSV to: %s\n", csvResp.FullName)

// Insert data into an existing table (CSV format)
insertResp, err := client.InsertIntoUpload(
	"my_user",
	"interest_rates",
	"2024-01-01,3.5\n2024-01-02,3.6",
	"text/csv",
)
if err != nil {
	// handle error
}
fmt.Printf("Inserted %d rows (%d bytes)\n", insertResp.RowsWritten, insertResp.BytesWritten)

// Insert data in NDJSON format
ndjsonData := `{"date":"2024-01-03","rate":3.7}
{"date":"2024-01-04","rate":3.8}`
insertResp, err = client.InsertIntoUpload(
	"my_user",
	"interest_rates",
	ndjsonData,
	"application/x-ndjson",
)

// Clear all data from a table (preserves schema)
clearResp, err := client.ClearUpload("my_user", "interest_rates")
if err != nil {
	// handle error
}

// Delete a table permanently
deleteResp, err := client.DeleteUpload("my_user", "interest_rates")
if err != nil {
	// handle error
}
```

#### Deprecated Table Endpoints

⚠️ **DEPRECATION NOTICE**: The following methods use deprecated `/v1/table/*` endpoints and will be removed on **March 1, 2026**. Please migrate to the new methods shown above.

```go
// DEPRECATED: Use ListUploads instead
tables, err := client.ListTables(50, 0)

// DEPRECATED: Use CreateUpload instead
createResp, err := client.CreateTable(req)

// DEPRECATED: Use UploadCSV instead
csvResp, err := client.UploadCSVDeprecated(req)

// DEPRECATED: Use DeleteUpload instead
deleteResp, err := client.DeleteTable("my_user", "table_name")

// DEPRECATED: Use ClearUpload instead
clearResp, err := client.ClearTable("my_user", "table_name")

// DEPRECATED: Use InsertIntoUpload instead
insertResp, err := client.InsertTable("my_user", "table_name", data, contentType)
```

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
