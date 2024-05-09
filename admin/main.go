package main

import (
	"molta9a/server/db"
	"os"

	"encoding/json"
	"log"

	"gorm.io/gorm/clause"
)

func main() {
	db.Init()

	updateTracks()
	updateQuests()
	updateModules()
}

func updateTracks() {
	content, err := os.ReadFile("./admin/tracks.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var tracks []db.Track
	err = json.Unmarshal(content, &tracks)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&tracks)

	log.Println("updated tracks")
}

func updateQuests() {
	content, err := os.ReadFile("./admin/quests.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var quests []db.Quest
	err = json.Unmarshal(content, &quests)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&quests)

	log.Println("updated quests")
}

func updateModules() {
	content, err := os.ReadFile("./admin/modules.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var modules []db.Module
	err = json.Unmarshal(content, &modules)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&modules)

	log.Println("updated modules")
}
