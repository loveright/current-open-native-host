package service

import (
	"fmt"
	"notion-native/config"
	"notion-native/model/entity"
	"notion-native/utils"
	"strconv"
)

func CalculateHours(pageId, orderId, notionVersion string) (int, error) {
	dataSourceId := getDataSourceId(pageId, notionVersion)

	if dataSourceId == "" {
		return 0, fmt.Errorf("dataSourceId is empty")
	}

	total, err := getDataSources(dataSourceId, orderId, notionVersion)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func getDataSourceId(pageId, nodeVersion string) string {

	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageId)
	token := config.GetNotionToken()
	page, err := utils.HttpRequest[entity.NotionPage]("GET", url, token, nodeVersion, nil)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	fmt.Println("data_source_id:", page.Parent.DataSourceID)
	if page.Parent.DataSourceID != "" {
		return page.Parent.DataSourceID
	}
	return ""
}

func getDataSources(dataSourceId, orderId, notionVersion string) (int, error) {

	url := fmt.Sprintf("https://api.notion.com/v1/data_sources/%s/query", dataSourceId)
	token := config.GetNotionToken()

	jsonBody := []byte(fmt.Sprintf(`{
		"filter": {
			"property": "工单",
			"rich_text": {
				"equals": "%s"
			}
		}
	}`, orderId))

	dataSources, err := utils.HttpRequest[entity.DataSources](
		"POST",
		url,
		token,
		notionVersion,
		jsonBody,
	)
	if err != nil {
		return 0, err
	}

	total := 0

	for _, page := range dataSources.Results {
		if prop, ok := page.Properties["工时"]; ok {
			for _, item := range prop.MultiSelect {

				val, err := strconv.Atoi(item.Name)
				if err != nil {
					continue
				}

				total += val
			}
		}
	}

	fmt.Println("总工时:", total)
	return total, nil
}
