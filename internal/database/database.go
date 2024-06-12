package database

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
    mu sync.Mutex
    conn *gorm.DB
}

func Connect(path string) Database {
    conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    log.Println("Connected to db.")
    return Database{ conn: conn, }
}
