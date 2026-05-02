package workflow

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"easyssl/server/internal/accessprovider"
	"easyssl/server/internal/certacme"
	"easyssl/server/internal/deployer"
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

		if err := d.executeRun(taskCtx, runID); err != nil {
			if taskCtx.Err() != nil {
				_ = d.repo.UpdateWorkflowRunStatus(context.Background(), runID, "canceled", "")
				return
			}
			_ = d.repo.UpdateWorkflowRunStatus(context.Background(), runID, "failed", err.Error())
			return
		}
		_ = d.repo.UpdateWorkflowRunStatus(context.Background(), runID, "succeeded", "")
	}()
}

type graphNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	} `json:"data"`
}

type graphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type workflowGraph struct {
	Nodes []graphNode `json:"nodes"`
	Edges []graphEdge `json:"edges"`
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func readProvider(config map[string]interface{}) string {
	return strings.ToLower(toString(config["provider"]))
}

func normalizeDNSProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "aliyun", "alidns":
		return accessprovider.ProviderAliyun
	case "tencent", "tencentcloud", "dnspod":
		return accessprovider.ProviderTencentCloud
	default:
		return strings.ToLower(strings.TrimSpace(provider))
	}
}

func readAccessID(config map[string]interface{}) string {
	for _, key := range []string{"accessId", "accessID", "access_id", "providerAccessId"} {
		if id := toString(config[key]); id != "" {
			return id
		}
	}
	return ""
}

func inferNodeAction(node graphNode) string {
	t := strings.ToLower(strings.TrimSpace(node.Type))
	id := strings.ToLower(strings.TrimSpace(node.ID))
	name := strings.ToLower(strings.TrimSpace(node.Data.Name))
	provider := readProvider(node.Data.Config)

	if t == "input" || t == "output" || t == "start" || t == "end" {
		return ""
	}
	if id == "start" || id == "end" || name == "start" || name == "end" {
		return ""
	}

	if t == "apply" || t == "bizapply" || strings.Contains(id, "apply") || strings.Contains(name, "apply") {
		return "apply"
	}
	if t == "deploy" || t == "bizdeploy" || strings.Contains(id, "deploy") || strings.Contains(name, "deploy") {
		return "deploy"
	}

	if provider == accessprovider.ProviderAliyun {
		if _, ok := node.Data.Config["domains"]; ok {
			return "apply"
		}
	}
	if provider == deployer.ProviderAliyunCAS || provider == deployer.ProviderQiniu || provider == deployer.ProviderSSH {
		return "deploy"
	}
	return ""
}

