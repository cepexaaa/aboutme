package infrastructure

import (
	"fmt"
	"sync"
	"time"
)

type Progress struct {
	total       int
	current     int64
	lastPercent int
	startTime   time.Time
	logger      *Logger
	mutex       sync.Mutex
}

func NewProgress(total int, logger *Logger) *Progress {
	return &Progress{
		total:     total,
		startTime: time.Now(),
		logger:    logger,
	}
}

func (p *Progress) Update(count ...int64) {
	var cnt int64
	if len(count) == 0 {
		cnt = 1
	} else {
		cnt = count[0]
	}

	p.mutex.Lock()
	defer func() { p.mutex.Unlock() }()

	p.current += cnt
	percent := int(float64(p.current) / float64(p.total) * 100)

	if percent > p.lastPercent && percent%10 == 0 {
		elapsed := time.Since(p.startTime)
		estimatedTotal := time.Duration(float64(elapsed) / float64(percent) * 100)
		remaining := estimatedTotal - elapsed

		p.logger.Info(fmt.Sprintf("Progress: %d%% (%d/%d), Elapsed: %v, Remaining: ~%v",
			percent, p.current, p.total, elapsed.Round(time.Second), remaining.Round(time.Second)))

		p.lastPercent = percent
	}
}

func (p *Progress) Finish() {
	elapsed := time.Since(p.startTime)
	p.logger.Info(fmt.Sprintf("Generation completed in %v", elapsed.Round(time.Millisecond)))
}
