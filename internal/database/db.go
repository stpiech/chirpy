package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var databaseFile string = "database.json"
var DatabaseMux sync.RWMutex

type Chirp struct {
  Id int `json:"id"`
  Body string `json:"body"`
}

type databaseStructure struct {
  Chirps []Chirp `json:"chirps"`
}

func Data() (databaseStructure, error) {
  DatabaseMux.RLock()
  defer DatabaseMux.RUnlock()

  dbJson, err := readDatabaseFile()  
  if err != nil {
    return databaseStructure{}, err
  }

  structuredData := databaseStructure{}
  err = json.Unmarshal(dbJson, &structuredData) 
  if err != nil {
    return databaseStructure{}, err
  }

  return structuredData, nil
}

func WriteChirp(dataToWrite Chirp) (Chirp, error) {
  data, err := Data()
  if err != nil {
    return Chirp{}, err
  }

  DatabaseMux.Lock()
  defer DatabaseMux.Unlock()

  highestId := 0
  for _, chirp := range data.Chirps {
    if chirp.Id > highestId {
      highestId = chirp.Id
    }
  }
  recordId := highestId + 1
  dataToWrite.Id = recordId
  data.Chirps = append(data.Chirps, dataToWrite)

  jsonData, err := json.Marshal(data)

  if err != nil {
    return Chirp{}, err
  }

  writeDatabaseFile(jsonData) 

  return dataToWrite, nil
}

func readDatabaseFile() ([]byte, error) {
  fileData, err := os.ReadFile(databaseFile)

  if err != nil {
    if os.IsNotExist(err) {
      createDatabaseFile()
      fileData, err = os.ReadFile(databaseFile)
      if err != nil {
        return nil, errors.New("Can't connect to DB")
      }
    } else {
      return nil, errors.New("Can't connect to DB")
    }
  }

  return fileData, nil
}

func writeDatabaseFile(data []byte) {
  os.WriteFile(databaseFile, data, 0666)
}

func createDatabaseFile() {
  os.WriteFile(databaseFile, []byte(`{"chirps": []}`), 0666)
}
