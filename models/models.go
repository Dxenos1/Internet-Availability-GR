package models

type State struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Municipality struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	StateId int    `json:"stateId"`
}

type Prefecture struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	StateId        int    `json:"stateId"`
	MunicipalityId int    `json:"municipalityId"`
}

type PackageInfo struct {
	Telephone string `json:"telephone"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	Speed     string `json:"speed"`
}

type Input struct {
	Telephones []string `yaml:"telephones"`
}
