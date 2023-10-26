package wt

import (
	"encoding/json"
	"fmt"
	"github.com/peterbourgon/diskv/v3"
)

func NewCorpAuthInfoDao(store *diskv.Diskv) *CorpAuthInfoDao {
	return &CorpAuthInfoDao{store: store}
}

type CorpAuthInfoDao struct {
	store *diskv.Diskv
}

func (s *CorpAuthInfoDao) Write(info *AuthInfo) error {
	key := corpAuthInfoKey(info.AuthCorpInfo.CorpId)

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if err = s.Remove(info.AuthCorpInfo.CorpId); err != nil {
		return err
	}

	if err = s.store.Write(key, jsonBytes); err != nil {
		return fmt.Errorf("failed to write auth info: %w", err)
	}

	return nil
}

func (s *CorpAuthInfoDao) Get(corpId string) (*AuthInfo, error) {
	jsonBytes, err := s.store.Read(corpAuthInfoKey(corpId))
	if err != nil {
		return nil, fmt.Errorf("failed to read auth info: %w", err)
	}

	var info AuthInfo
	if err = json.Unmarshal(jsonBytes, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth info: %w", err)
	}

	return &info, nil
}

func (s *CorpAuthInfoDao) Remove(corpId string) error {
	key := corpAuthInfoKey(corpId)

	if s.store.Has(key) {
		if err := s.store.Erase(key); err != nil {
			return fmt.Errorf("failed to erase auth info: %w", err)
		}
	}

	return nil
}

type AuthInfo struct {
	DealerCorpInfo struct {
		CorpId   string `json:"corpid"`
		CorpName string `json:"corp_name"`
	} `json:"dealer_corp_info"`
	AuthCorpInfo struct {
		CorpId            string `json:"corpid"`
		CorpName          string `json:"corp_name"`
		CorpType          string `json:"corp_type"`
		CorpSquareLogoUrl string `json:"corp_square_logo_url"`
		CorpUserMax       int    `json:"corp_user_max"`
		CorpFullName      string `json:"corp_full_name"`
		VerifiedEndTime   int64  `json:"verified_end_time"`
		SubjectType       int    `json:"subject_type"`
		CorpWxQrCode      string `json:"corp_wxqrcode"`
		CorpScale         string `json:"corp_scale"`
		CorpIndustry      string `json:"corp_industry"`
		CorpSubIndustry   string `json:"corp_sub_industry"`
	} `json:"auth_corp_info"`
	AuthInfo struct {
		Agent []struct {
			AgentId          int    `json:"agentid"`
			Name             string `json:"name"`
			RoundLogoUrl     string `json:"round_logo_url"`
			SquareLogoUrl    string `json:"square_logo_url"`
			AuthMode         int    `json:"auth_mode"`
			IsCustomizedApp  bool   `json:"is_customized_app"`
			AuthFromThirdApp bool   `json:"auth_from_thirdapp"`
			Privilege        struct {
				Level      int      `json:"level"`
				AllowParty []int    `json:"allow_party"`
				AllowUser  []string `json:"allow_user"`
				AllowTag   []int    `json:"allow_tag"`
			} `json:"privilege"`
			SharedFrom struct {
				CorpId    string `json:"corpid"`
				ShareType int    `json:"share_type"`
			} `json:"shared_from"`
		} `json:"agent"`
	} `json:"auth_info"`
	AuthUserInfo struct {
		UserId     string `json:"userid"`
		OpenUserId string `json:"open_userid"`
		Name       string `json:"name"`
		Avatar     string `json:"avatar"`
	} `json:"auth_user_info"`
}

func corpAuthInfoKey(corpId string) string {
	return fmt.Sprintf("%s/auth_info.json", corpId)
}
