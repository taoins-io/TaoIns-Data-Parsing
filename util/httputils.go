package util

import (
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"tao/logger"
	"time"
)

type Res struct {
	StatusCode int
	Content    []byte
}

func GetMethod(url string) (res Res, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		logger.GetLogger().Errorf("GetMethod err %v %v %v", zap.Error(err), zap.Any("url", url), zap.Any("response", r))
		return
	}
	defer r.Body.Close()
	response, _ := ioutil.ReadAll(r.Body)
	res.StatusCode = r.StatusCode
	res.Content = response
	return
}
