package utils

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/port"
	"time"
)

type uTime struct{}

func (t *uTime) Now() *time.Time {
	now := time.Now()
	return &now
}

func New() port.IUtils {
	return &uTime{}
}
