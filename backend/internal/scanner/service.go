package scanner

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	registry    *Registry
	mu          sync.RWMutex
	scans       map[string]*ScanRecord
	subscribers map[string][]chan Event
	counter     atomic.Uint64
}

func NewService(registry *Registry) *Service {
	return &Service{
		registry:    registry,
		scans:       make(map[string]*ScanRecord),
		subscribers: make(map[string][]chan Event),
	}
}

func (s *Service) Start(req ScanRequest) (*ScanRecord, error) {
	if err := ValidateTarget(req.TargetURL); err != nil {
		record := &ScanRecord{
			ID:        s.nextID(),
			TargetURL: req.TargetURL,
			Status:    StatusRestricted,
			Grade:     "F",
			Summary:   "Hedef SSRF koruması nedeniyle reddedildi.",
			Error:     err.Error(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.mu.Lock()
		s.scans[record.ID] = record
		s.mu.Unlock()
		return record, nil
	}

	record := &ScanRecord{
		ID:         s.nextID(),
		TargetURL:  req.TargetURL,
		Status:     StatusQueued,
		Grade:      "Pending",
		Progress:   0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Modules:    []ModuleResult{},
		Meta:       map[string]interface{}{"module_count": len(s.registry.Resolve(req.Modules))},
	}

	s.mu.Lock()
	s.scans[record.ID] = record
	s.mu.Unlock()

	go s.runScan(record.ID, req)
	return record, nil
}

func (s *Service) Get(id string) (*ScanRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.scans[id]
	if !ok {
		return nil, false
	}

	copyRecord := *record
	copyRecord.Modules = append([]ModuleResult{}, record.Modules...)
	copyRecord.Recommendations = append([]string{}, record.Recommendations...)
	return &copyRecord, true
}

func (s *Service) Subscribe(id string) (<-chan Event, func()) {
	ch := make(chan Event, 10)
	s.mu.Lock()
	s.subscribers[id] = append(s.subscribers[id], ch)
	s.mu.Unlock()

	cancel := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		list := s.subscribers[id]
		for idx, item := range list {
			if item == ch {
				s.subscribers[id] = append(list[:idx], list[idx+1:]...)
				break
			}
		}
		close(ch)
	}

	return ch, cancel
}

func (s *Service) nextID() string {
	return fmt.Sprintf("scan-%04d", s.counter.Add(1))
}

func (s *Service) runScan(id string, req ScanRequest) {
	scanners := s.registry.Resolve(req.Modules)
	if len(scanners) == 0 {
		s.update(id, func(record *ScanRecord) {
			record.Status = StatusFailed
			record.Grade = "F"
			record.Error = "no valid scanner modules requested"
		})
		return
	}

	s.update(id, func(record *ScanRecord) {
		record.Status = StatusRunning
		record.Summary = "Tarama başlatıldı."
	})
	s.publish(id, Event{Type: "status", Message: "scan started"})

	results := make(chan ModuleResult, len(scanners))
	var wg sync.WaitGroup

	for _, mod := range scanners {
		wg.Add(1)
		go func(mod Scanner) {
			defer wg.Done()
			results <- mod.Run(req.TargetURL)
		}(mod)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	completed := 0
	moduleResults := make([]ModuleResult, 0, len(scanners))
	recommendations := make([]string, 0, len(scanners))
	for result := range results {
		completed++
		moduleResults = append(moduleResults, result)
		recommendations = append(recommendations, result.Findings...)
		progress := completed * 100 / len(scanners)

		s.update(id, func(record *ScanRecord) {
			record.Modules = append([]ModuleResult{}, moduleResults...)
			record.Progress = progress
			record.UpdatedAt = time.Now()
		})
		s.publish(id, Event{Type: "progress", Message: fmt.Sprintf("%s tamamlandı", result.Name), Data: result})
	}

	s.update(id, func(record *ScanRecord) {
		record.Status = StatusCompleted
		record.Progress = 100
		record.Recommendations = uniqueNonEmpty(recommendations)
		record.Grade = deriveGrade(record.Modules)
		record.Summary = buildSummary(record.Modules)
		record.UpdatedAt = time.Now()
	})
	s.publish(id, Event{Type: "completed", Message: "scan completed"})
}

func (s *Service) update(id string, fn func(record *ScanRecord)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if record, ok := s.scans[id]; ok {
		fn(record)
		record.UpdatedAt = time.Now()
	}
}

func (s *Service) publish(id string, event Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sub := range s.subscribers[id] {
		select {
		case sub <- event:
		default:
		}
	}
}

func deriveGrade(results []ModuleResult) string {
	score := 100
	for _, result := range results {
		switch result.Grade {
		case "A+":
			score -= 0
		case "A":
			score -= 4
		case "B":
			score -= 10
		case "C":
			score -= 18
		case "D":
			score -= 28
		default:
			score -= 35
		}
	}

	switch {
	case score >= 92:
		return "A+"
	case score >= 84:
		return "A"
	case score >= 72:
		return "B"
	case score >= 60:
		return "C"
	case score >= 45:
		return "D"
	default:
		return "F"
	}
}

func buildSummary(results []ModuleResult) string {
	if len(results) == 0 {
		return "Henüz modül sonucu yok."
	}

	weak := 0
	for _, result := range results {
		if result.Grade == "D" || result.Grade == "F" {
			weak++
		}
	}

	if weak == 0 {
		return "Tarama tamamlandı. Kritik seviyede zayıf modül tespit edilmedi."
	}

	return fmt.Sprintf("Tarama tamamlandı. %d modül düşük not aldı ve öncelikli iyileştirme gerektiriyor.", weak)
}

func uniqueNonEmpty(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

