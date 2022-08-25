package sportradar

import (
	"testing"
	"time"

	sr "github.com/playback-sports/sportradar/pkg/base"
	assert "github.com/stretchr/testify/require"
)

const mlbKey = "cms9sf848zh9tnptnp62m3ts"
const nflKey = "usdsgdvez9tbdfmk38yhv3ar"

func TestMLBDailySummary(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
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
		Keys: LeagueKeys{
			MLB: mlbKey,
		},
	})

	now := time.Now()
	schedule, err := client.MLBSchedule(now, sr.SeasonTypeRegular)
	assert.NoError(t, err)
	assert.Equal(t, "MLB", schedule.League.Alias)
	assert.True(t, len(schedule.Games) > 1)
	assert.NotEmpty(t, schedule.Games[0].ID)
	assert.NotEmpty(t, schedule.Games[0].Broadcast)
	assert.NotEmpty(t, schedule.Games[0].Broadcast.Network)
}

func TestNFLSchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NFL: nflKey,
		},
	})

	now := time.Now()
	schedule, err := client.NFLSchedule(now, sr.SeasonTypeRegular)
	assert.NoErrorf(t, err, "%s", err)
	assert.NotEmpty(t, schedule.ID)
	assert.True(t, len(schedule.Weeks) > 1)
	assert.NotEmpty(t, schedule.Weeks[0].ID)
	assert.NotEmpty(t, schedule.Weeks[0].Games)
	assert.NotEmpty(t, schedule.Weeks[0].Games[0].ID)
}
