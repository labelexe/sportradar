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
const ncaafKey = "rdg8a4vhc3wmqte9jc4gwx52"
const ncaambKey = "9jq83aq4z5v4sbqdrssu7g4p"

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

	schedule, err := client.MLBSchedule(2022, sr.SeasonTypeRegular)
	assert.NoError(t, err)
	assert.Equal(t, "MLB", schedule.League.Alias)
	assert.NotEmpty(t, schedule.Games)
	assert.NotEmpty(t, schedule.Games[0].ID)
	assert.NotEmpty(t, schedule.Games[0].Broadcast)
	assert.NotEmpty(t, schedule.Games[0].Broadcast.Network)
}

func TestNBASchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NBA: nbaKey,
		},
	})

	schedule, err := client.NBASchedule(2022, sr.SeasonTypeRegular)
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

	schedule, err := client.NHLSchedule(2022, sr.SeasonTypeRegular)
	assert.NoErrorf(t, err, "%s", err)
	assert.NotEmpty(t, schedule.Season)
	assert.NotEmpty(t, schedule.Season.ID)
	assert.NotEmpty(t, schedule.Games[0])
	assert.NotEmpty(t, schedule.Games[0].ID)
	assert.NotEmpty(t, schedule.Games[0].HomeTeam)
	assert.NotEmpty(t, schedule.Games[0].AwayTeam)
	assert.NotEmpty(t, schedule.Games[0].Broadcasts)
	assert.NotEmpty(t, schedule.Games[0].Broadcasts[0].Network)
}

func TestNCAAFSchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NCAAF: ncaafKey,
		},
	})

	schedule, err := client.NCAAFSchedule(2022, sr.SeasonTypeRegular)
	assert.NoErrorf(t, err, "%s", err)
	assert.NotEmpty(t, schedule.ID)
	assert.NotEmpty(t, schedule.Weeks)
	assert.NotEmpty(t, schedule.Weeks[0].ID)
	assert.NotEmpty(t, schedule.Weeks[0].Games)
	assert.NotEmpty(t, schedule.Weeks[0].Games[0].ID)
	assert.NotEmpty(t, schedule.Weeks[0].Games[0].Broadcast)
}

func TestNCAAMBSchedule(t *testing.T) {
	client := NewClient(ClientConfig{
		Keys: LeagueKeys{
			NCAAMB: ncaambKey,
		},
	})

	schedule, err := client.NCAAMBSchedule(2022, sr.SeasonTypeRegular)
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
