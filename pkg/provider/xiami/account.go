package xiami

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/sreq"
)

func LoginRaw(account string, password string) (*LoginResponse, error) {
	return std.LoginRaw(account, password)
}

// 登录接口，account 可为邮箱/手机号码
func (a *API) LoginRaw(account string, password string) (*LoginResponse, error) {
	token, err := a.getToken(APILogin)
	if err != nil {
		return nil, err
	}

	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	model := map[string]string{
		"account":  account,
		"password": password,
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(LoginResponse)
	err = a.Request(sreq.MethodGet, APILogin, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return resp, nil
}
