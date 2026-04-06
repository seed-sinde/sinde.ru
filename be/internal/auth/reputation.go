package auth

import "context"

type IPReputation interface {
	Score(ctx context.Context, ip string) (int, string, error)
}
type NoopReputation struct{}

func (NoopReputation) Score(context.Context, string) (int, string, error) {
	return 0, "", nil
}
