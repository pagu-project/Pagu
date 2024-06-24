package job

import (
	"sync"
)

type Scheduler interface {
	Submit(job Job)
	Run()
	Shutdown()
}

type scheduler struct {
	jobs []Job
	mu   *sync.Mutex
	wg   *sync.WaitGroup
}

func NewScheduler() Scheduler {
	return &scheduler{
		jobs: make([]Job, 0),
		mu:   &sync.Mutex{},
		wg:   &sync.WaitGroup{},
	}
}

func (s *scheduler) Submit(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *scheduler) Run() {
	for _, j := range s.jobs {
		s.wg.Add(1)
		go j.Start()
	}
}

func (s *scheduler) Shutdown() {
	for _, j := range s.jobs {
		j.Stop()
		s.wg.Done()
	}
}
