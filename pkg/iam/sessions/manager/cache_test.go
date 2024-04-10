package manager

import (
	"testing"

	"github.com/codfrm/cago/database/cache/memory"
)

func TestNewCacheSessionManager(t *testing.T) {
	m, _ := memory.NewMemoryCache()
	testExpireSession(t, NewCacheSessionManagerWithExpire(m, 60))
}
