package ftx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	bsk "github.com/yin75620/crypto-berserker/setting"
)

type Ftx struct {
}

var (
	apiUrl = "https://ftx.com/api"
)

func (ftx *Ftx) GetAccountInfo() map[string]string {
	map[string][]string
	return ftx.get("account", &url.Values)
}

func (ftx *Ftx) buildRequest(reqMethod, path string, postForm *url.Values) error {
	httpRequest
}

func (ftx *Ftx) get(path string, params *url.Values) {
	return doMyRequest("GET", path, params)
}

func (ftx *Ftx) doMyRequest(method string, path string, params *url.Values) error {

	err := buildPostForm(method, path, params)
	if err != nil {
		return err
	}
	url := apiUrl + path

	//log.Println(sign, timestamp)
	resp, err := httpRequest(http.DefaultClient, method, url, reqBody, params)

	if err != nil {
		//log.Println(err)
		return err
	} else {
		//	log.Println(string(resp))
		return json.Unmarshal(resp, &response)
	}
}

func (ftx *Ftx) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("FTX-KEY", bsk.FTX_KEY)
	nanos := time.Now().UnixNano()
	ts := string(nanos / 1000000)
	postForm.Set("FTX-TS", ts)
	//domain := strings.Replace(ftx.config.Endpoint, "https://", "", len(ftx.config.Endpoint))
	boyd := "" //之後再用
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, path, body)
	sign, _ := GetParamHmacSHA256Base64Sign(bsk.FTX_API_SECRET_KEY, payload)
	postForm.Set("FTX-SIGN", sign)

	return nil
}

func httpRequest(client *http.Client, reqtype, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}
