package gzstat

import (
	"time"
)

type GzipStatus struct {
	Size              uint64    `json:"size"`
	Modified          time.Time `json:"modified"`
	Comment           string    `json:"comment"`
	Name              string    `json:"name"`
	OriginalSizeMod32 uint32    `json:"original_size_mod32"`
	OS                byte      `json:"os"`
}
