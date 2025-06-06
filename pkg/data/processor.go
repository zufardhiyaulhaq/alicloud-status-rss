package data

import (
	"slices"

	"github.com/zufardhiyaulhaq/alicloud-status-rss/pkg/model"
)

func ProcessRSS(rssData []model.RSS) ([]model.RSSItem, error) {
	guidMap := make(map[string]*model.RSSItem)

	for _, rss := range rssData {
		for _, item := range rss.Channel.Items {
			if existing, ok := guidMap[item.GUID]; ok {
				if !slices.Contains(existing.Type, rss.Type) {
					existing.Type = append(existing.Type, rss.Type)
				}
			} else {
				guidMap[item.GUID] = &model.RSSItem{
					Title:   item.Title,
					Link:    item.Link,
					Content: item.Content,
					Type:    []string{rss.Type},
					GUID:    item.GUID,
				}
			}
		}
	}

	result := make([]model.RSSItem, 0, len(guidMap))
	for _, v := range guidMap {
		result = append(result, *v)
	}

	return result, nil
}
