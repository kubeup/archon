package pbtime

import (
	"fmt"

	"go.pedge.io/pb/go/google/protobuf"
)

// Duration converts t to a google_protobuf.Duration.
func (t *TimestampRange) Duration() (*google_protobuf.Duration, error) {
	if t.Start == nil {
		return nil, fmt.Errorf("pbtime: %v has no start timestamp", t)
	}
	if t.End == nil {
		return nil, fmt.Errorf("pbtime: %v has no end timestamp", t)
	}
	if !t.Start.Before(t.End) {
		return nil, fmt.Errorf("pbtime: %v has a start before end", t)
	}
	return google_protobuf.DurationToProto(t.End.GoTime().Sub(t.Start.GoTime())), nil
}
