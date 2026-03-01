package application

import (
	"sync"
	"sync/atomic"
	"testing"

	"fractalflame/internal/infrastructure"
)

func BenchmarkParallelIterations(b *testing.B) {
	logger := infrastructure.NewLogger()
	generator := NewGenerator(logger)

	tests := []struct {
		name    string
		threads int
		batch   int
	}{
		{"1-thread-100", 1, 100},
		{"1-thread-1000", 1, 1000},
		{"1-thread-10000", 1, 10000},
		{"4-thread-100", 4, 100},
		{"4-thread-1000", 4, 1000},
		{"4-thread-10000", 4, 10000},
		{"8-thread-100", 8, 100},
		{"8-thread-1000", 8, 1000},
		{"8-thread-10000", 8, 10000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			config := benchmarkConfig(100, 100, tt.batch*100, tt.threads)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := generator.Generate(config)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	tests := []struct {
		name    string
		workers int
		tasks   int
	}{
		{"1-worker-100", 1, 100},
		{"1-worker-1000", 1, 1000},
		{"4-worker-100", 4, 100},
		{"4-worker-1000", 4, 1000},
		{"8-worker-100", 8, 100},
		{"8-worker-1000", 8, 1000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				tasks := make(chan int, tt.tasks)
				results := make(chan int, tt.tasks)

				var wg sync.WaitGroup
				for w := 0; w < tt.workers; w++ {
					wg.Add(1)
					go func(workerID int) {
						defer wg.Done()
						for task := range tasks {

							result := task * task
							results <- result
						}
					}(w)
				}

				for t := 0; t < tt.tasks; t++ {
					tasks <- t
				}
				close(tasks)

				wg.Wait()
				close(results)

				total := 0
				for range results {
					total++
				}

				if total != tt.tasks {
					b.Errorf("Expected %d results, got %d", tt.tasks, total)
				}
			}
		})
	}
}

func BenchmarkAtomicOperations(b *testing.B) {
	tests := []struct {
		name string
		ops  int
	}{
		{"100-ops", 100},
		{"1000-ops", 1000},
		{"10000-ops", 10000},
		{"100000-ops", 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var counter int64

				for j := 0; j < tt.ops; j++ {
					atomic.AddInt64(&counter, 1)
				}

				_ = atomic.LoadInt64(&counter)
			}
		})
	}
}

func BenchmarkMutexVsAtomic(b *testing.B) {
	b.Run("Mutex", func(b *testing.B) {
		var mu sync.Mutex
		var counter int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})

	b.Run("Atomic", func(b *testing.B) {
		var counter int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			atomic.AddInt64(&counter, 1)
		}
	})

	b.Run("RWMutex-Read", func(b *testing.B) {
		var mu sync.RWMutex
		var counter int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.RLock()
			_ = counter
			mu.RUnlock()
		}
	})
}
