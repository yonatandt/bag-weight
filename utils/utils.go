package utils

import "sync"

// SafeHistoryMap is a struct that wraps a int->bool map
// with a RW Mutex Lock.
type SafeHistoryMap struct {
	HistoryMap map[int]bool
	Lock       sync.RWMutex
}
