package models

type CharacterResponseDetailed struct {
	Guild Guild `json:"guild"`
}

type Guild struct {
	Name string `json:"name"`
}
