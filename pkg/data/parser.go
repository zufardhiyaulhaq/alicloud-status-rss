package data

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
)

func ParseRSS(ctx context.Context, url string) (*model.RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rss model.RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}
	return &rss, nil
}
