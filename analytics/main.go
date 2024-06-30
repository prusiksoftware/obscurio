package analytics

import "time"

type DurationType string

type Event struct {
	Profile       string
	OriginalQuery string
	ModifiedQuery string
	Durations     map[DurationType]time.Duration
	Created       time.Time
}

type Analytics struct {
	retentionSeconds int
	Events           []Event
}

func New(retentionSeconds int) *Analytics {
	return &Analytics{
		retentionSeconds: retentionSeconds,
	}
}

func (a *Analytics) TrackQuery(profile, originalQuery, modifiedQuery string, durations map[DurationType]time.Duration) {
	//	enumerate tables and columns accessed
	event := Event{
		Profile:       profile,
		OriginalQuery: originalQuery,
		ModifiedQuery: modifiedQuery,
		Durations:     durations,
		Created:       time.Now(),
	}
	a.Events = append(a.Events, event)
}
