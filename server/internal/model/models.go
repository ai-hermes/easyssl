package model

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type AuthContext struct {
	UserID   string `json:"userId"`
	Role     string `json:"role"`
	AuthType string `json:"authType"`
	APIKeyID string `json:"apiKeyId,omitempty"`
}

func (a AuthContext) IsAdmin() bool {
	return a.Role == RoleAdmin
}

type Admin struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Access struct {
	ID          string                 `json:"id"`
	OwnerUserID string                 `json:"ownerUserId"`
	Name        string                 `json:"name"`
	Provider    string                 `json:"provider"`
	Config      map[string]interface{} `json:"config"`
	Reserve     string                 `json:"reserve,omitempty"`
	DeletedAt   *time.Time             `json:"deleted,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type Certificate struct {
	ID               string     `json:"id"`
	OwnerUserID      string     `json:"ownerUserId"`
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
	OwnerUserID   string                 `json:"ownerUserId"`
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
	ID          string                 `json:"id"`
	OwnerUserID string                 `json:"ownerUserId"`
	WorkflowID  string                 `json:"workflowId"`
	Status      string                 `json:"status"`
	Trigger     string                 `json:"trigger"`
	StartedAt   time.Time              `json:"startedAt"`
	EndedAt     *time.Time             `json:"endedAt"`
	Graph       map[string]interface{} `json:"graph"`
	Error       string                 `json:"error"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type WorkflowRunNode struct {
	ID        string                 `json:"id"`
	RunID     string                 `json:"runId"`
	NodeID    string                 `json:"nodeId"`
	NodeName  string                 `json:"nodeName"`
	Action    string                 `json:"action"`
	Provider  string                 `json:"provider"`
	Status    string                 `json:"status"`
	StartedAt *time.Time             `json:"startedAt"`
	EndedAt   *time.Time             `json:"endedAt"`
	Error     string                 `json:"error"`
	Output    map[string]interface{} `json:"output"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type WorkflowRunEvent struct {
	ID        string                 `json:"id"`
	RunID     string                 `json:"runId"`
	NodeID    string                 `json:"nodeId"`
	EventType string                 `json:"eventType"`
	Message   string                 `json:"message"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"createdAt"`
}

type APIKey struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	KeyHash    string     `json:"-"`
	Status     string     `json:"status"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty"`
}

type Statistics struct {
	CertificateTotal        int `json:"certificateTotal"`
	CertificateExpiringSoon int `json:"certificateExpiringSoon"`
	CertificateExpired      int `json:"certificateExpired"`
	WorkflowTotal           int `json:"workflowTotal"`
	WorkflowEnabled         int `json:"workflowEnabled"`
	WorkflowDisabled        int `json:"workflowDisabled"`
}

type OpenApplyCertificateRequest struct {
	Provider              string   `json:"provider"`
	AccessID              string   `json:"accessId"`
	Domains               []string `json:"domains"`
	CAProvider            string   `json:"caProvider"`
	ContactEmail          string   `json:"contactEmail"`
	KeyAlgorithm          string   `json:"keyAlgorithm"`
	DNSPropagationTimeout int      `json:"dnsPropagationTimeout"`
	DNSTTL                int      `json:"dnsTTL"`
}
