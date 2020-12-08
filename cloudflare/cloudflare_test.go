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
	cloudflareSetting.APIKey = "2271dd6dacc6d007c16dd2fbe56855f21fb67"
	cloudflareSetting.Email = "aa71435723@gmail.com"
	cloudflareSetting.AccountID = "55a3cef3d255aabc5a510c7ba1caead5"
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"

	cloudflareService, err := NewService(cloudflareSetting)
	_, err = cloudflareService.Search("BDSR-254-2-NS-TS.ts")
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

func TestCloudflareService_GetSignedURL(t *testing.T) {

	videoUID := "6efb8dd5f97c16a8e935e46316d0c673"

	var cloudflareSetting CloudflareSetting
	cloudflareSetting.APIKey = "2876bc25dd10fee6a8f1baec1012a2d31aaa1"
	cloudflareSetting.Email = "ziweiyuntltd@gmail.com"
	cloudflareSetting.AccountID = "4d1046c79adc4fbc33c93cdb2caacc3c"
	cloudflareSetting.APIDomain = "api.cloudflare.com"
	cloudflareSetting.APIVersion = "v4"
	cloudflareSetting.Pem = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBdmt2c21Fc21pODNmU244d1pkS215NEZJUFdMdlFiVU15UUx3Mm5ucGRwNW1uZ0RyCnB4bjZDY3FGMXZ1L0oyQkJjeVJvZUZPZkxsWUZBZVFSMXJYdUxHZ3hvR0Nya2xIZnBlLzhKbngyOEU3SFhRQlUKaUpuQllIa3RFWmk0d0NCbDZvcytoVnVtWXZuQjQxamV2VEEvTnVIVlRGUDhBWU5JSGZCVFJEY2JiNG5ZRmhaTQo5ZzlSUXRPRXBERmdja3NZeVBGU0ozbmxUeGZEUDRMRjdjUWFEUkdtcWI1R1lXOTBtalBqQkpzOTUzUlYxV0MxCnJCTE4wOEFHZDh4UFJ2aU1aOGNQam1QamZ1b2NqTExRTjRHSmdhR1p4SkRXQW0zQ3VDUHNGRW9hRmd2cGt6NzgKb0RaUitLdXVZVmxxVnpOQ3JoQlJLeVBzVGw1ZDdHZGpIQW1NQ3dJREFRQUJBb0lCQVFDYTl3Y25pZU5NN0F6UApCTDVyM053NVV3RjZBK3dralFScFdQeThYWlR5Sk5JYUQxUFgwejZiNUpHVFhaVHZ3dUhwbXhkOERWVE9qZndyCjZ3ZGYydTJtdWY3WHhJRlRlVnJ3TFhzZitERi9SaGZ4czBnanFWb2hidXgxclBHZWU0T2pPVnRqakJ2MTg3K2gKblFoZDlrRTBOem5VbTN2WDI0bkozNkJmSjZVdThTQUhhVTBvSWVVcTErKzF1RUIrc2x4WmRFQ2plYzJaWlo1MwpVeTBaT2dnT0ZaWEZVVElHRkpHZVNRUGw1M3R6dE1zdTl4YWU0b2VnMWo3a3ZGcUFldTRwYVR6cGxmcUVCZ1ZpCkZUMFUxSUVqanJ3VVFMekk5dzVncnZ5RXAvNkZKaFd0ejZ6OG94SzZ6TlhGNkQ2VkplVEZERlltK1FXYXlQKysKSWwxS0xjVEJBb0dCQU1ycGgzOXFHNVdBOWZIQk45dThENWpuejA1eXhjRnJVaGQ1L1dJY21xc09maWRuTXpwQgo1UGhGZGZyRyswWFpMUVRmc29QaFQ5VnJub3o2Rm1FRWtuQnBIZXlOQzBibW5aT21wb1dhUTUxVDF6M3d5V08vCk1vYWVmTkFMU1JTSlBhTlBPQ21nM0xYRzc2ajIxNVpWWnlaVmJrZ2EwaWdJUTJweHl6UzMxUzFyQW9HQkFQQVYKYjRvY2ZscS9peHJMdUZMOUVUMVVZZXYwRmh3SlV6KzcwdHNualN3QzhDdjg4NUtYSWJtazhpbWQyTUdveEljRwpob1d3SkxzMld1UzZiV0lFTCt2MVF3a1VWbEJFc3FsTHNHWUYxQTEvUVVjWWVGUHZGam5TcnJkY243OG0wRDh0Ck02cGtJQ1NPM3Yzbnp4R01TSUJrcUhLZ0VaOElVSjExR2xSOEpTUGhBb0dCQUxnMGpveHQ0RUsxd3hCSVR4Uk0Kd29BV0dRMW5oZjFVRno4MndIOGI1cEZaWTg4VGtkN1dTUzNWcVFnVE1iTTBOL2xQdG5pZ3gxL0JCanVISVYvTAp6Y0Q4dkd5dGtrbzRPMTc2RC93RGtsUTE4NVhJakpyZnpOZUc4MW5PbFBadXJLVWYycVYzNGtXbkpwUm1Ha3JnCmx2YW00YW5WcDJrdUx6MW50b2pTUmxXbkFvR0JBSyswVUdRNGhEU3YrQU1OVXdIUldidVR0UEoxT1hVZFVnTFQKMS9ZeDFQeC96ZnV6YlNNOFhoODZXMHdmekZHMnpOV3c3ZVNMUythRFdqUUpTQ0l5eEV1Z3ZJVzVqNDNCS1N3RApTNzd1eHdsMXQzVnJzQ3hsVHRQVW42OXNKekZERzZjUTByNEI5eEFxUzRKeEV6ZFpmbm9Rc01McTZOcUZ3RkhzCk1PL2h4MkNoQW9HQU16Q294b1RnSFk5dXFWajBvZ3UvWHZoWjFqbFJSRXdCMGNCMEMyMmExOHVUSXRHaWdhZEQKMGx0aHBVVjNjdnlNQUlMQTRSQWtaaVlGVWZ5b2ZRQ0gwSXBZak9mWmRpOGlTbnhkWDVlYTVSaENscGRubTJ3Vgo5Uzh0STF1T2VqSlZyQUh6RjZjSkMvWVJiT0RvNXBNaWtIZW82UzJxZHZlMTdXeGNNdDFySzlRPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	cloudflareSetting.UtilDomain = "util.cloudflarestream.com"
	cloudflareSetting.KeyID = "09adbe4d12fe7c40d0ecb790c3da6b41"
	cloudflareSetting.StreamDomain = "watch.cloudflarestream.com"

	cloudflareService, err := NewService(cloudflareSetting)
	signedURL, err := cloudflareService.GetSignedURL(videoUID)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(signedURL)

	t.Log("OK")
}
