package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	actual := GetPathParam("/api/v1/agentPools/{poolId}/updates")
	assert.Equal(t, []string{"poolId"}, actual)
}

func TestParseMulti(t *testing.T) {
	actual := GetPathParam("/api/v1/agentPools/{poolId}/updates/{updateId}")
	assert.Equal(t, []string{"poolId", "updateId"}, actual)
}

func TestParseMultis(t *testing.T) {
	actual := GetPathParam("/api/v1/agentPools/{poolId}/updates/{updateId}/activate")
	assert.Equal(t, []string{"poolId", "updateId"}, actual)
}

func TestFirstToLower(t *testing.T) {
	actual := FirstToLower("ApiServiceIntegrations")
	assert.Equal(t, "apiServiceIntegrations", actual)
}

func TestExtractMap(t *testing.T) {
	actual := `{ "profile": { "description": "test", "name": "Test1" }, "type": "OKTA_GROUP"}`
	m, err := ExtractMap(actual)
	assert.NoError(t, err)
	assert.Equal(t, m["profile"].(map[string]interface{})["name"], "Test1")
}
