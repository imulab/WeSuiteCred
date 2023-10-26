package wt

import (
	"fmt"
	"github.com/peterbourgon/diskv/v3"
)

func NewCorpSecretDao(store *diskv.Diskv) *CorpSecretDao {
	return &CorpSecretDao{store: store}
}

type CorpSecretDao struct {
	store *diskv.Diskv
}

func (s *CorpSecretDao) Write(corpId string, corpSecret string) error {
	key := corpSecretKey(corpId)

	if err := s.Remove(corpId); err != nil {
		return err
	}

	if err := s.store.WriteString(key, corpSecret); err != nil {
		return fmt.Errorf("failed to write corp_secret: %w", err)
	}

	return nil
}

func (s *CorpSecretDao) Get(corpId string) (string, error) {
	v := s.store.ReadString(corpSecretKey(corpId))
	if len(v) == 0 {
		return "", fmt.Errorf("corp_secret not found")
	}

	return v, nil
}

func (s *CorpSecretDao) Remove(corpId string) error {
	key := corpSecretKey(corpId)

	if s.store.Has(key) {
		if err := s.store.Erase(key); err != nil {
			return fmt.Errorf("failed to erase corp_secret: %w", err)
		}
	}

	return nil
}

func corpSecretKey(corpId string) string {
	return fmt.Sprintf("%s/corp_secret", corpId)
}
