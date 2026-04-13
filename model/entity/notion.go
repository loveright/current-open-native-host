package entity

type NotionPage struct {
	ID         string              `json:"id"`
	Parent     Parent              `json:"parent"`
	Properties map[string]Property `json:"properties"`
}

type Parent struct {
	Type         string `json:"type"`
	DataSourceID string `json:"data_source_id"`
	DatabaseID   string `json:"database_id"`
}

type DataSources struct {
	Results []NotionPage `json:"results"`
}

type Property struct {
	MultiSelect []MultiSelect `json:"multi_select"`
}

type MultiSelect struct {
	Name string `json:"name"`
}
