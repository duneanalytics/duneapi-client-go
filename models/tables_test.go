package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTableListResponseValidation(t *testing.T) {
	t.Run("valid table list response", func(t *testing.T) {
		resp := TableListResponse{
			Tables: []TableListElement{},
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("missing tables array", func(t *testing.T) {
		resp := TableListResponse{}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "tables array")
	})
}

func TestTableCreateResponseValidation(t *testing.T) {
	t.Run("valid create response", func(t *testing.T) {
		resp := TableCreateResponse{
			Namespace:    "my_user",
			TableName:    "test_table",
			FullName:     "dune.my_user.test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("missing namespace", func(t *testing.T) {
		resp := TableCreateResponse{
			TableName:    "test_table",
			FullName:     "dune.my_user.test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "namespace")
	})

	t.Run("missing table_name", func(t *testing.T) {
		resp := TableCreateResponse{
			Namespace:    "my_user",
			FullName:     "dune.my_user.test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "table_name")
	})

	t.Run("missing full_name", func(t *testing.T) {
		resp := TableCreateResponse{
			Namespace:    "my_user",
			TableName:    "test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "full_name")
	})
}

func TestCSVUploadResponseValidation(t *testing.T) {
	t.Run("valid CSV upload response", func(t *testing.T) {
		resp := CSVUploadResponse{
			TableName:    "test_table",
			FullName:     "dune.my_user.test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("missing table_name", func(t *testing.T) {
		resp := CSVUploadResponse{
			FullName:     "dune.my_user.test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "table_name")
	})

	t.Run("missing full_name", func(t *testing.T) {
		resp := CSVUploadResponse{
			TableName:    "test_table",
			ExampleQuery: "SELECT * FROM dune.my_user.test_table",
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "full_name")
	})
}

func TestTableInsertResponseValidation(t *testing.T) {
	t.Run("valid insert response", func(t *testing.T) {
		resp := TableInsertResponse{
			BytesWritten: 1024,
			RowsWritten:  10,
			TableName:    "test_table",
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("missing table_name", func(t *testing.T) {
		resp := TableInsertResponse{
			BytesWritten: 1024,
			RowsWritten:  10,
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "table_name")
	})
}

func TestTableDeleteResponseValidation(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		resp := TableDeleteResponse{
			Success: true,
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("failed deletion", func(t *testing.T) {
		resp := TableDeleteResponse{
			Success: false,
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "deletion failed")
	})
}

func TestTableClearResponseValidation(t *testing.T) {
	t.Run("successful clear", func(t *testing.T) {
		resp := TableClearResponse{
			Success: true,
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("failed clear", func(t *testing.T) {
		resp := TableClearResponse{
			Success: false,
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "clear failed")
	})
}

func TestTableCreateRequestJSONMarshaling(t *testing.T) {
	t.Run("marshal create request", func(t *testing.T) {
		req := TableCreateRequest{
			Namespace:   "my_user",
			TableName:   "interest_rates",
			Description: "10 year daily interest rates",
			IsPrivate:   false,
			Schema: []TableColumn{
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
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled TableCreateRequest
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Equal(t, req.Namespace, unmarshaled.Namespace)
		require.Equal(t, req.TableName, unmarshaled.TableName)
		require.Equal(t, req.Description, unmarshaled.Description)
		require.Equal(t, req.IsPrivate, unmarshaled.IsPrivate)
		require.Len(t, unmarshaled.Schema, 2)
		require.Equal(t, "date", unmarshaled.Schema[0].Name)
		require.Equal(t, "timestamp", unmarshaled.Schema[0].Type)
	})
}

func TestCSVUploadRequestJSONMarshaling(t *testing.T) {
	t.Run("marshal CSV upload request", func(t *testing.T) {
		req := CSVUploadRequest{
			TableName:   "my_table",
			Data:        "col1,col2\nval1,val2",
			Description: "Test table",
			IsPrivate:   false,
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled CSVUploadRequest
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Equal(t, req.TableName, unmarshaled.TableName)
		require.Equal(t, req.Data, unmarshaled.Data)
		require.Equal(t, req.Description, unmarshaled.Description)
		require.Equal(t, req.IsPrivate, unmarshaled.IsPrivate)
	})
}

func TestTableListResponseJSONMarshaling(t *testing.T) {
	t.Run("marshal and unmarshal table list", func(t *testing.T) {
		original := TableListResponse{
			Tables: []TableListElement{
				{
					FullName:       "dune.my_user.test_table",
					IsPrivate:      false,
					TableSizeBytes: "1024",
					CreatedAt:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Owner: TableOwner{
						Handle: "my_user",
						Type:   "user",
					},
					Columns: []TableColumn{
						{
							Name:     "id",
							Type:     "bigint",
							Nullable: false,
						},
						{
							Name:     "name",
							Type:     "varchar",
							Nullable: true,
						},
					},
				},
			},
		}

		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled TableListResponse
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Len(t, unmarshaled.Tables, 1)
		require.Equal(t, original.Tables[0].FullName, unmarshaled.Tables[0].FullName)
		require.Equal(t, original.Tables[0].IsPrivate, unmarshaled.Tables[0].IsPrivate)
		require.Equal(t, original.Tables[0].Owner.Handle, unmarshaled.Tables[0].Owner.Handle)
		require.Len(t, unmarshaled.Tables[0].Columns, 2)
	})
}

func TestTableColumnSchema(t *testing.T) {
	t.Run("table column with all fields", func(t *testing.T) {
		col := TableColumn{
			Name:        "blockchain",
			Type:        "varchar",
			Nullable:    false,
			Description: "The blockchain name",
			Metadata: map[string]any{
				"filtering_column": true,
			},
		}

		jsonData, err := json.Marshal(col)
		require.NoError(t, err)

		var unmarshaled TableColumn
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Equal(t, col.Name, unmarshaled.Name)
		require.Equal(t, col.Type, unmarshaled.Type)
		require.Equal(t, col.Nullable, unmarshaled.Nullable)
		require.Equal(t, col.Description, unmarshaled.Description)
		require.NotNil(t, unmarshaled.Metadata)
		require.Equal(t, true, unmarshaled.Metadata["filtering_column"])
	})
}
