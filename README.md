# ⛁ Pkg DB

This library contains a test helper to set up a package-level
database server using [testcontainers](https://golang.testcontainers.org/modules/postgres/),
and provides a function to create databases within it for individual tests.
This is an improvement over [survivorbat/go-db-suite](https://github.com/survivorbat/go-db-suite)
which relies on testify's suites and as a result makes tests run in sequentially.
By using TestMain, tests can be run in parallel.

## ⬇️ Installation

```bash
go get github.com/survivorbat/go-pkg-db
```

## 📋 Usage

```go
import (
 "context"
 "os"
 "testing"

 pkgdb "github.com/survivorbat/go-pkg-db"
)

var createDB pkgdb.CreateFunc

func TestMain(m *testing.M) {
  ctx, cancel := context.WithCancel(context.Background())

  // Setup the postgres server
  createFunc, teardown, err := pkgdb.SetupPostgres(ctx, "postgres:18-alpine", nil)
  if err != nil {
    panic(err)
  }

  // Set the creation func as a global, allowing tests to use it
  createDB = createFunc

  // Run the unit tests
  code := m.Run()

  // Cleanup
  err = teardown()
  if err != nil {
    panic(err)
  }

  cancel()

  // Exit the process with the test's code
  os.Exit(code)
}

```

## 🔭 Plans

- More databases
