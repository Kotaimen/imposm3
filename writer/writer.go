package writer

import (
	"imposm3/cache"
	"imposm3/database"
	"imposm3/expire"
	"imposm3/geom/limit"
	"imposm3/logging"
	"imposm3/stats"
	"runtime"
	"sync"
)

var log = logging.NewLogger("writer")

type ErrorLevel interface {
	Level() int
}

type looper interface {
	loop()
}

type OsmElemWriter struct {
	osmCache   *cache.OSMCache
	diffCache  *cache.DiffCache
	progress   *stats.Statistics
	inserter   database.Inserter
	wg         *sync.WaitGroup
	limiter    *limit.Limiter
	writer     looper
	srid       int
	expireor   expire.Expireor
	concurrent bool
}

func (writer *OsmElemWriter) SetLimiter(limiter *limit.Limiter) {
	writer.limiter = limiter
}

func (writer *OsmElemWriter) EnableConcurrent() {
	writer.concurrent = true
}

func (writer *OsmElemWriter) Start() {
	concurrency := 1
	if writer.concurrent {
		concurrency = runtime.NumCPU()
	}
	for i := 0; i < concurrency; i++ {
		writer.wg.Add(1)
		go writer.writer.loop()
	}
}

func (writer *OsmElemWriter) SetExpireor(exp expire.Expireor) {
	writer.expireor = exp
}

func (writer *OsmElemWriter) Wait() {
	writer.wg.Wait()
}
