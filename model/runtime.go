package model

type Runtime struct {
	Id              string `json:"id"`
	Label           string `json:"label"`
	LatestVersionId string `json:"latest_version_id"`
	Author          Author `json:"author"`
	Description     string `json:"description"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Url   string `json:"url,omitempty"`
}
