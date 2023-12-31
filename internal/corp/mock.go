package corp

import (
	"github.com/jarcoal/httpmock"
	"net/http"
)

func MockGetPermanentCodeEndpoint() {
	httpmock.RegisterResponder(
		http.MethodPost,
		"=~^"+getPermanentCodeUrl+".*",
		httpmock.NewStringResponder(http.StatusOK, `
{
	"errcode":0,
	"errmsg":"ok",
	"permanent_code": "xxxx", 
	"dealer_corp_info": 
	{
		"corpid": "xxxx",
		"corp_name": "name"
	},
	"auth_corp_info": 
	{
		"corpid": "xxxx",
		"corp_name": "name",
		"corp_type": "verified",
		"corp_square_logo_url": "yyyyy",
		"corp_user_max": 50,
		"corp_full_name":"full_name",
		"verified_end_time":1431775834,
		"subject_type": 1,
		"corp_wxqrcode": "zzzzz",
		"corp_scale": "1-50人",
		"corp_industry": "IT服务",
		"corp_sub_industry": "计算机软件/硬件/信息服务"
	},
	"auth_info":
	{
		"agent" :
		[
			{
				"agentid":1,
				"name":"NAME",
				"round_logo_url":"xxxxxx",
				"square_logo_url":"yyyyyy",
				"auth_mode":1,
				"is_customized_app":false,
				"auth_from_thirdapp":false,
				"privilege":
				{
					"level":1,
					"allow_party":[1,2,3],
					"allow_user":["zhansan","lisi"],
					"allow_tag":[1,2,3]
				},
				"shared_from":
				{
					"corpid":"wwyyyyy",
					"share_type": 1
				}
			}
		]
	},
	"auth_user_info":
	{
		"userid":"aa",
		"open_userid":"xxxxxx",
		"name":"xxx",
		"avatar":"http://xxx"
	},
	"state":"state001"
}
`),
	)
}

func MockGetAuthorizationInfoEndpoint() {
	httpmock.RegisterResponder(
		http.MethodPost,
		"=~^"+getAuthInfoUrl+".*",
		httpmock.NewStringResponder(http.StatusOK, `
{
    "errcode":0,
    "errmsg":"ok",
	"dealer_corp_info": 
	{
		"corpid": "xxxx",
		"corp_name": "name"
	},
	"auth_corp_info": 
	{
		"corpid": "xxxx",
		"corp_name": "name",
		"corp_type": "verified",
		"corp_square_logo_url": "yyyyy",
		"corp_user_max": 50,
		"corp_full_name":"full_name",
		"verified_end_time":1431775834,
		"subject_type": 1,
		"corp_wxqrcode": "zzzzz",
		"corp_scale": "1-50人",
		"corp_industry": "IT服务",
		"corp_sub_industry": "计算机软件/硬件/信息服务"
	},
	"auth_info":
	{
		"agent" :
		[
			{
				"agentid":1,
				"name":"NAME",
				"round_logo_url":"xxxxxx",
				"square_logo_url":"yyyyyy",
				"appid":1,
				"auth_mode":1,
				"is_customized_app":false,
				"auth_from_thirdapp":false,
				"privilege":
				{
					"level":1,
					"allow_party":[1,2,3],
					"allow_user":["zhansan","lisi"],
					"allow_tag":[1,2,3],
					"extra_party":[4,5,6],
					"extra_user":["wangwu"],
					"extra_tag":[4,5,6]
				},
				"shared_from":
				{
					"corpid":"wwyyyyy",
					"share_type": 1
				}
			},
			{
				"agentid":2,
				"name":"NAME2",
				"round_logo_url":"xxxxxx",
				"square_logo_url":"yyyyyy",
				"appid":5,
				"shared_from":
				{
					"corpid":"wwyyyyy",
					"share_type": 0
				}
			}
		]
	}
}
`),
	)
}

func MockGetPermissionsEndpoint() {
	httpmock.RegisterResponder(
		http.MethodPost,
		"=~^"+getPermissionsUrl+".*",
		httpmock.NewStringResponder(http.StatusOK, `
{
	"errcode":0,
	"errmsg":"ok",
	"app_permissions":["contact:sensitive:user_name"]
}
`),
	)
}

func MockGetCorpAccessTokenEndpoint() {
	httpmock.RegisterResponder(
		http.MethodGet,
		"=~^"+getCorpAccessTokenUrl+".*",
		httpmock.NewStringResponder(http.StatusOK, `
{
   "errcode": 0,
   "errmsg": "ok",
   "access_token": "accesstoken000001",
   "expires_in": 7200
}
`),
	)
}
