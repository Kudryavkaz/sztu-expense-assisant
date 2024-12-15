package lock

import (
	"strings"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/cache"
	"github.com/go-redsync/redsync/v4"
	"go.uber.org/zap"
)

type Mutex struct {
	lock *redsync.Mutex
}

func GetLock(mutexName string, retries int, outTime time.Duration) (m *Mutex, err error) {
	m = &Mutex{}

	m.lock = cache.Rs.NewMutex("mutex:"+mutexName, redsync.WithExpiry(outTime))

	if retries < 0 {
		retries = 999999
	}

	err = m.lock.Lock()
	for i := 0; i < retries; i++ {
		if err != nil {
			if strings.Contains(err.Error(), "lock already taken") {
				log.Logger().Warn("lock already taken, waiting...", zap.Int("retry", i), zap.String("mutexName", mutexName))
				time.Sleep(1 * time.Second)
				err = m.lock.Lock()
				continue
			}
			return
		}
		break
	}

	return
}

func (m *Mutex) Release() (ok bool, err error) {
	ok, err = m.lock.Unlock()
	return
}
