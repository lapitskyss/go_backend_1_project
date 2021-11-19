package linksrv

import "context"

func (s *LinkService) List(ctx context.Context, hashes []string) (<-chan Link, <-chan error) {
	return s.ls.GetByHashes(ctx, hashes)
}
