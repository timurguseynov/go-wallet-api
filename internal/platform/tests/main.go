package tests

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/timurguseynov/go-wallet-api/internal/platform/db"
	"github.com/timurguseynov/go-wallet-api/internal/platform/rest"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// Test owns state for running/shutting down tests.
type Test struct {
	Log      *log.Logger
	MasterDB *db.DB
}

// New is the entry point for tests.
func New() *Test {
	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// ============================================================
	// Start database

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal("main : couldn't connect to database", err)
	}

	mustSeed(context.TODO(), dbConn)

	return &Test{Log: log, MasterDB: dbConn}
}

// TearDown is used for shutting down tests. Calling this should be
// done in a defer immediately after calling New.
func (t *Test) TearDown() {}

// Recover is used to prevent panics from allowing the test to cleanup.
func Recover(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal("Unhandled Exception:", string(debug.Stack()))
	}
}

// Context returns an app level context for testing.
func Context() context.Context {
	values := rest.Values{
		TraceID: uuid.New(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), rest.KeyValues, &values)
}
