package request

type CalculateHours struct {
	PageId        string `json:"page_id"`
	OrderId       string `json:"order_id"`
	NotionVersion string `json:"notion_version"`
}
