// Copyright © 2022 zc2638 <zc2638@qq.com>.

package core

import (
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"
)

const tempToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQ1MDI0NDUsIm5iZiI6MTY2MzQwMDI0NSwiaWF0IjoxNjYzMzk4NDQ1LCJqdGkiOiJDTTpjYXRfbWF0Y2g6bHQxMjM0NTYiLCJvcGVuX2lkIjoiIiwidWlkIjo0NTk0MjYwMiwiZGVidWciOiIiLCJsYW5nIjoiIn0.1lXIcb1WL_SdsXG5N_i1drjjACRhRZUS2uadHlT6zIY"

func SendByUID(uid string) error {
	openID, err := getOpenID(uid)
	if err != nil {
		return fmt.Errorf("get OpenID failed: %v", err)
	}
	token, err := getTokenByOpenID(openID)
	if err != nil {
		return fmt.Errorf("get token failed: %v", err)
	}
	return Send(token, uid)
}

func getOpenID(uid string) (string, error) {
	resp, err := client.R().
		SetHeader("t", tempToken).
		SetQueryParam("uid", uid).
		Get("/sheep/v1/game/user_info")

	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("[%d] 请求错误: %s", resp.StatusCode(), resp.String())
	}

	jsonResult := gjson.Parse(resp.String())
	openid := jsonResult.Get("data.wx_open_id")
	if jsonResult.Get("err_code").Int() != 0 || openid.String() == "" {
		return "", fmt.Errorf("请求错误: %s", resp.String())
	}
	return openid.String(), nil
}

type User struct {
	UID      string `json:"uid"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nick_name"`
	Sex      int    `json:"sex"`
}

func getTokenByOpenID(openID string) (string, error) {
	resp, err := client.R().SetBody(&User{
		UID:      openID,
		Avatar:   "1",
		NickName: "1",
		Sex:      1,
	}).Post("/sheep/v1/user/login_oppo")

	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("[%d] 请求错误: %s", resp.StatusCode(), resp.String())
	}

	jsonResult := gjson.Parse(resp.String())
	token := jsonResult.Get("data.token")
	if jsonResult.Get("err_code").Int() != 0 || token.String() == "" {
		return "", fmt.Errorf("请求错误: %s", resp.String())
	}
	return token.String(), nil
}
