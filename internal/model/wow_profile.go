package model

type WowProfile struct {
	ID          int          `json:"id"`
	WowAccounts []WowAccount `json:"wow_accounts"`
}

type WowAccount struct {
	Characters []CharacterResponseSimple `json:"characters"`
}

type CharacterResponseSimple struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Realm Realm  `json:"realm"`
}

type Realm struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CharacterResponseDetailed struct {
	Guild Guild `json:"guild"`
}

type Guild struct {
	Name string `json:"name"`
}
