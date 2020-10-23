package caresdk

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	core "github.com/byte-care/care-core"
)

// Client ...
type Client struct {
	credential *core.Credential
}

// NewClient ...
func NewClient(accessKey, secretKey string) (client *Client, err error) {
	credential, err := core.NewCredential(accessKey, secretKey)
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

func (client *Client) consPostData(specificParameter map[string]string) (data url.Values, commonParameter *core.CommonParameter, signature string) {
	commonParameter = &core.CommonParameter{
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

func consAPIURL(path string) string {
	return fmt.Sprintf("https://api.bytecare.xyz/%s", path)
}

func (client *Client) post(path string, specificParameter map[string]string) (err error) {
	data, commonParameter, signature := client.consPostData(specificParameter)
	body := strings.NewReader(data.Encode())

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", consAPIURL(path), body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Care-Key", commonParameter.AccessKey)
	req.Header.Set("Care-Timestamp", commonParameter.Timestamp)
	req.Header.Set("Care-Nonce", commonParameter.SignatureNonce)
	req.Header.Set("Care-Signature", signature)

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

// Email ...
func (client *Client) Email(topic string, msg string) (err error) {
	specificParameter := map[string]string{
		"topic": topic,
		"msg":   msg,
	}

	return client.post("send-email", specificParameter)
}

// LogClient ...
type LogClient struct {
	*Client
	conn *websocket.Conn
}

func logFail() {
	log.Println("❌ Fail to ByteCare")
}

func readLoop(conn *websocket.Conn) {
	for {
		if messageType, _, err := conn.ReadMessage(); err != nil {
			log.Println(messageType)
			log.Println(err)

			break
		}
	}
}

// NewLogClient ...
func NewLogClient(accessKey, secretKey, topic string, isPro bool) (logClient *LogClient, err error) {
	client, err := NewClient(accessKey, secretKey)
	if err != nil {
		logFail()
		return nil, err
	}

	url_ := url.URL{Scheme: "wss", Host: "live.bytecare.xyz", Path: "/log/pub"}

	_, commonParameter, signature := client.consPostData(nil)

	header := http.Header{
		"Care-Key":       {commonParameter.AccessKey},
		"Care-Timestamp": {commonParameter.Timestamp},
		"Care-Nonce":     {commonParameter.SignatureNonce},
		"Care-Signature": {signature},
	}

	conn, resp, err := websocket.DefaultDialer.Dial(url_.String(), header)
	if err != nil {
		logFail()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, errors.New(bodyString)
	}

	go readLoop(conn)

	err = conn.WriteMessage(websocket.TextMessage, []byte(topic))
	if err != nil {
		logFail()
		return
	}

	var proString string
	if isPro {
		proString = "1"
	} else {
		proString = "0"
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(proString))
	if err != nil {
		logFail()
		return
	}

	logClient = &LogClient{
		client,
		conn,
	}

	log.Println("✅ ByteCare")

	return
}

// PubLog ...
func (logClient *LogClient) PubLog(content string) (err error) {
	err = logClient.conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		return
	}

	return
}

// CloseLog ...
func (logClient *LogClient) CloseLog(isSuccessful bool) (err error) {
	var closeCode int = websocket.CloseNormalClosure

	if !isSuccessful {
		closeCode = 4000
	}

	err = logClient.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, ""))
	if err != nil {
		return
	}

	return
}
