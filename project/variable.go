package project

import (
	"time"
)

type variable struct {
	Type  DataType
	Value interface{}
	Time  time.Time
}
