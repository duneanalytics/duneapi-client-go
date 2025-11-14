package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDatasetResponseValidation(t *testing.T) {
	t.Run("valid dataset response", func(t *testing.T) {
		ds := DatasetResponse{
			Slug:      "dex.trades",
			Namespace: "dex",
			TableName: "trades",
			Owner: DatasetOwner{
				Handle: "dune",
				Type:   "team",
			},
		}
		require.NoError(t, ds.HasError())
	})

	t.Run("missing slug", func(t *testing.T) {
		ds := DatasetResponse{
			Namespace: "dex",
			TableName: "trades",
			Owner: DatasetOwner{
				Handle: "dune",
				Type:   "team",
			},
		}
		require.Error(t, ds.HasError())
		require.Contains(t, ds.HasError().Error(), "slug")
	})

	t.Run("missing namespace", func(t *testing.T) {
		ds := DatasetResponse{
			Slug:      "dex.trades",
			TableName: "trades",
			Owner: DatasetOwner{
				Handle: "dune",
				Type:   "team",
			},
		}
		require.Error(t, ds.HasError())
		require.Contains(t, ds.HasError().Error(), "namespace")
	})

	t.Run("missing table_name", func(t *testing.T) {
		ds := DatasetResponse{
			Slug:      "dex.trades",
			Namespace: "dex",
			Owner: DatasetOwner{
				Handle: "dune",
				Type:   "team",
			},
		}
		require.Error(t, ds.HasError())
		require.Contains(t, ds.HasError().Error(), "table_name")
	})

	t.Run("missing owner handle", func(t *testing.T) {
		ds := DatasetResponse{
			Slug:      "dex.trades",
			Namespace: "dex",
			TableName: "trades",
			Owner: DatasetOwner{
				Type: "team",
			},
		}
		require.Error(t, ds.HasError())
		require.Contains(t, ds.HasError().Error(), "owner handle")
	})
}

func TestListDatasetsResponseValidation(t *testing.T) {
	t.Run("valid list response", func(t *testing.T) {
		resp := ListDatasetsResponse{
			Datasets: []DatasetResponse{},
			Total:    0,
		}
		require.NoError(t, resp.HasError())
	})

	t.Run("missing datasets array", func(t *testing.T) {
		resp := ListDatasetsResponse{
			Total: 10,
		}
		require.Error(t, resp.HasError())
		require.Contains(t, resp.HasError().Error(), "datasets array")
	})
}

func TestDatasetResponseJSONMarshaling(t *testing.T) {
	t.Run("marshal and unmarshal dataset response", func(t *testing.T) {
		original := DatasetResponse{
			Slug:        "dex.trades",
			Name:        "DEX Trades",
			Namespace:   "dex",
			TableName:   "trades",
			Type:        "spell",
			IsPrivate:   false,
			Description: "Decentralized exchange trades",
			Columns: []DatasetColumn{
				{
					Name:        "blockchain",
					Type:        "varchar",
					Nullable:    false,
					Description: "Blockchain name",
				},
				{
					Name:     "amount",
					Type:     "double",
					Nullable: true,
				},
			},
			Owner: DatasetOwner{
				Handle: "dune",
				Type:   "team",
			},
			CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled DatasetResponse
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Equal(t, original.Slug, unmarshaled.Slug)
		require.Equal(t, original.Name, unmarshaled.Name)
		require.Equal(t, original.Namespace, unmarshaled.Namespace)
		require.Equal(t, original.TableName, unmarshaled.TableName)
		require.Equal(t, original.Type, unmarshaled.Type)
		require.Equal(t, original.IsPrivate, unmarshaled.IsPrivate)
		require.Equal(t, original.Description, unmarshaled.Description)
		require.Len(t, unmarshaled.Columns, 2)
		require.Equal(t, original.Owner.Handle, unmarshaled.Owner.Handle)
	})
}

func TestListDatasetsResponseJSONMarshaling(t *testing.T) {
	t.Run("marshal and unmarshal list response", func(t *testing.T) {
		original := ListDatasetsResponse{
			Datasets: []DatasetResponse{
				{
					Slug:      "dex.trades",
					Namespace: "dex",
					TableName: "trades",
					Owner: DatasetOwner{
						Handle: "dune",
						Type:   "team",
					},
				},
				{
					Slug:      "tokens.erc20",
					Namespace: "tokens",
					TableName: "erc20",
					Owner: DatasetOwner{
						Handle: "community",
						Type:   "user",
					},
				},
			},
			Total: 2,
		}

		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled ListDatasetsResponse
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		require.Equal(t, original.Total, unmarshaled.Total)
		require.Len(t, unmarshaled.Datasets, 2)
		require.Equal(t, original.Datasets[0].Slug, unmarshaled.Datasets[0].Slug)
		require.Equal(t, original.Datasets[1].Slug, unmarshaled.Datasets[1].Slug)
	})
}
