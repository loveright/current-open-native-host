package service

import (
	"fmt"
	"notion-native/config"
	"notion-native/model/entity"
	"notion-native/utils"
	"os"
	"strconv"
)

func CalculateHours(pageId, orderId, notionVersion string) (float64, error) {
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
		fmt.Fprintln(os.Stderr, "error:", err)
		return ""
	}
	fmt.Fprintln(os.Stderr, "data_source_id:", page.Parent.DataSourceID)
	if page.Parent.DataSourceID != "" {
		return page.Parent.DataSourceID
	}
	return ""
}

func getDataSources(dataSourceId, orderId, notionVersion string) (float64, error) {

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

	total := 0.0

	for _, page := range dataSources.Results {
		if prop, ok := page.Properties["工时"]; ok {
			for _, item := range prop.MultiSelect {

				val, err := strconv.ParseFloat(item.Name, 64)
				if err != nil {
					continue
				}

				total += val
			}
		}
	}

	fmt.Fprintln(os.Stderr, "总工时:", total)
	return total, nil
}
