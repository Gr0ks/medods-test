package models

type AccessPair struct {
	AccessKey string `json:"accessKey"`
	RefreshKey string `json:"refreshKey"`
}