package speech

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsHost  = "iat-api.xfyun.cn"
	wsPath  = "/v2/iat"
	hostUrl = "wss://iat-api.xfyun.cn/v2/iat"
)

const (
	STATUS_FIRST_FRAME    = 0
	STATUS_CONTINUE_FRAME = 1
	STATUS_LAST_FRAME     = 2
)

// Config 语音识别配置
type Config struct {
	APIKey    string
	APISecret string
	AppID     string
}

// Recognizer 语音识别器
type Recognizer struct {
	config Config
}

// NewRecognizer 创建新的语音识别器
func NewRecognizer(config Config) *Recognizer {
	return &Recognizer{
		config: config,
	}
}

// RecognizeFile 识别音频文件
func (r *Recognizer) RecognizeFile(audioFile string) (*string, error) {
	fmt.Println(r.config)
	st := time.Now()
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl(hostUrl, r.config.APIKey, r.config.APISecret), nil)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %v, status=%s, body=%s", err, resp.Status, readResp(resp))
	} else if resp.StatusCode != 101 {
		return nil, fmt.Errorf("连接失败: %v, status=%s, body=%s", err, resp.Status, readResp(resp))
	}
	//打开音频文件

	var frameSize = 1280 //每一帧的音频大小
	// var intervel = 40 * time.Millisecond //发送音频间隔
	//开启协程，发送数据
	ctx, _ := context.WithCancel(context.Background())
	defer conn.Close()
	var status = 0
	go func() {
		//	start:
		audioFile, err := os.Open(audioFile)
		if err != nil {
			panic(err)
		}
		status = STATUS_FIRST_FRAME //音频的状态信息，标识音频是第一帧，还是中间帧、最后一帧
		//		time.Sleep(20*time.Second)
		var buffer = make([]byte, frameSize)
		for {
			len, err := audioFile.Read(buffer)
			if err != nil {
				if err == io.EOF { //文件读取完了，改变status = STATUS_LAST_FRAME
					status = STATUS_LAST_FRAME
				} else {
					panic(err)
				}
			}
			select {
			case <-ctx.Done():
				fmt.Println("session end ---")
				return
			default:
			}
			switch status {
			case STATUS_FIRST_FRAME: //发送第一帧音频，带business 参数
				frameData := map[string]interface{}{
					"common": map[string]interface{}{
						"app_id": r.config.AppID, //appid 必须带上，只需第一帧发送
					},
					"business": map[string]interface{}{ //business 参数，只需一帧发送
						"language": "zh_cn",
						"domain":   "iat",
						"accent":   "mandarin",
					},
					"data": map[string]interface{}{
						"status":   STATUS_FIRST_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				fmt.Println("send first")
				conn.WriteJSON(frameData)
				status = STATUS_CONTINUE_FRAME
			case STATUS_CONTINUE_FRAME:
				frameData := map[string]interface{}{
					"data": map[string]interface{}{
						"status":   STATUS_CONTINUE_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				conn.WriteJSON(frameData)
			case STATUS_LAST_FRAME:
				frameData := map[string]interface{}{
					"data": map[string]interface{}{
						"status":   STATUS_LAST_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				conn.WriteJSON(frameData)
				fmt.Println("send last")
				return
				//	goto start
			}

		}

	}()

	//获取返回的数据
	var decoder Decoder
	for {
		var resp = RespData{}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}
		json.Unmarshal(msg, &resp)
		//fmt.Println(string(msg))
		fmt.Println(resp.Data.Result.String(), resp.Sid)
		if resp.Code != 0 {
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			return nil, fmt.Errorf("识别失败: %v", resp.Message)
		}
		decoder.Decode(&resp.Data.Result)
		if resp.Data.Status == 2 {
			//cf()
			//fmt.Println("final:",decoder.String())
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			break
			//return
		}

	}

	res := decoder.String()
	return &res, nil
}

// // RecognizeAudio 直接识别音频数据
// func (r *Recognizer) RecognizeAudio(audioData []byte) (<-chan Result, error) {
// 	// 建立WebSocket连接（签名 URL 放在查询参数）
// 	wsURL := "wss://" + wsHost + wsPath
// 	signedURL := assembleAuthUrl(wsURL, r.config.APIKey, r.config.APISecret)
// 	conn, resp, err := websocket.DefaultDialer.Dial(signedURL, nil)
// 	if err != nil {
// 		if resp != nil {
// 			body, _ := io.ReadAll(resp.Body)
// 			return nil, fmt.Errorf("连接失败: %v, status=%s, body=%s", err, resp.Status, string(body))
// 		}
// 		return nil, fmt.Errorf("连接失败: %v", err)
// 	}

// 	resultChan := make(chan Result)

// 	go func() {
// 		defer conn.Close()
// 		defer close(resultChan)

// 		frameSize := 1280
// 		interval := 40 * time.Millisecond
// 		total := len(audioData)

// 		// 发送首帧（含 common/business）
// 		firstEnd := frameSize
// 		if firstEnd > total {
// 			firstEnd = total
// 		}
// 		first := map[string]any{
// 			"common": map[string]any{
// 				"app_id": r.config.AppID,
// 			},
// 			"business": map[string]any{
// 				"language": "zh_cn",
// 				"domain":   "iat",
// 				"accent":   "mandarin",
// 			},
// 			"data": map[string]any{
// 				"status":   0,
// 				"format":   "audio/L16;rate=16000",
// 				"audio":    base64.StdEncoding.EncodeToString(audioData[:firstEnd]),
// 				"encoding": "raw",
// 			},
// 		}
// 		_ = conn.WriteJSON(first)
// 		time.Sleep(interval)

// 		// 中间帧
// 		for off := firstEnd; off+frameSize < total; off += frameSize {
// 			mid := map[string]any{
// 				"data": map[string]any{
// 					"status":   1,
// 					"format":   "audio/L16;rate=16000",
// 					"audio":    base64.StdEncoding.EncodeToString(audioData[off : off+frameSize]),
// 					"encoding": "raw",
// 				},
// 			}
// 			_ = conn.WriteJSON(mid)
// 			time.Sleep(interval)
// 		}

// 		// 尾帧（带最后一段音频）
// 		lastStart := firstEnd
// 		if lastStart+frameSize < total {
// 			lastStart = total - frameSize
// 		}
// 		if lastStart < total {
// 			last := map[string]any{
// 				"data": map[string]any{
// 					"status":   2,
// 					"format":   "audio/L16;rate=16000",
// 					"audio":    base64.StdEncoding.EncodeToString(audioData[lastStart:]),
// 					"encoding": "raw",
// 				},
// 			}
// 			_ = conn.WriteJSON(last)
// 		}

// 		// 接收识别结果
// 		for {
// 			_, message, err := conn.ReadMessage()
// 			if err != nil {
// 				resultChan <- Result{Err: fmt.Errorf("接收失败: %v", err)}
// 				return
// 			}
// 			resultChan <- Result{Text: string(message)}
// 		}
// 	}()

//		return resultChan, nil
//	}
type RespData struct {
	Sid     string `json:"sid"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Result Result `json:"result"`
	Status int    `json:"status"`
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

// 解析返回数据，仅供demo参考，实际场景可能与此不同。
type Decoder struct {
	results []*Result
}

func (d *Decoder) Decode(result *Result) {
	if len(d.results) <= result.Sn {
		d.results = append(d.results, make([]*Result, result.Sn-len(d.results)+1)...)
	}
	if result.Pgs == "rpl" {
		for i := result.Rg[0]; i <= result.Rg[1]; i++ {
			d.results[i] = nil
		}
	}
	d.results[result.Sn] = result
}

func (d *Decoder) String() string {
	var r string
	for _, v := range d.results {
		if v == nil {
			continue
		}
		r += v.String()
	}
	return r
}

type Result struct {
	Ls  bool   `json:"ls"`
	Rg  []int  `json:"rg"`
	Sn  int    `json:"sn"`
	Pgs string `json:"pgs"`
	Ws  []Ws   `json:"ws"`
}

func (t *Result) String() string {
	var wss string
	for _, v := range t.Ws {
		wss += v.String()
	}
	return wss
}

type Ws struct {
	Bg int  `json:"bg"`
	Cw []Cw `json:"cw"`
}

func (w *Ws) String() string {
	var wss string
	for _, v := range w.Cw {
		wss += v.W
	}
	return wss
}

type Cw struct {
	Sc int    `json:"sc"`
	W  string `json:"w"`
}
