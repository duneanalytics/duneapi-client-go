package e2e

import (
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Submitting is not exercised here: every SubmitContracts call creates a real
// decoding submission for the key's owner.
func TestListContractSubmissions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	result, err := client.ListContractSubmissions(models.ListContractSubmissionsOptions{Limit: 5})
	require.NoError(t, err)
	assert.NotNil(t, result.Submissions)
	assert.GreaterOrEqual(t, result.Total, 0)
	assert.LessOrEqual(t, len(result.Submissions), 5)

	for _, submission := range result.Submissions {
		assert.NotEmpty(t, submission.ID)
		assert.NotEmpty(t, submission.BlockchainName)
		assert.NotEmpty(t, submission.Status)
	}
}

func TestListContractSubmissionsWithStatusFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	result, err := client.ListContractSubmissions(models.ListContractSubmissionsOptions{
		Limit:  5,
		Status: "rejected",
	})
	require.NoError(t, err)

	for _, submission := range result.Submissions {
		assert.Equal(t, "rejected", submission.Status)
	}
}
