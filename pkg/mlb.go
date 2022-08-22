package sportradar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type srMLBGameSummary struct {
	ID           string           `json:"id"`
	Status       string           `json:"status"`
	Coverage     string           `json:"coverage"`
	GameNumber   int64            `json:"game_number"`
	DayNight     string           `json:"day_night"`
	Scheduled    string           `json:"scheduled"`
	HomeTeam     string           `json:"home_team"`
	AwayTeam     string           `json:"away_team"`
	DoubleHeader bool             `json:"double_header"`
	EntryMode    string           `json:"entry_mode"`
	Reference    string           `json:"reference"`
	Venue        Venue            `json:"venue"`
	Broadcast    Broadcast        `json:"broadcast"`
	Weather      Weather          `json:"weather"`
	Home         srMLBTeamSummary `json:"home"`
	Away         srMLBTeamSummary `json:"away"`
}

type MLBGameSummary struct {
	ID           string         `json:"id"`
	Status       GameStatus     `json:"status"`
	Coverage     string         `json:"coverage"`
	GameNumber   int64          `json:"game_number"`
	DayNight     string         `json:"day_night"`
	Scheduled    string         `json:"scheduled"`
	HomeTeam     string         `json:"home_team"`
	AwayTeam     string         `json:"away_team"`
	DoubleHeader bool           `json:"double_header"`
	EntryMode    string         `json:"entry_mode"`
	Reference    string         `json:"reference"`
	Venue        Venue          `json:"venue"`
	Broadcast    Broadcast      `json:"broadcast"`
	Weather      Weather        `json:"weather"`
	Home         MLBTeamSummary `json:"home"`
	Away         MLBTeamSummary `json:"away"`
}

func srMLBGameSummaryConvert(d srMLBGameSummary) MLBGameSummary {
	return MLBGameSummary{
		ID:           d.ID,
		Status:       parseGameStatus(d.Status),
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
		Home:         srMLBTeamSummaryConvert(d.Home),
		Away:         srMLBTeamSummaryConvert(d.Away),
	}
}

func srMLBGameSummaryBatchConvert(b []srMLBGameSummary) []MLBGameSummary {
	out := make([]MLBGameSummary, 0, len(b))
	for _, d := range b {
		out = append(out, srMLBGameSummaryConvert(d))
	}
	return out
}

type srMLBTeamSummary struct {
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

type MLBTeamSummary struct {
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

func srMLBTeamSummaryConvert(d srMLBTeamSummary) MLBTeamSummary {
	return MLBTeamSummary{
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

type srMLBLeagueSummary struct {
	Alias string             `json:"alias"`
	Name  string             `json:"name"`
	ID    string             `json:"id"`
	Date  string             `json:"date"`
	Games []srMLBGameSummary `json:"games"`
}

type MLBLeagueSummary struct {
	Alias string           `json:"alias"`
	Name  string           `json:"name"`
	ID    string           `json:"id"`
	Date  string           `json:"date"`
	Games []MLBGameSummary `json:"games"`
}

func srMLBLeagueSummaryConvert(d srMLBLeagueSummary) MLBLeagueSummary {
	return MLBLeagueSummary{
		Alias: d.Alias,
		Name:  d.Name,
		ID:    d.ID,
		Date:  d.Date,
		Games: srMLBGameSummaryBatchConvert(d.Games),
	}
}

type srMLBDailySummary struct {
	League  srMLBLeagueSummary `json:"league"`
	Comment string             `json:"_comment"`
}

type MLBDailySummary struct {
	League  MLBLeagueSummary `json:"league"`
	Comment string           `json:"_comment"`
}

func srMLBDailySummaryConvert(d srMLBDailySummary) MLBDailySummary {
	return MLBDailySummary{
		League:  srMLBLeagueSummaryConvert(d.League),
		Comment: d.Comment,
	}
}

var sportradarMLBDailySummaryURLTemplate string = "http://api.sportradar.us/mlb/production/v7/en/games/%d/%d/%d/summary.json?api_key=%s" // year, month, day, apiKey

func mlbDailySummaryURL(t time.Time, apiKey string) string {
	return fmt.Sprintf(sportradarMLBDailySummaryURLTemplate, t.Year(), t.Month(), t.Day(), apiKey)
}

func FetchMLBDailySummary(t time.Time, apiKey string) (MLBDailySummary, error) {
	url := mlbDailySummaryURL(t, apiKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MLBDailySummary{}, err
	}
	client := http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return MLBDailySummary{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return MLBDailySummary{}, fmt.Errorf("invalid response code from sport radar game summary request: %s %d", url, resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MLBDailySummary{}, err
	}

	log.Info().Msgf("sport radar resp: %s", string(respBytes))
	summary := srMLBDailySummary{}
	err = json.Unmarshal(respBytes, &summary)
	if err != nil {
		return MLBDailySummary{}, nil
	}

	return srMLBDailySummaryConvert(summary), nil
}