func parseGraphOrder(graph workflowGraph) []graphNode {
	if len(graph.Nodes) <= 1 || len(graph.Edges) == 0 {
		return graph.Nodes
	}

	inDegree := map[string]int{}
	nodeByID := map[string]graphNode{}
	adj := map[string][]string{}
	indexByID := map[string]int{}
	for i, n := range graph.Nodes {
		nodeByID[n.ID] = n
		indexByID[n.ID] = i
		if _, ok := inDegree[n.ID]; !ok {
			inDegree[n.ID] = 0
		}
	}
	for _, e := range graph.Edges {
		if _, ok := nodeByID[e.Source]; !ok {
			continue
		}
		if _, ok := nodeByID[e.Target]; !ok {
			continue
		}
		adj[e.Source] = append(adj[e.Source], e.Target)
		inDegree[e.Target]++
	}

	queue := make([]string, 0)
	for _, n := range graph.Nodes {
		if inDegree[n.ID] == 0 {
			queue = append(queue, n.ID)
		}
	}
	slices.SortFunc(queue, func(a, b string) int {
		return indexByID[a] - indexByID[b]
	})

	ordered := make([]graphNode, 0, len(graph.Nodes))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		ordered = append(ordered, nodeByID[cur])
		for _, next := range adj[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(ordered) != len(graph.Nodes) {
		return graph.Nodes
	}
	return ordered
}

func parseDomains(v interface{}) []string {
	if v == nil {
		return nil
	}
	domains := make([]string, 0)
	switch vv := v.(type) {
	case string:
		s := strings.NewReplacer("\n", ";", ",", ";", " ", ";").Replace(vv)
		for _, p := range strings.Split(s, ";") {
			p = strings.TrimSpace(p)
			if p != "" {
				domains = append(domains, p)
			}
		}
	case []interface{}:
		for _, item := range vv {
			s := strings.TrimSpace(fmt.Sprintf("%v", item))
			if s != "" {
				domains = append(domains, s)
			}
		}
	}
	return domains
}

func parseInt(v interface{}, fallback int) int {
	s := toString(v)
	if s == "" {
		return fallback
	}
	var i int
	_, _ = fmt.Sscanf(s, "%d", &i)
	if i == 0 {
		return fallback
	}
	return i
}

func (d *Dispatcher) executeRun(ctx context.Context, runID string) error {
	run, err := d.repo.GetWorkflowRun(ctx, runID)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(run.Graph)
	if err != nil {
		return fmt.Errorf("marshal graph failed: %w", err)
	}
	var graph workflowGraph
	if err := json.Unmarshal(raw, &graph); err != nil {
		return fmt.Errorf("parse graph failed: %w", err)
	}

	nodes := parseGraphOrder(graph)
	certByNodeID := map[string]*model.Certificate{}
	var latestCert *model.Certificate
	applyCount := 0

	upsertNode := func(node graphNode, action, provider, status string, startedAt, endedAt *time.Time, errMsg string, output map[string]interface{}) {
		_, _ = d.repo.UpsertWorkflowRunNode(context.Background(), model.WorkflowRunNode{
			RunID:     run.ID,
			NodeID:    node.ID,
			NodeName:  node.Data.Name,
			Action:    action,
			Provider:  provider,
			Status:    status,
			StartedAt: startedAt,
			EndedAt:   endedAt,
			Error:     errMsg,
			Output:    output,
		})
	}
	appendEvent := func(nodeID, eventType, message string, payload map[string]interface{}) {
		_, _ = d.repo.AppendWorkflowRunEvent(context.Background(), model.WorkflowRunEvent{
			RunID:     run.ID,
			NodeID:    nodeID,
			EventType: eventType,
			Message:   message,
			Payload:   payload,
		})
	}

	for _, node := range nodes {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		action := inferNodeAction(node)
		cfg := node.Data.Config
		if cfg == nil || action == "" {
			continue
		}
		nodeProvider := readProvider(cfg)
		if action == "apply" {
			nodeProvider = normalizeDNSProvider(nodeProvider)
		}
		startedAt := time.Now()
		upsertNode(node, action, nodeProvider, "running", &startedAt, nil, "", nil)
		appendEvent(node.ID, "started", "node started", map[string]interface{}{"action": action, "provider": nodeProvider})

		failNode := func(err error) error {
			endedAt := time.Now()
			upsertNode(node, action, nodeProvider, "failed", &startedAt, &endedAt, err.Error(), nil)
			appendEvent(node.ID, "failed", err.Error(), nil)
			return err
		}

		switch action {
		case "apply":
			nodeProvider = normalizeDNSProvider(readProvider(cfg))
			accessID := readAccessID(cfg)
			if accessID == "" {
				return failNode(fmt.Errorf("apply: node %s requires config.accessId", node.ID))
			}
			access, err := d.repo.GetAccessByID(ctx, accessID)
			if err != nil {
				return failNode(fmt.Errorf("apply: load access %s failed: %w", accessID, err))
			}
			accessProvider := normalizeDNSProvider(access.Provider)
			if nodeProvider == "" {
				nodeProvider = accessProvider
			}
			if nodeProvider != accessProvider {
				return failNode(fmt.Errorf("apply: node provider %s does not match access provider %s", nodeProvider, accessProvider))
			}
			if nodeProvider != accessprovider.ProviderAliyun && nodeProvider != accessprovider.ProviderTencentCloud {
				return failNode(fmt.Errorf("apply: unsupported provider %s", nodeProvider))
			}
			appendEvent(node.ID, "log", "dns provider validated", map[string]interface{}{"provider": nodeProvider})

			domains := parseDomains(cfg["domains"])
			if len(domains) == 0 {
				return failNode(fmt.Errorf("apply: node %s requires non-empty domains", node.ID))
			}
			appendEvent(node.ID, "log", "requesting certificate", map[string]interface{}{"domains": domains})

			acmeResp, err := certacme.Obtain(certacme.ObtainRequest{
				Email:                 toString(cfg["contactEmail"]),
				Domains:               domains,
				CAProvider:            toString(cfg["caProvider"]),
				DNSProvider:           nodeProvider,
				AliyunAccessKeyID:     toString(access.Config["accessKeyId"]),
				AliyunAccessKeySecret: toString(access.Config["accessKeySecret"]),
				TencentSecretID:       toString(access.Config["secretId"]),
				TencentSecretKey:      toString(access.Config["secretKey"]),
				TencentRegion:         toString(access.Config["region"]),
				TencentSessionToken:   toString(access.Config["sessionToken"]),
				DNSPropagationTimeout: parseInt(cfg["dnsPropagationTimeout"], 120),
				DNSTTL:                parseInt(cfg["dnsTTL"], 0),
				KeyAlgorithm:          toString(cfg["keyAlgorithm"]),
			})
			if err != nil {
				return failNode(fmt.Errorf("apply: node %s request certificate failed: %w", node.ID, err))
			}

			parsed, err := parseCertificateMeta(acmeResp.Certificate, acmeResp.PrivateKey)
			if err != nil {
				return failNode(fmt.Errorf("apply: parse certificate failed: %w", err))
			}

			saved, err := d.repo.SaveCertificate(ctx, model.Certificate{
				Source:           "request",
				SubjectAltNames:  parsed.SubjectAltNames,
				SerialNumber:     parsed.SerialNumber,
				Certificate:      acmeResp.Certificate,
				PrivateKey:       acmeResp.PrivateKey,
				IssuerOrg:        parsed.IssuerOrg,
				KeyAlgorithm:     parsed.KeyAlgorithm,
				ValidityNotAfter: parsed.ValidityNotAfter,
				WorkflowID:       run.WorkflowID,
				WorkflowRunID:    run.ID,
			})
			if err != nil {
				return failNode(fmt.Errorf("apply: save certificate failed: %w", err))
			}
			certByNodeID[node.ID] = saved
			latestCert = saved
			applyCount++
			endedAt := time.Now()
			output := map[string]interface{}{
				"certificateId": saved.ID,
				"domains":       parseDomains(cfg["domains"]),
				"notAfter":      parsed.ValidityNotAfter,
			}
			upsertNode(node, action, nodeProvider, "succeeded", &startedAt, &endedAt, "", output)
			appendEvent(node.ID, "succeeded", "certificate obtained", output)

		case "deploy":
			if latestCert == nil {
				return failNode(fmt.Errorf("deploy: node %s requires an apply node before it", node.ID))
			}
			provider := readProvider(cfg)
			accessID := readAccessID(cfg)
			if accessID == "" {
				return failNode(fmt.Errorf("deploy: node %s requires config.accessId", node.ID))
			}
			access, err := d.repo.GetAccessByID(ctx, accessID)
			if err != nil {
				return failNode(fmt.Errorf("deploy: load access %s failed: %w", accessID, err))
			}
			if provider == "" {
				provider = strings.ToLower(strings.TrimSpace(access.Provider))
			}
			if provider == accessprovider.ProviderAliyun {
				provider = deployer.ProviderAliyunCAS
			}
			nodeProvider = provider
			appendEvent(node.ID, "log", "deploy provider resolved", map[string]interface{}{"provider": provider})

			certToDeploy := latestCert
			if fromNode := toString(cfg["certificateOutputNodeId"]); fromNode != "" {
				if c, ok := certByNodeID[fromNode]; ok {
					certToDeploy = c
				} else {
					return failNode(fmt.Errorf("deploy: node %s references unknown certificateOutputNodeId %s", node.ID, fromNode))
				}
			}

			if err := deployer.Execute(ctx, deployer.Request{
				Provider:     provider,
				AccessConfig: access.Config,
				Config:       cfg,
				Certificate:  certToDeploy.Certificate,
				PrivateKey:   certToDeploy.PrivateKey,
			}); err != nil {
				return failNode(fmt.Errorf("deploy: node %s failed: %w", node.ID, err))
			}

			endedAt := time.Now()
			output := map[string]interface{}{
				"provider":      nodeProvider,
				"certificateId": certToDeploy.ID,
			}
			upsertNode(node, action, nodeProvider, "succeeded", &startedAt, &endedAt, "", output)
			appendEvent(node.ID, "succeeded", "certificate deployed", output)
		}
	}

	if applyCount == 0 {
		return fmt.Errorf("workflow has no apply node")
	}
	return nil
}

type certMeta struct {
	SubjectAltNames  string
	SerialNumber     string
	IssuerOrg        string
	KeyAlgorithm     string
	ValidityNotAfter *time.Time
}

func parseCertificateMeta(certPEM, keyPEM string) (*certMeta, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid certificate pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	serial := serialToHex(cert.SerialNumber)

	sans := make([]string, 0, len(cert.DNSNames)+1)
	if cert.Subject.CommonName != "" {
		sans = append(sans, cert.Subject.CommonName)
	}
	sans = append(sans, cert.DNSNames...)
	sans = dedupStrings(sans)

	issuer := ""
	if len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	}

	alg := inferKeyAlgorithm(keyPEM)
	notAfter := cert.NotAfter
	return &certMeta{
		SubjectAltNames:  strings.Join(sans, ";"),
		SerialNumber:     serial,
		IssuerOrg:        issuer,
		KeyAlgorithm:     alg,
		ValidityNotAfter: &notAfter,
	}, nil
}

func serialToHex(n *big.Int) string {
	if n == nil {
		return ""
	}
	return strings.ToUpper(n.Text(16))
}

func dedupStrings(items []string) []string {
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

func inferKeyAlgorithm(privateKeyPEM string) string {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return ""
	}
	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		switch key := pk.(type) {
		case *rsa.PrivateKey:
			return fmt.Sprintf("RSA%d", key.N.BitLen())
		case *ecdsa.PrivateKey:
			if key.Curve != nil {
				return fmt.Sprintf("EC%d", key.Curve.Params().BitSize)
			}
			return "EC"
		}
	}
	if rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return fmt.Sprintf("RSA%d", rsaKey.N.BitLen())
	}
	if ecKey, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		if ecKey.Curve != nil {
			return fmt.Sprintf("EC%d", ecKey.Curve.Params().BitSize)
		}
		return "EC"
	}
	return ""
}
