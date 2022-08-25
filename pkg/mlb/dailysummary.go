package mlb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	sportradar "github.com/playback-sports/sportradar/pkg/base"
)

type srDailySummaryGameMap struct {
	Game srDailySummaryGame `json:"game"`
}

type srDailySummaryGame struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	Coverage     string               `json:"coverage"`
	GameNumber   int64                `json:"game_number"`
	DayNight     string               `json:"day_night"`
	Scheduled    string               `json:"scheduled"`
	HomeTeam     string               `json:"home_team"`
	AwayTeam     string               `json:"away_team"`
	DoubleHeader bool                 `json:"double_header"`
	EntryMode    string               `json:"entry_mode"`
	Reference    string               `json:"reference"`
	Venue        sportradar.Venue     `json:"venue"`
	Broadcast    sportradar.Broadcast `json:"broadcast"`
	Weather      sportradar.Weather   `json:"weather"`
	Home         srDailySummaryTeam   `json:"home"`
	Away         srDailySummaryTeam   `json:"away"`
}

type DailySummaryGame struct {
	ID           string                `json:"id"`
	Status       sportradar.GameStatus `json:"status"`
	Coverage     string                `json:"coverage"`
	GameNumber   int64                 `json:"game_number"`
	DayNight     string                `json:"day_night"`
	Scheduled    string                `json:"scheduled"`
	HomeTeam     string                `json:"home_team"`
	AwayTeam     string                `json:"away_team"`
	DoubleHeader bool                  `json:"double_header"`
	EntryMode    string                `json:"entry_mode"`
	Reference    string                `json:"reference"`
	Venue        sportradar.Venue      `json:"venue"`
	Broadcast    sportradar.Broadcast  `json:"broadcast"`
	Weather      sportradar.Weather    `json:"weather"`
	Home         DailySummaryTeam      `json:"home"`
	Away         DailySummaryTeam      `json:"away"`
}

func srDailySummaryGameConvert(d srDailySummaryGame) DailySummaryGame {
	return DailySummaryGame{
		ID:           d.ID,
		Status:       sportradar.ParseGameStatus(d.Status),
		Coverage:     d.Coverage,
		GameNumber:   d.GameNumber,
		DayNight:     d.DayNight,
		Scheduled:    d.Scheduled,
		HomeTeam:     d.HomeTeam,
		AwayTeam:     d.AwayTeam,
		DoubleHeader: d.DoubleHeader,
		EntryMode:    d.EntryMode,
		Reference:    d.Reference,
		Venue:        d.Venue,
		Broadcast:    d.Broadcast,
		Weather:      d.Weather,
		Home:         srDailySummaryTeamConvert(d.Home),
		Away:         srDailySummaryTeamConvert(d.Away),
	}
}

func srDailySummaryGameBatchConvert(b []srDailySummaryGameMap) []DailySummaryGame {
	out := make([]DailySummaryGame, 0, len(b))
	for _, d := range b {
		out = append(out, srDailySummaryGameConvert(d.Game))
	}
	return out
}

type srDailySummaryTeam struct {
	Name   string `json:"name"`
	Market string `json:"market"`
	Abbr   string `json:"abbr"`
	ID     string `json:"id"`
	Runs   int    `json:"runs"`
	Hits   int    `json:"hits"`
	Errors int    `json:"errors"`
	Win    int    `json:"win"`
	Loss   int    `json:"loss"`
}

type DailySummaryTeam struct {
	Name   string `json:"name"`
	Market string `json:"market"`
	Abbr   string `json:"abbr"`
	ID     string `json:"id"`
	Runs   int    `json:"runs"`
	Hits   int    `json:"hits"`
	Errors int    `json:"errors"`
	Win    int    `json:"win"`
	Loss   int    `json:"loss"`
}

func srDailySummaryTeamConvert(d srDailySummaryTeam) DailySummaryTeam {
	return DailySummaryTeam{
		Name:   d.Name,
		Market: d.Market,
		Abbr:   d.Abbr,
		ID:     d.ID,
		Runs:   d.Runs,
		Hits:   d.Hits,
		Errors: d.Errors,
		Win:    d.Win,
		Loss:   d.Loss,
	}
}

type srDailySummaryLeague struct {
	Alias string                  `json:"alias"`
	Name  string                  `json:"name"`
	ID    string                  `json:"id"`
	Date  string                  `json:"date"`
	Games []srDailySummaryGameMap `json:"games"`
}

type DailySummaryLeague struct {
	Alias string             `json:"alias"`
	Name  string             `json:"name"`
	ID    string             `json:"id"`
	Date  string             `json:"date"`
	Games []DailySummaryGame `json:"games"`
}

func srDailySummaryLeagueConvert(d srDailySummaryLeague) DailySummaryLeague {
	return DailySummaryLeague{
		Alias: d.Alias,
		Name:  d.Name,
		ID:    d.ID,
		Date:  d.Date,
		Games: srDailySummaryGameBatchConvert(d.Games),
	}
}

type srDailySummary struct {
	League  srDailySummaryLeague `json:"league"`
	Comment string               `json:"_comment"`
}

type DailySummary struct {
	League  DailySummaryLeague `json:"league"`
	Comment string             `json:"_comment"`
}

func srDailySummaryConvert(d srDailySummary) DailySummary {
	return DailySummary{
		League:  srDailySummaryLeagueConvert(d.League),
		Comment: d.Comment,
	}
}

var sportradarDailySummaryURLTemplate string = "http://api.sportradar.us/mlb/production/v7/en/games/%d/%d/%d/summary.json?api_key=%s" // year, month, day, apiKey

func dailySummaryURL(t time.Time, apiKey string) string {
	return fmt.Sprintf(sportradarDailySummaryURLTemplate, t.Year(), t.Month(), t.Day(), apiKey)
}

func FetchDailySummary(c http.Client, t time.Time, apiKey string) (DailySummary, error) {
	url := dailySummaryURL(t, apiKey)
	// fmt.Printf("daily summary url: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return DailySummary{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return DailySummary{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return DailySummary{}, fmt.Errorf("invalid response code from sport radar game summary request: %s %d", url, resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DailySummary{}, err
	}

	summary := srDailySummary{}
	err = json.Unmarshal(respBytes, &summary)
	if err != nil {
		return DailySummary{}, err
	}

	return srDailySummaryConvert(summary), nil
}
