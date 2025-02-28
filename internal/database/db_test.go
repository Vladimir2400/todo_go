package db

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	InitDB()
	err := Migrate()
	if err != nil {
		t.Fatalf("DB migration is falled: %s", err)
	}
}
