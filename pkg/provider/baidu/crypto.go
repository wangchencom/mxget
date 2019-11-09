package baidu

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/winterssy/mxget/pkg/cryptography"
	"github.com/winterssy/sreq"
)

func signPayload(params sreq.Params) sreq.Params {
	q := params.Encode()
	ts := fmt.Sprintf("%d", time.Now().Unix())
	r := fmt.Sprintf("baidu_taihe_music_secret_key%s", ts)
	key := fmt.Sprintf("%x", md5.Sum([]byte(r)))[8:24]
	param := base64.StdEncoding.EncodeToString(cryptography.AESCBCEncrypt([]byte(q), []byte(key), []byte(key)))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("baidu_taihe_music%s%s", param, ts))))
	return sreq.Params{
		"timestamp": ts,
		"param":     param,
		"sign":      sign,
	}
}
