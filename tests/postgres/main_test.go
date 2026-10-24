package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkgdb "github.com/survivorbat/go-pkg-db"
)

var createDB pkgdb.CreateFunc

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	createFunc, teardown, err := pkgdb.SetupPostgres(ctx, "postgres:18-alpine", nil)
	if err != nil {
		panic(err)
	}

	createDB = createFunc

	code := m.Run()

	err = teardown()
	if err != nil {
		panic(err)
	}

	cancel()

	os.Exit(code)
}

func TestDatabaseIsCreated(t *testing.T) {
	t.Parallel()
	// Act
	db := createDB(t)

	// Assert
	row := db.QueryRowContext(t.Context(), "select current_database()")
	require.NoError(t, row.Err())

	var result string

	err := row.Scan(&result)
	require.NoError(t, err)

	assert.Equal(t, "testdatabaseiscreated", result)
}

func TestDatabaseIsCreatedInSubTests(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"foo": "testdatabaseiscreatedinsubtestsfoo",
		"bar": "testdatabaseiscreatedinsubtestsbar",
	}

	for name, expected := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			db := createDB(t)

			// Assert
			row := db.QueryRowContext(t.Context(), "select current_database()")
			require.NoError(t, row.Err())

			var result string

			err := row.Scan(&result)
			require.NoError(t, err)

			assert.Equal(t, expected, result)
		})
	}
}
