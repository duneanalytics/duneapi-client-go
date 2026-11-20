package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateQueryRequestTagsAndParameters(t *testing.T) {
	name := "n"
	cases := []struct {
		name string
		req  UpdateQueryRequest
		want string
	}{
		{"nil keeps the current value", UpdateQueryRequest{Name: &name}, `{"name":"n"}`},
		{"empty clears", UpdateQueryRequest{Tags: []string{}, Parameters: []QueryParameter{}},
			`{"tags":[],"parameters":[]}`},
		{"set replaces", UpdateQueryRequest{
			Tags:       []string{"a"},
			Parameters: []QueryParameter{{Key: "k", Type: "text", Value: "v"}},
		},
			`{"tags":["a"],"parameters":[{"key":"k","type":"text","value":"v"}]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.req)
			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(got))
		})
	}
}
