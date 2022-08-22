package sportradar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type MLBGameSummary struct {
	ID           string         `json:"id"`
	Status       string         `json:"status"`
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

type MLBLeagueSummary struct {
	Alias string           `json:"alias"`
	Name  string           `json:"name"`
	ID    string           `json:"id"`
	Date  string           `json:"date"`
	Games []MLBGameSummary `json:"games"`
}

type MLBDailySummary struct {
	League  MLBLeagueSummary `json:"league"`
	Comment string           `json:"_comment"`
}

var sportradarMLBDailySummaryURLTemplate string = "http://api.sportradar.us/mlb/trial/v7/en/games/%d/%d/%d/summary.json?api_key=%s" // year, month, day, apiKey

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
	summary := MLBDailySummary{}
	err = json.Unmarshal(respBytes, &summary)
	if err != nil {
		return MLBDailySummary{}, nil
	}

	return summary, nil
}
