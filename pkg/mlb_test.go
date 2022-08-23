package sportradar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const mlbKey = "cms9sf848zh9tnptnp62m3ts"

func TestFetchMLBDailySummary(t *testing.T) {
	now := time.Now()
	summary, err := FetchMLBDailySummary(now, mlbKey)
	assert.NoError(t, err)
	assert.Equal(t, "MLB", summary.League.Alias)
	assert.True(t, len(summary.League.Games) > 1)
	assert.NotEmpty(t, summary.League.Games[0].ID)
}
