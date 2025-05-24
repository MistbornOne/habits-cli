package main

import(
	"encoding/json"
	"os"
	"path/filepath"
)

type HabitData struct {
	Dates []string `json:"dates"`
	Streak int `json:"streak"`
	Longest int `json:"longest"`
}

type HabitStore map[string]HabitData

var habitFile string

func init() {
	exePath, err := os.Executable()
	if err != nil {
		panic("Could not determine executable path: " + err.Error())
	}

	dir := filepath.Dir(exePath)
	habitFile = filepath.Join(dir, "habits.json")
}

func loadHabits() (HabitStore, error) {
	data := HabitStore{}

	file, err := os.ReadFile(habitFile)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil // first run
		}
		return nil, err
	}

	err = json.Unmarshal(file, &data)
	return data, err
}

func saveHabits(data HabitStore) error {
	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(habitFile, bytes, 0644)
}

