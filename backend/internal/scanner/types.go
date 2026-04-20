package scanner

import "time"

type ModuleName string

const (
	ModulePorts   ModuleName = "ports"
	ModuleHeaders ModuleName = "headers"
	ModuleTLS     ModuleName = "tls"
	ModuleFuzz    ModuleName = "fuzz"
	ModuleXSS     ModuleName = "xss"
	ModuleSQLi    ModuleName = "sqli"
	ModuleCVE     ModuleName = "cve"
)

var DefaultModules = []ModuleName{
	ModulePorts,
	ModuleHeaders,
	ModuleTLS,
	ModuleFuzz,
	ModuleXSS,
	ModuleSQLi,
	ModuleCVE,
}

type ScanStatus string

const (
	StatusQueued     ScanStatus = "queued"
	StatusRunning    ScanStatus = "running"
	StatusCompleted  ScanStatus = "completed"
	StatusFailed     ScanStatus = "failed"
	StatusRestricted ScanStatus = "restricted"
)

type ScanRequest struct {
	TargetURL string       `json:"target_url"`
	Modules   []ModuleName `json:"modules"`
}

type ModuleResult struct {
	Name      ModuleName      `json:"name"`
	Status    string          `json:"status"`
	Grade     string          `json:"grade,omitempty"`
	Summary   string          `json:"summary"`
	OWASP     []string        `json:"owasp,omitempty"`
	Findings  []string        `json:"findings,omitempty"`
	Metadata  map[string]any  `json:"metadata,omitempty"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
}

type ScanRecord struct {
	ID            string                 `json:"id"`
	TargetURL     string                 `json:"target_url"`
	Status        ScanStatus             `json:"status"`
	Grade         string                 `json:"grade"`
	Progress      int                    `json:"progress"`
	Modules       []ModuleResult         `json:"modules"`
	Summary       string                 `json:"summary"`
	Recommendations []string             `json:"recommendations"`
	Error         string                 `json:"error,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Meta          map[string]interface{} `json:"meta,omitempty"`
}

type Event struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Scanner interface {
	Name() ModuleName
	Run(target string) ModuleResult
}

