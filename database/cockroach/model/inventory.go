package model

import (
	"time"

	"github.com/google/uuid"
)

type Building struct {
	ID         uuid.UUID
	Blueprint  uuid.UUID
	Owner      uuid.UUID
	CreatedAt  time.Time
	FinishedAt time.Time
	Active     bool
}
