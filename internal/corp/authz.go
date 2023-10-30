package corp

import (
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"github.com/uptrace/bun"
)

// Authorization models the corp_authz table which stores the authorization information of the client corporation.
type Authorization struct {
	bun.BaseModel `bun:"table:corp_authz"`

	ID            int64                       `bun:"id,pk,autoincrement"`
	CorpID        string                      `bun:"corp_id,notnull"`
	CorpName      string                      `bun:"corp_name,notnull"`
	PermanentCode string                      `bun:"perm_code,notnull"`
	AuthInfo      *sqlitedb.JSON[AuthInfo]    `bun:"auth_info,notnull"`
	Permissions   *sqlitedb.JSON[Permissions] `bun:"perm,notnull"`
}

// AuthInfo contains the authorization event information from WeCom as a result of the client corporation authorizes
// access by the suite. This information is embedded in the payload of both create_auth_info event and change_auth_info
// event.
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

// Permissions contains the list of access permissions granted by the admin from the authorizing corporation.
// It is recommended for developers to monitor the changes of permissions and contact the administrator if any
// required permission was not granted.
type Permissions struct {
	AppPermissions []string `json:"app_permissions"`
}
