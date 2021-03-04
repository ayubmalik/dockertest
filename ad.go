package dockertestspike

import (
	"github.com/google/uuid"
	"time"
)

type Ad struct {
	ID      uuid.UUID
	Content string
	Created time.Time
	StartAt time.Time
	EndAt   time.Time
}

func NewAd(content string, startAt, endAt time.Time) Ad {
	return Ad{
		ID:      uuid.New(),
		Content: content,
		StartAt: startAt,
		EndAt:   endAt,
		Created: time.Now(),
	}
}
