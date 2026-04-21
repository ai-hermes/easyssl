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
};

export type Certificate = {
  id: string;
  subjectAltNames: string;
  serialNumber: string;
  keyAlgorithm: string;
  validityNotAfter?: string;
  isRevoked: boolean;
};
