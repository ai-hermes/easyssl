package service

import "easyssl/server/internal/providercatalog"

func (s *Service) ListProviderDefinitions(kind string) []providercatalog.Definition {
	return providercatalog.ListByKind(kind)
}
