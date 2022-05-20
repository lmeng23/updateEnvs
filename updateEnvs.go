package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type Option struct {
	Value   string `json:"value"`
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
	Id      int    `json:"id"`
}

func main() {
	const sourceUrl = '';
	const sourceClientID = '';
	const sourceClientSecret = ''

	const distUrl = '';
	const distClientID = ''
	const distClientSecret = ''

	sourceToken := getToken(sourceUrl, sourceClientID, sourceClientSecret)
	var cookies []string

	cookieSlice := getEnvs(sourceUrl, sourceToken)

	for _, v := range cookieSlice {
		data := v.(map[string]interface{})
		cookies = append(cookies, data["value"].(string))
	}

	cookieStr := strings.Join(cookies, "&")

	distToken := getToken(distUrl, distClientID, distClientSecret)
	cookiesInfo := getEnvs(distUrl, distToken)

	var id int
	for _, v := range cookiesInfo {
		data := v.(map[string]interface{})
		if data["name"] == "JD_COOKIE" {
			id = int(data["id"].(float64))
		}
	}
	res := updateEnvs(id, distUrl, distToken, cookieStr)
	if int(res["code"].(float64)) == 200 {
		fmt.Println("更新JD_COOKIE环境变量成功!")
	} else {
		fmt.Println("更新JD_COOKIE环境变量失败!")
	}
}

func errorInfo(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func conJson(body []byte) map[string]interface{} {
	data := make(map[string]interface{})
	err := json.Unmarshal(body, &data)
	errorInfo(err)
	return data
}

func getToken(url string, client_id, client_secret string) string {
	url = fmt.Sprintf("%v/open/auth/token?client_id=%v&client_secret=%v", url, client_id, client_secret)
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, payload)
	errorInfo(err)
	res, err := client.Do(req)
	errorInfo(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		errorInfo(err)
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	errorInfo(err)
	data := conJson(body)
	data = data["data"].(map[string]interface{})
	return fmt.Sprintf("%v %v", data["token_type"], data["token"])
}

func getEnvs(url, token string) []interface{} {
	url = fmt.Sprintf("%v/open/envs", url)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	errorInfo(err)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, payload)
	errorInfo(err)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	errorInfo(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		errorInfo(err)
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	errorInfo(err)
	data := conJson(body)
	return data["data"].([]interface{})
}

func updateEnvs(id int, url, token, value string) map[string]interface{} {
	url = fmt.Sprintf("%v/open/envs", url)

	fmt.Println(token)

	marshal, err := json.Marshal(Option{
		Value:   value,
		Name:    "JD_COOKIE",
		Remarks: "京东Cookie",
		Id:      id,
	})
	errorInfo(err)
	option := string(marshal)

	/*option := fmt.Sprintf(`
	{
		"value":"%v",
		"name":"JD_COOKIE",
		"remarks":"京东Cookie",
		"id":%v
	}`, value, id)*/

	payload := strings.NewReader(option)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, payload)
	errorInfo(err)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	errorInfo(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		errorInfo(err)
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	errorInfo(err)
	data := conJson(body)
	return data
}
