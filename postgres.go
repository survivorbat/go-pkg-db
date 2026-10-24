package pkgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"testing"

	_ "github.com/lib/pq" // Ensure postgres is available
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

// CreateFunc describes a function to create a fresh database for the test
type CreateFunc func(t *testing.T) *sql.DB

// SetupPostgres will spin upo a postgres container using testcontainers and return a
// CreateFunc that can be used in tests to create fresh databases in the postgres container.
// The second argument is the cancellation function.
func SetupPostgres(ctx context.Context, image string, callback func(*sql.DB)) (CreateFunc, func() error, error) {
	container, err := postgrescontainer.Run(ctx, image, postgrescontainer.BasicWaitStrategies())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run postgres container: %w", err)
	}

	connectionString := container.MustConnectionString(ctx, "sslmode=disable")

	connectionURL, err := url.Parse(connectionString)
	if err != nil {
		termErr := container.Terminate(ctx)
		if termErr != nil {
			err = errors.Join(err, termErr)
		}

		return nil, nil, fmt.Errorf("failed to parse connection URL: %w", err)
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		termErr := container.Terminate(ctx)
		if termErr != nil {
			err = errors.Join(err, termErr)
		}

		return nil, nil, fmt.Errorf("failed to open connection: %w", err)
	}

	cancel := func() error {
		return errors.Join(
			db.Close(),
			container.Terminate(ctx),
		)
	}

	//nolint:contextcheck // The test's context is used
	return NewPostgresFactory(db, connectionURL, callback), cancel, nil
}

// NewPostgresFactory creates a new db creation func to create databases for tests. The final callback may be
// added to apply operations on the database after its creation.
func NewPostgresFactory(rootDB *sql.DB, connectionURL *url.URL, callback func(*sql.DB)) CreateFunc {
	return func(t *testing.T) *sql.DB {
		t.Helper()

		db, cancel, err := ensureDatabase(t.Context(), rootDB, connectionURL, t.Name())
		if err != nil {
			t.Fatal("failed to create a database", err)
		}

		if callback != nil {
			callback(db)
		}

		t.Cleanup(func() { _ = cancel })

		return db
	}
}
