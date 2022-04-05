package secrets

import (
	"context"
	"errors"
)

func (s *secretsManager) getPublicKey(ctx context.Context) (string, string, error) {
	pk, _, err := s.c.Actions.GetRepoPublicKey(ctx, s.owner, s.repo)
	if err != nil {
		return "", "", err
	}

	if pk.GetKey() == "" {
		return "", "", errors.New("could not get public key")
	}

	return pk.GetKey(), pk.GetKeyID(), nil
}
