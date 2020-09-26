package cloudflare

// CloudflareSetting 設置與 cloudflare 溝通的基本參數
type CloudflareSetting struct {
	APIKey       string
	Email        string
	AccountID    string
	APIDomain    string
	UtilDomain   string
	APIVersion   string
	KeyID        string
	Pem          string
	StreamDomain string
}

// HttpDoSetting 請求參數設置
type HttpDoSetting struct {
	Method      string
	Endpoint    string
	Body        []byte
	ContentType string
	AuthEmail   string
	AuthKey     string
}

type UploadReturnModel struct {
	Filename string
	UID      string
}

type SignedURLBody struct {
	ID  string `json:"id"`
	Pem string `json:"pem"`
	Exp int64  `json:"exp"`
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

// VideoSearchResponse 開始區塊
type VideoSearchResponse struct {
	Result   []VideoSearchResultModel `json:"result"`
	Success  bool                     `json:"success"`
	Errors   []Error                  `json:"errors"`
	Messages interface{}              `json:"messages"`
}

type VideoSearchResultModel struct {
	AllowedOrigins        []string                  `json:"allowedOrigins"`
	Created               string                    `json:"created"`
	Duration              int                       `json:"duration"`
	Input                 VideoSearchInputModel     `json:"input"`
	MaxDurationSeconds    int                       `json:"maxDurationSeconds"`
	Meta                  map[string]string         `json:"meta"`
	Modified              string                    `json:"modified"`
	UploadExpiry          string                    `json:"uploadExpiry"`
	Playback              VideoSearchPlaybackModel  `json:"playback"`
	Preview               string                    `json:"preview"`
	ReadyToStream         bool                      `json:"readyToStream"`
	RequireSignedURLs     bool                      `json:"requireSignedURLs"`
	Size                  int                       `json:"size"`
	Status                VideoSearchStatusModel    `json:"status"`
	Thumbnail             string                    `json:"thumbnail"`
	ThumbnailTimestampPct float32                   `json:"thumbnailTimestampPct"`
	UID                   string                    `json:"uid"`
	Uploaded              string                    `json:"uploaded"`
	Watermark             VideoSearchWatermarkModel `json:"watermark"`
}

type VideoSearchInputModel struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VideoSearchPlaybackModel struct {
	Hls  int `json:"hls"`
	Dash int `json:"dash"`
}

type VideoSearchStatusModel struct {
	State       string `json:"state"`
	PctComplete int    `json:"pctComplete"`
}

type VideoSearchWatermarkModel struct {
	UID            string  `json:"uid"`
	Size           int     `json:"size"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	Created        string  `json:"created"`
	DownloadedFrom string  `json:"downloadedFrom"`
	Name           string  `json:"name"`
	Opacity        float32 `json:"opacity"`
	Padding        float32 `json:"padding"`
	Scale          float32 `json:"scale"`
	Position       string  `json:"position"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// VideoSearchResponse 結束區塊
