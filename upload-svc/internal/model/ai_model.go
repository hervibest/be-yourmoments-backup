package model

type ProcessFacecam struct {
	UserId    string
	CreatorId string
	FileURL   string
}

type ProcessPhoto struct {
	PhotoId          string
	CreatorId        string
	FileURL          string
	OriginalFilename string
}
