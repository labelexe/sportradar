package sportradar

import (
	"encoding/json"
	"testing"
	"time"

	sr "github.com/playback-sports/sportradar/pkg/base"
	"github.com/stretchr/testify/assert"
)

const mlbKey = "cms9sf848zh9tnptnp62m3ts"

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func TestMLBDailySummary(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeageuKeys{
			MLB: mlbKey,
		},
	})

	now := time.Now()
	summary, err := client.MLBDailySummary(now)
	assert.NoError(t, err)
	assert.Equal(t, "MLB", summary.League.Alias)
	assert.True(t, len(summary.League.Games) > 1)
	assert.NotEmpty(t, summary.League.Games[0].ID)
}

func TestMLBSchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeageuKeys{
			MLB: mlbKey,
		},
	})

	now := time.Now()
	summary, err := client.MLBSchedule(now, sr.SeasonTypeRegular)
	assert.NoError(t, err)
	assert.Equal(t, "MLB", summary.League.Alias)
	assert.True(t, len(summary.Games) > 1)
	assert.NotEmpty(t, summary.Games[0].ID)
	assert.NotEmpty(t, summary.Games[0].Broadcast)
	assert.NotEmpty(t, summary.Games[0].Broadcast.Network)
}
