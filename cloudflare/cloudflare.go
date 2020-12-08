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
	store        *Store
	tusClient    *tus.Client
}

// NewService 初始化服務
func NewService(cloudflareSetting CloudflareSetting) (*CloudflareService, error) {
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

	// 組合 API
	cloudflareStreamAPI := "https://" + c.apiDomain + "/client/" + c.apiVersion + "/accounts/" + c.accountID + "/stream"

	// 建立 header 認證資訊
	headers := make(http.Header)
	headers.Add("X-Auth-Email", c.email)
	headers.Add("X-Auth-Key", c.apiKey)

	// 初始化 tus config
	c.store = new(Store)
	config := &tus.Config{
		ChunkSize:           50 * 1024 * 1024, // Cloudflare Stream requires a minimum chunk size of 5MB.
		Resume:              true,
		OverridePatchMethod: false,
		Store:               c.store,
		Header:              headers,
		HttpClient:          nil,
	}

	// 初始化 client
	var err error
	c.tusClient, err = tus.NewClient(cloudflareStreamAPI, config)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Upload 上傳影音檔案
func (c *CloudflareService) Upload(uploadParameter UploadParameter) (*UploadReturnModel, error) {

	var err error

	// 把檔案打包到
	upload := tus.NewUpload(uploadParameter.Reader, uploadParameter.Size, uploadParameter.Metadata, uploadParameter.Fingerprint)

	// 建立上傳工作
	uploader, err := c.tusClient.CreateOrResumeUpload(upload)
	if err != nil {
		return nil, err
	}

	process := time.Now().Unix()
	go func() {
		for {
			fmt.Println(process, uploadParameter.Filename, upload.Progress(), "/100")
			time.Sleep(5 * time.Second)
			if upload.Finished() || upload.Progress() >= 99 {
				break
			}
			if uploader == nil {
				fmt.Println("uploader process", process, "上傳發生錯誤，停止進度")
				break
			}
		}
	}()

	// 開始上傳
	err = uploader.Upload()
	if err != nil {

		fmt.Println("[cf lib] 等待10秒")
		time.Sleep(10 * time.Second)

		// 確認真的上傳失敗，所以查詢看看
		videoSearchResponse, searchErr := c.Search(uploadParameter.Filename)
		if searchErr != nil {
			return nil, searchErr
		}
		if videoSearchResponse.Success && len(videoSearchResponse.Result) > 0 {
			fmt.Println("[cf lib] 上傳失敗，但有查到影片上傳成功了")
			var uploadReturnModel UploadReturnModel
			uploadReturnModel.Filename = uploadParameter.Filename
			uploadReturnModel.UID = videoSearchResponse.Result[0].UID
			return &uploadReturnModel, nil
		}

		fmt.Println("[cf lib] 上傳失敗，也查過沒有影片在 CF 上")
		uploader = nil
		return nil, err
	}

	// 上傳成功後，查詢結果
	videoSearchResponse, err := c.Search(uploadParameter.Filename)
	if err != nil {
		return nil, err
	}

	if videoSearchResponse.Success {
		var uploadReturnModel UploadReturnModel
		uploadReturnModel.Filename = uploadParameter.Filename
		uploadReturnModel.UID = videoSearchResponse.Result[0].UID
		return &uploadReturnModel, nil
	}

	return nil, errors.New("上傳結果異常")
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

// Search 查影片資訊
func (c *CloudflareService) AdvanceSearch(status, after string) ([]VideoSearchResponse, error) {
	endpoint := "https://" + c.apiDomain + "/client/" + c.apiVersion + "/accounts/" + c.accountID + "/stream?"

	if status != "" {
		endpoint = endpoint + "status=" + status
	}
	if after != "" {
		endpoint = endpoint + "&after=" + after
	}

	var videoSearchResponse []VideoSearchResponse
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

	fmt.Println(string(resp))
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
