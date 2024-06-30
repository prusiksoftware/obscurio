package analytics

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAnalytics_TrackQuery(t *testing.T) {
	at := New(60)
	durations := make(map[DurationType]time.Duration)
	durations["query"] = time.Second

	at.TrackQuery("profile", "original", "default", durations)

	require.Equal(t, 1, len(at.Events))
	require.Equal(t, "profile", at.Events[0].Profile)
	require.Equal(t, "original", at.Events[0].OriginalQuery)
	require.Equal(t, "default", at.Events[0].ModifiedQuery)
	require.Equal(t, time.Second, at.Events[0].Durations["query"])
}
