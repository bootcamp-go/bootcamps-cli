package secrets

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/google/go-github/v43/github"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

const (
	keySize   = 32
	nonceSize = 24
)

// Test indirection
var generateKey = box.GenerateKey

type SecretManager interface {
	SetSecret(context.Context, string, string) error
}

type secretsManager struct {
	c     *github.Client
	repo  string
	owner string
}

func NewSecretsManager(client *github.Client, owner, repo string) SecretManager {
	return &secretsManager{
		c:     client,
		repo:  repo,
		owner: owner,
	}
}

func (s *secretsManager) SetSecret(ctx context.Context, key, value string) error {
	ev, pkid, err := s.encryptValue(ctx, value)
	if err != nil {
		return err
	}

	if ev == "" {
		return nil
	}

	return s.sendSecretRequest(ctx, pkid, key, ev)
}

func (s *secretsManager) encryptValue(ctx context.Context, value string) (string, string, error) {
	if value == "" {
		return "", "", nil
	}

	pk, pkid, err := s.getPublicKey(ctx)
	if err != nil {
		return "", "", err
	}

	data, err := encryptWithLibSodium(pk, value)
	if err != nil {
		return "", "", err
	}

	return data, pkid, nil
}

func encryptWithLibSodium(pk, value string) (string, error) {
	// decode the provided public key from base64
	recipientKey := new([keySize]byte)
	b, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return "", err
	} else if size := len(b); size != keySize {
		return "", fmt.Errorf("recipient public key has invalid length (%d bytes)", size)
	}

	copy(recipientKey[:], b)

	// create an ephemeral key pair
	pubKey, privKey, err := generateKey(rand.Reader)
	if err != nil {
		return "", err
	}

	// create the nonce by hashing together the two public keys
	nonce := new([nonceSize]byte)
	nonceHash, err := blake2b.New(nonceSize, nil)
	if err != nil {
		return "", err
	}

	if _, err := nonceHash.Write(pubKey[:]); err != nil {
		return "", err
	}

	if _, err := nonceHash.Write(recipientKey[:]); err != nil {
		return "", err
	}

	copy(nonce[:], nonceHash.Sum(nil))

	// begin the output with the ephemeral public key and append the encrypted content
	out := box.Seal(pubKey[:], []byte(value), nonce, recipientKey, privKey)

	// base64-encode the final output
	return base64.StdEncoding.EncodeToString(out), nil
}

func (s *secretsManager) sendSecretRequest(ctx context.Context, pkid, name, value string) error {
	es := &github.EncryptedSecret{
		Name:           name,
		EncryptedValue: strings.Trim(value, "\n"),
		KeyID:          pkid,
	}

	color.Print("cyan", fmt.Sprintf("Sending %s secret request to GitHub...", name))

	res, err := s.c.Actions.CreateOrUpdateRepoSecret(ctx, s.owner, s.repo, es)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		if res.StatusCode != 204 {
			return errors.New("could not create or update secret")
		}
	}

	return nil
}
