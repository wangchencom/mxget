package netease

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/sreq"
)

func EmailLoginRaw(email string, password string) (*LoginResponse, error) {
	return std.EmailLoginRaw(email, password)
}

// 邮箱登录
func (a *API) EmailLoginRaw(email string, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"username":      email,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	err := a.Request(sreq.MethodPost, APIEmailLogin,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("email login: %s", resp.Msg)
	}

	return resp, nil
}

func CellphoneLoginRaw(countryCode int, phone int, password string) (*LoginResponse, error) {
	return std.CellphoneLoginRaw(countryCode, phone, password)
}

// 手机登录
func (a *API) CellphoneLoginRaw(countryCode int, phone int, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"phone":         phone,
		"countrycode":   countryCode,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	err := a.Request(sreq.MethodPost, APICellphoneLogin,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("cellphone login: %s", resp.Msg)
	}

	return resp, nil
}

func RefreshLoginRaw() (*CommonResponse, error) {
	return std.RefreshLoginRaw()
}

// 刷新登录状态
func (a *API) RefreshLoginRaw() (*CommonResponse, error) {
	resp := new(CommonResponse)
	err := a.Request(sreq.MethodPost, APIRefreshLogin,
		sreq.WithForm(weapi(struct{}{})),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("refresh login: %s", resp.Msg)
	}

	return resp, nil
}

func LogoutRaw() (*CommonResponse, error) {
	return std.LogoutRaw()
}

// 退出登录
func (a *API) LogoutRaw() (*CommonResponse, error) {
	resp := new(CommonResponse)
	err := a.Request(sreq.MethodPost, APILogout,
		sreq.WithForm(weapi(struct{}{})),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("logout: %s", resp.Msg)
	}

	return resp, nil
}
