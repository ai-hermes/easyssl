package model

import "time"

type Admin struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Access struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Config    map[string]interface{} `json:"config"`
	Reserve   string                 `json:"reserve,omitempty"`
	DeletedAt *time.Time             `json:"deleted,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type Certificate struct {
	ID               string     `json:"id"`
	Source           string     `json:"source"`
	SubjectAltNames  string     `json:"subjectAltNames"`
	SerialNumber     string     `json:"serialNumber"`
	Certificate      string     `json:"certificate"`
	PrivateKey       string     `json:"privateKey"`
	IssuerOrg        string     `json:"issuerOrg"`
	KeyAlgorithm     string     `json:"keyAlgorithm"`
	ValidityNotAfter *time.Time `json:"validityNotAfter"`
	IsRevoked        bool       `json:"isRevoked"`
	WorkflowID       string     `json:"workflowId"`
	WorkflowRunID    string     `json:"workflowRunId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type Workflow struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Trigger       string                 `json:"trigger"`
	TriggerCron   string                 `json:"triggerCron"`
	Enabled       bool                   `json:"enabled"`
	GraphDraft    map[string]interface{} `json:"graphDraft"`
	GraphContent  map[string]interface{} `json:"graphContent"`
	HasDraft      bool                   `json:"hasDraft"`
	HasContent    bool                   `json:"hasContent"`
	LastRunID     *string                `json:"lastRunId"`
	LastRunStatus *string                `json:"lastRunStatus"`
	LastRunTime   *time.Time             `json:"lastRunTime"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

type WorkflowRun struct {
	ID         string                 `json:"id"`
	WorkflowID string                 `json:"workflowId"`
	Status     string                 `json:"status"`
	Trigger    string                 `json:"trigger"`
	StartedAt  time.Time              `json:"startedAt"`
	EndedAt    *time.Time             `json:"endedAt"`
	Graph      map[string]interface{} `json:"graph"`
	Error      string                 `json:"error"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

type Statistics struct {
	CertificateTotal        int `json:"certificateTotal"`
	CertificateExpiringSoon int `json:"certificateExpiringSoon"`
	CertificateExpired      int `json:"certificateExpired"`
	WorkflowTotal           int `json:"workflowTotal"`
	WorkflowEnabled         int `json:"workflowEnabled"`
	WorkflowDisabled        int `json:"workflowDisabled"`
}
