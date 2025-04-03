package main

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"log"
	"os"
	"path/filepath"
)

func AddTileSetTypeProperty(filePath string) {
	// Define the JSON file path
	println("updating", filePath)
	// Read the existing JSON file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal JSON into a map
	var jsonData map[string]interface{}
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	tiles, ok := jsonData["tiles"].([]interface{})
	if !ok {
		jsonData["properties"] = []map[string]interface{}{
			{
				"name":  "tileSetType",
				"type":  "string",
				"value": "uniform",
			},
		}
		println("updated uniform tileset")
	}

	// Add a new field
	if len(tiles) > 0 {
		if _, exists := tiles[0].(map[string]interface{})["animation"]; exists {
			jsonData["properties"] = []map[string]interface{}{
				{
					"name":  "tileSetType",
					"type":  "string",
					"value": "animated",
				},
			}
		} else if len(tiles) > 0 {
			jsonData["properties"] = []map[string]interface{}{
				{
					"name":  "tileSetType",
					"type":  "string",
					"value": "dynamic",
				},
			}
		}
	}

	// Marshal the updated JSON
	updatedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write the updated JSON back to the file
	if err := os.WriteFile(filePath, updatedJSON, os.ModePerm); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("✅ JSON file updated successfully!")
}

func AddAdHocPropertyToMap(layerFilter, typeFilter string, newProp map[string]interface{}) {
	// Define the JSON file path
	// Read the existing JSON file
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	file, err := os.ReadFile(workDir + "/assets/map/town1Map.json")

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal JSON into a map
	var mapData map[string]interface{}

	if err := json.Unmarshal(file, &mapData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	layers, ok := mapData["layers"].([]interface{})
	if !ok {
		log.Fatalf("Failed to convert to list data structure: %v", err)
	}

	var targetLayer map[string]interface{}

	for _, layer := range layers {
		layerMap, ok := layer.(map[string]interface{})
		if !ok {
			println("failed to convert to map")
		}

		// Check if the "name" field matches the target
		if name, ok := layerMap["name"].(string); ok && name == layerFilter {
			targetLayer = layerMap
		}

	}

	objects := targetLayer["objects"].([]interface{})

	for _, obj := range objects {
		objMap := obj.(map[string]interface{})
		if objMap["type"] != typeFilter {
			continue
		}
		if props, exists := objMap["properties"].([]interface{}); exists {
			println("appending prop to", objMap["Name"])
			objMap["properties"] = append(props, newProp)
		} else {
			objMap["properties"] = []interface{}{newProp}
		}
	}

	// Marshal the updated JSON

	updatedJSON, err := json.Marshal(mapData)
	// Write the updated JSON back to the file

	if err := os.WriteFile(workDir+"/assets/map/town1Map.json", updatedJSON, os.ModePerm); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("✅ JSON file updated successfully!")
}

func processListofFiles() {
	files, err := os.ReadDir("tilesets/assets/map/town1Map.json")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join("tilesets/", file.Name()) // Process only JSON files
			AddTileSetTypeProperty(filePath)
		}
	}
}

func main() {

	exitDoorProp := map[string]interface{}{
		"Name":     "sprite",
		"drawType": "string",
		"Value":    "none",
	}

	AddAdHocPropertyToMap("doors", "exitDoor", exitDoorProp)

}

func UpdateTileMapProperty(contents []byte, inputProp gameObjects.PropertiesJSON, layerName string, objectName string) {

	var tileMap struct {
		Layers []gameObjects.TileMapLayer `json:"layers"`
		Raw    json.RawMessage            `json:"-"` // Preserve extra JSON fields
	}

	if err := json.Unmarshal(contents, &tileMap); err != nil {
		log.Fatal("Error reading tileMap:", err)
	}

	var targetLayer *gameObjects.TileMapLayer
	for _, layer := range tileMap.Layers {
		if layer.Name == layerName {
			targetLayer = &layer
		}
	}

	if targetLayer == nil {
		log.Fatal("targetLayer is null")
	}

	for i := range targetLayer.Objects {
		println(targetLayer.Objects[i].Name)
		if targetLayer.Objects[i].Type == objectName {
			propThere := false
			for j := range targetLayer.Objects[i].Properties {
				if targetLayer.Objects[i].Properties[j].Name == inputProp.Name {
					targetLayer.Objects[i].Properties[j] = inputProp
					propThere = true
				}
				println("updating existing prop to new prop")
			}
			if !propThere {
				targetLayer.Objects[i].Properties = append(targetLayer.Objects[i].Properties, inputProp)
				println("adding new prop")
			}
		}
	}

	updatedJSON, err := json.MarshalIndent(tileMap, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	// Write the updated JSON back to the file

	workDir, err := os.Getwd()

	if err := os.WriteFile(workDir+"/assets/map/town1Map.json", updatedJSON, os.ModePerm); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

}
