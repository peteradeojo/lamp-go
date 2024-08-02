package tasks

import (
	"github.com/google/uuid"
)

type TaskComment struct {
	Message  string
	SenderId int
	TaskId   uuid.UUID
}
