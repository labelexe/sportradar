package sportradar

import (
	"testing"
	"time"

	sr "github.com/playback-sports/sportradar/pkg/base"
	assert "github.com/stretchr/testify/require"
)

const mlbKey = "cms9sf848zh9tnptnp62m3ts"
const nflKey = "usdsgdvez9tbdfmk38yhv3ar"
const nbaKey = "ugfpbjbn9d8npnmrhqzq7g49"
const nhlKey = "msuuhzsh7v3bfqw24m4nze6s"

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
	assert.NotEmpty(t, summary.League.Games)
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
	assert.NotEmpty(t, schedule.Games)
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
	assert.NotEmpty(t, schedule.Weeks)
	assert.NotEmpty(t, schedule.Weeks[0].ID)
	assert.NotEmpty(t, schedule.Weeks[0].Games)
	assert.NotEmpty(t, schedule.Weeks[0].Games[0].ID)
}

func TestNBASchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NBA: nbaKey,
		},
	})

	now := time.Now()
	schedule, err := client.NBASchedule(now, sr.SeasonTypeRegular)
	assert.NoErrorf(t, err, "%s", err)
	assert.NotEmpty(t, schedule.Season)
	assert.NotEmpty(t, schedule.Season.ID)
	assert.NotEmpty(t, schedule.Games[0])
	assert.NotEmpty(t, schedule.Games[0].ID)
	assert.NotEmpty(t, schedule.Games[0].Broadcasts)
	assert.NotEmpty(t, schedule.Games[0].Broadcasts[0].Network)
	assert.NotEmpty(t, schedule.Games[0].HomeTeam)
	assert.NotEmpty(t, schedule.Games[0].AwayTeam)
}

func TestNHLSchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NHL: nhlKey,
		},
	})

	now := time.Now()
	schedule, err := client.NHLSchedule(now, sr.SeasonTypeRegular)
	assert.NoErrorf(t, err, "%s", err)
	assert.NotEmpty(t, schedule.Season)
	assert.NotEmpty(t, schedule.Season.ID)
	assert.NotEmpty(t, schedule.Games[0])
	assert.NotEmpty(t, schedule.Games[0].ID)
	assert.NotEmpty(t, schedule.Games[0].HomeTeam)
	assert.NotEmpty(t, schedule.Games[0].AwayTeam)
}
