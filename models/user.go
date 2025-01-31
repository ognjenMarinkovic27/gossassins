package models

type User struct {
	Uid      string `json:"uid"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
}
