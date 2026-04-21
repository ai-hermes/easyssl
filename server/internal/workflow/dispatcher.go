package workflow

import (
	"context"
	"sync"
	"time"

	"easyssl/server/internal/model"
	"easyssl/server/internal/repository"
)

type Dispatcher struct {
	repo *repository.Repository

	mu         sync.Mutex
	maxWorkers int
	pending    []string
	processing map[string]context.CancelFunc
}

func NewDispatcher(repo *repository.Repository, maxWorkers int) *Dispatcher {
	if maxWorkers <= 0 {
		maxWorkers = 2
	}
	return &Dispatcher{repo: repo, maxWorkers: maxWorkers, pending: make([]string, 0), processing: make(map[string]context.CancelFunc)}
}

func (d *Dispatcher) Stats() (int, []string, []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	pending := append([]string(nil), d.pending...)
	processing := make([]string, 0, len(d.processing))
	for id := range d.processing {
		processing = append(processing, id)
	}
	return d.maxWorkers, pending, processing
}

func (d *Dispatcher) Start(ctx context.Context, runID string) {
	d.mu.Lock()
	d.pending = append(d.pending, runID)
	d.mu.Unlock()
	go d.tryNext(context.Background())
}

func (d *Dispatcher) Cancel(ctx context.Context, runID string) {
	d.mu.Lock()
	if cancel, ok := d.processing[runID]; ok {
		cancel()
		delete(d.processing, runID)
	}
	for i, id := range d.pending {
		if id == runID {
			d.pending = append(d.pending[:i], d.pending[i+1:]...)
			break
		}
	}
	d.mu.Unlock()
	_ = d.repo.UpdateWorkflowRunStatus(ctx, runID, "canceled", "")
}

func (d *Dispatcher) tryNext(ctx context.Context) {
	d.mu.Lock()
	if len(d.processing) >= d.maxWorkers || len(d.pending) == 0 {
		d.mu.Unlock()
		return
	}
	runID := d.pending[0]
	d.pending = d.pending[1:]
	taskCtx, cancel := context.WithCancel(context.Background())
	d.processing[runID] = cancel
	d.mu.Unlock()

	go func() {
		defer func() {
			d.mu.Lock()
			delete(d.processing, runID)
			d.mu.Unlock()
			go d.tryNext(context.Background())
		}()

		_ = d.repo.UpdateWorkflowRunStatus(taskCtx, runID, "processing", "")

		select {
		case <-taskCtx.Done():
			_ = d.repo.UpdateWorkflowRunStatus(context.Background(), runID, "canceled", "")
			return
		case <-time.After(3 * time.Second):
			_ = d.repo.UpdateWorkflowRunStatus(context.Background(), runID, "succeeded", "")
			_ = d.persistMockCertificate(context.Background(), runID)
		}
	}()
}

func (d *Dispatcher) persistMockCertificate(ctx context.Context, runID string) error {
	run, err := d.repo.GetWorkflowRun(ctx, runID)
	if err != nil {
		return err
	}
	now := time.Now()
	expire := now.Add(90 * 24 * time.Hour)
	_, err = d.repo.SaveCertificate(ctx, model.Certificate{
		Source:           "request",
		SubjectAltNames:  "example.com;*.example.com",
		SerialNumber:     runID[:8],
		Certificate:      "-----BEGIN CERTIFICATE-----\nMOCK\n-----END CERTIFICATE-----",
		PrivateKey:       "-----BEGIN PRIVATE KEY-----\nMOCK\n-----END PRIVATE KEY-----",
		IssuerOrg:        "Mock CA",
		KeyAlgorithm:     "RSA2048",
		ValidityNotAfter: &expire,
		WorkflowID:       run.WorkflowID,
		WorkflowRunID:    run.ID,
	})
	return err
}
