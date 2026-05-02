export type ApiResp<T> = { code: number; msg: string; data: T };

export type Access = {
  id?: string;
  name: string;
  provider: string;
  config: Record<string, unknown>;
  reserve?: string;
};

export type Workflow = {
  id?: string;
  name: string;
  description: string;
  trigger: string;
  triggerCron: string;
  enabled: boolean;
  graphDraft: Record<string, unknown>;
  graphContent: Record<string, unknown>;
  hasDraft: boolean;
  hasContent: boolean;
  lastRunId?: string;
  lastRunStatus?: string;
  lastRunTime?: string;
};

export type WorkflowRun = {
  id: string;
  workflowId: string;
  status: string;
  trigger: string;
  startedAt: string;
  endedAt?: string;
  graph?: Record<string, unknown>;
  error: string;
};

export type WorkflowRunNode = {
  id: string;
  runId: string;
  nodeId: string;
  nodeName: string;
  action: string;
  provider: string;
  status: string;
  startedAt?: string;
  endedAt?: string;
  error: string;
  output?: Record<string, unknown>;
};

export type WorkflowRunEvent = {
  id: string;
  runId: string;
  nodeId: string;
  eventType: string;
  message: string;
  payload?: Record<string, unknown>;
  createdAt: string;
};

export type Certificate = {
  id: string;
  subjectAltNames: string;
  serialNumber: string;
  keyAlgorithm: string;
  validityNotAfter?: string;
  isRevoked: boolean;
};
