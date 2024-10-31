package main

type characterJson struct {
	Name           string            `json:"name"`
	Stats          map[string]int    `json:"skills"`
	Comabt_skills  map[string]string `json:"cSkills"`
	Dalogue_skills map[string]string `json:"dSkills"`
}

type character struct {
	Name           string
	Stats          map[string]int
	CombatSkills   map[string]string
	DialogueSkills map[string]string
}

type skill struct {
	skillName string
	effect    string
}

//func character_loader(characterJson) character {

//}
