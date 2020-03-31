package kan_sdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	sign "github.com/kan-fun/kan-core"
)

type Client struct {
	credential *sign.Credential
}

func NewClient(accessKey, secretKey string) (client *Client, err error) {
	credential, err := sign.NewCredential(accessKey, secretKey)
	if err != nil {
		return
	}

	client = &Client{
		credential,
	}

	return
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (client *Client) consPostData(specificParameter map[string]string) (data url.Values, commonParameter *sign.CommonParameter, signature string) {
	commonParameter = &sign.CommonParameter{
		AccessKey:      client.credential.AccessKey,
		SignatureNonce: uuid.New().String(),
		Timestamp:      strconv.FormatInt(makeTimestamp(), 10),
	}

	signature = client.credential.Sign(*commonParameter, specificParameter)

	data = map[string][]string{}

	for k, v := range specificParameter {
		s := make([]string, 1)
		s[0] = v

		data[k] = s
	}

	return
}

func (client *Client) post(path string, specificParameter map[string]string) (err error) {
	data, commonParameter, signature := client.consPostData(specificParameter)
	body := strings.NewReader(data.Encode())

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.kan-fun.com/%s", path), body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("X-Ca-Key", commonParameter.AccessKey)
	req.Header.Set("X-Ca-Timestamp", commonParameter.Timestamp)
	req.Header.Set("X-Ca-Nonce", commonParameter.SignatureNonce)
	req.Header.Set("X-Ca-Signature", signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return errors.New(buf.String())
	}

	return
}

func (client *Client) Email(topic string, msg string) (err error) {
	specificParameter := map[string]string{
		"topic": topic,
		"msg":   msg,
	}

	return client.post("send-email", specificParameter)
}
