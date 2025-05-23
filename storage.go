package main

import(
	"encoding/json"
	"os"
	"time"
)

type HabitData struct {
	Dates []string `json:"dates"`
	Streak int `json:"streak"`
}

type HabitStore map[string]HabitData

const habitFile = "habits.json"

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

func today() string {
	return time.Now().Format("2006-01-02")
}
