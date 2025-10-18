package models

type ProfileResponse struct {
	ID          int          `json:"id"`
	WowAccounts []WowAccount `json:"wow_accounts"`
}

type WowAccount struct {
	Characters []CharacterResponseSimple `json:"characters"`
}

type CharacterResponseSimple struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Level         int           `json:"level"`
	PlayableClass PlayableClass `json:"playable_class"`
	PlayableRace  PlayableRace  `json:"playable_race"`
	Realm         Realm         `json:"realm"`
}

type PlayableClass struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

type PlayableRace struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

type Realm struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
