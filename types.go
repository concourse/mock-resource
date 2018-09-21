package resource

type Source struct {
	FetchImage bool `json:"fetch_image"`
}

type Version struct {
	Version string `json:"version"`
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
