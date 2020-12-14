package cloudflare

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go"
)

func TestCloudflareService_Upload(t *testing.T) {

	f, err := os.Open("BD-M02-OS-TS.ts")
	if err != nil {
		t.Error(err.Error())
	}

	defer f.Close()

	fi, _ := f.Stat()

	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = ""
	cloudflareSetting.Email = ""
	cloudflareSetting.AccountID = ""
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"

	meta := make(map[string]string)
	meta["name"] = fi.Name()
	meta["requireSignedURLs"] = "true"

	var uploadParameter UploadParameter
	uploadParameter.Filename = fi.Name()
	uploadParameter.Fingerprint = fmt.Sprintf("%s-%d-%s", fi.Name(), fi.Size(), fi.ModTime())
	uploadParameter.Metadata = meta
	uploadParameter.Reader = f
	uploadParameter.Size = fi.Size()

	cloudflareService, err := NewService(cloudflareSetting)
	if err != nil {
		t.Error(err.Error())
	}

	var uploadReturnModel *UploadReturnModel
	_ = retry.Do(
		func() error {
			uploadReturnModel, err = cloudflareService.Upload(uploadParameter)
			if err != nil {
				fmt.Println(err.Error())
			}
			return err
		},
	)

	if uploadReturnModel != nil {
		fmt.Println("上傳成功：", uploadReturnModel.Filename, uploadReturnModel.UID)
	} else {
		fmt.Println(fi.Name(), "上傳失敗")
	}

	t.Log("OK")
}

func TestCloudflareService_Search(t *testing.T) {
	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = ""
	cloudflareSetting.Email = ""
	cloudflareSetting.AccountID = ""
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"

	cloudflareService, err := NewService(cloudflareSetting)
	_, err = cloudflareService.Search("SKYHD-042-OS-TS.ts")
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("OK")
}

func TestCloudflareService_AdvanceSearch(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Taipei")
	a := time.Date(2020, 12, 8, 15, 18, 0, 0, loc)
	after := a.UTC().Format("2006-01-02T15:04:05Z07:00")
	fmt.Println(after)

	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = ""
	cloudflareSetting.Email = ""
	cloudflareSetting.AccountID = ""
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"

	cloudflareService, err := NewService(cloudflareSetting)
	_, err = cloudflareService.AdvanceSearch("ready", after)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("OK")
}

func TestCloudflareService_Delete(t *testing.T) {

	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = ""
	cloudflareSetting.Email = ""
	cloudflareSetting.AccountID = ""
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"

	cloudflareService, err := NewService(cloudflareSetting)
	err = cloudflareService.Delete("336ff6d4156685322ec89dbe92ff333f")
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("OK")
}

func TestCloudflareService_GetSignedURL(t *testing.T) {

	videoUID := "6efb8dd5f97c16a8e935e46316d0c673"

	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = ""
	cloudflareSetting.Email = ""
	cloudflareSetting.AccountID = ""
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"
	cloudflareSetting.Pem = "="
	cloudflareSetting.UtilDomain = "util.cloudflarestream.com"
	cloudflareSetting.KeyID = ""
	cloudflareSetting.StreamDomain = "watch.cloudflarestream.com"

	cloudflareService, err := NewService(cloudflareSetting)
	signedURL, err := cloudflareService.GetSignedURL(videoUID)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(signedURL)

	t.Log("OK")
}
