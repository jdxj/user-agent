package db

import (
	"testing"
	"time"

	"github.com/jdxj/user-agent/control"
)

func TestMySQL(t *testing.T) {
	if err := MySQL.Ping(); err != nil {
		t.Fatalf("%s", err)
	}
	defer MySQL.Close()

	time.Sleep(10 * time.Second)
	close(control.Stop)
}
