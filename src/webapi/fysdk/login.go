package fysdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FYSDKLoginResp struct {
	Code        int    `json:"code"`
	Description string `json:"desc"`
}

func Login(platform string, token string) (string, error) {
	form := url.Values{}
	form.Add("token", token)
	form.Add("time", fmt.Sprintf("%d", time.Now().Unix()))

	sign := Sign(form, GetLoginSecret())
	form.Add("sign", sign)

	urlPath := fmt.Sprintf("https://sdk2-syapi.737.com/sdk/index/%s/%s/user_check", _appId, platform)
	resp, err := http.Get(fmt.Sprintf("%s?%s", urlPath, form.Encode()))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("请求返回码异常: %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(string(body), "{") {
		var payload FYSDKLoginResp
		err := json.Unmarshal(body, &payload)
		if err != nil {
			return "", err
		}

		return "", errors.New(fmt.Sprintf("(%d)%s", payload.Code, payload.Description))
	}

	return string(body), nil
}
