package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	tus "github.com/eventials/go-tus"
)

// CloudflareService 串流服務提供
type CloudflareService struct {
	apiKey       string
	email        string
	accountID    string
	apiDomain    string
	utilDomain   string
	apiVersion   string
	keyID        string
	pem          string
	streamDomain string
}

// NewService 初始化服務
func NewService(cloudflareSetting CloudflareSetting) *CloudflareService {
	c := new(CloudflareService)
	c.apiKey = cloudflareSetting.APIKey
	c.email = cloudflareSetting.Email
	c.accountID = cloudflareSetting.AccountID
	c.apiDomain = cloudflareSetting.APIDomain
	c.apiVersion = cloudflareSetting.APIVersion
	c.utilDomain = cloudflareSetting.UtilDomain
	c.keyID = cloudflareSetting.KeyID
	c.pem = cloudflareSetting.Pem
	c.streamDomain = cloudflareSetting.StreamDomain
	return c
}

// Upload 上傳影音檔案
func (c *CloudflareService) Upload(uploadParameter UploadParameter) (UploadReturnModel, error) {

	var uploadReturnModel UploadReturnModel

	headers := make(http.Header)
	headers.Add("X-Auth-Email", c.email)
	headers.Add("X-Auth-Key", c.apiKey)

	config := &tus.Config{
		ChunkSize:           5 * 1024 * 1024, // Cloudflare Stream requires a minimum chunk size of 5MB.
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              headers,
		HttpClient:          nil,
	}

	// 初始化 client
	client, err := tus.NewClient("https://"+c.apiDomain+"/client/"+c.apiVersion+"/accounts/"+c.accountID+"/stream", config)
	if err != nil {
		return uploadReturnModel, err
	}

	// 把檔案打包到
	upload := tus.NewUpload(uploadParameter.Reader, uploadParameter.Size, uploadParameter.Metadata, uploadParameter.Fingerprint)
	if err != nil {
		return uploadReturnModel, err
	}

	// upload.Metadata = meta

	// 建立上傳工作
	uploader, err := client.CreateUpload(upload)
	if err != nil {
		return uploadReturnModel, err
	}

	// 開始上傳
	fmt.Println("[cloudflare lib log] 開始上傳")
	err = uploader.Upload()
	if err != nil {
		return uploadReturnModel, err
	}
	fmt.Println("[cloudflare lib log] 結束上傳")

	// for {
	// 	progress := upload.Progress()
	// 	if progress == 100 {
	// 		break
	// 	}
	// 	fmt.Println(progress)
	// 	time.Sleep(1 * time.Second)
	// }

	// 上傳成功後，查詢結果
	videoSearchResponse, err := c.Search(uploadParameter.Filename)
	if err != nil {
		return uploadReturnModel, err
	}

	if videoSearchResponse.Success {
		uploadReturnModel.Filename = uploadParameter.Filename
		uploadReturnModel.UID = videoSearchResponse.Result[0].UID
		return uploadReturnModel, nil
	}

	return uploadReturnModel, errors.New("上傳結果異常")
}

// Search 查影片資訊
func (c *CloudflareService) Search(videoName string) (VideoSearchResponse, error) {
	endpoint := "https://" + c.apiDomain + "/client/" + c.apiVersion + "/accounts/" + c.accountID + "/stream?search=" + videoName

	var videoSearchResponse VideoSearchResponse
	var httpSetting HttpDoSetting
	httpSetting.AuthEmail = c.email
	httpSetting.AuthKey = c.apiKey
	httpSetting.Method = http.MethodGet
	httpSetting.Endpoint = endpoint
	httpSetting.Body = nil
	httpSetting.ContentType = "application/json"

	resp, err := c.httpDo(httpSetting)
	if err != nil {
		return videoSearchResponse, err
	}

	// fmt.Println(string(resp))
	json.Unmarshal(resp, &videoSearchResponse)
	// fmt.Println(videoSearchResponse)
	return videoSearchResponse, nil

}

// GetSignedURL 取得簽名過的影片網址
func (c *CloudflareService) GetSignedURL(videoUID string) (string, error) {

	signedURL := ""
	endpoint := "https://" + c.utilDomain + "/sign/" + videoUID

	// 組合 body 參數
	var signedURLBody SignedURLBody
	signedURLBody.ID = c.keyID
	signedURLBody.Pem = c.pem
	signedURLBody.Exp = time.Now().Add(2 * time.Hour).Unix()

	data, err := json.Marshal(signedURLBody)
	if err != nil {
		return signedURL, err
	}

	// 組合 http setting
	var httpSetting HttpDoSetting
	httpSetting.AuthEmail = c.email
	httpSetting.AuthKey = c.apiKey
	httpSetting.Method = http.MethodGet
	httpSetting.Endpoint = endpoint
	httpSetting.Body = data
	httpSetting.ContentType = "application/json"

	// 開始請求
	resp, err := c.httpDo(httpSetting)
	if err != nil {
		return signedURL, err
	}

	signedURL = "https://" + c.streamDomain + "/" + string(resp)

	return signedURL, nil
}

// httpDo 共用的 http 請求
func (c *CloudflareService) httpDo(httpDoSetting HttpDoSetting) (responseByte []byte, err error) {

	var request *http.Request
	var response *http.Response
	client := &http.Client{}

	if httpDoSetting.Body == nil {
		request, err = http.NewRequest(httpDoSetting.Method, httpDoSetting.Endpoint, nil)
		if err != nil {
			return
		}
	} else {
		request, err = http.NewRequest(httpDoSetting.Method, httpDoSetting.Endpoint, bytes.NewBuffer(httpDoSetting.Body))
		if err != nil {
			return
		}
	}

	if httpDoSetting.ContentType != "" {
		request.Header.Set("Content-Type", httpDoSetting.ContentType)
	}

	if httpDoSetting.AuthEmail != "" {
		request.Header.Set("X-Auth-Email", httpDoSetting.AuthEmail)
	}

	if httpDoSetting.AuthKey != "" {
		request.Header.Set("X-Auth-Key", httpDoSetting.AuthKey)
	}

	// 開始請求
	response, err = client.Do(request)
	if err != nil {
		return
	}

	// 讀取回傳資料
	responseByte, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// 關閉回傳的資料串流
	err = response.Body.Close()
	if err != nil {
		return
	}

	return
}
