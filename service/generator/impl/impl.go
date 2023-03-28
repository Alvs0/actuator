package impl

import (
	"sync"
)

type generator struct {
	mu         sync.Mutex
}
