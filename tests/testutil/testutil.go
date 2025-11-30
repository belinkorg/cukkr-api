package testutil

import (
	"bLink-app/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a mock database for testing
func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

// SetupTestDBWithTransaction creates a mock database for testing with transaction
func SetupTestDBWithTransaction(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true, // Important!
	})
	assert.NoError(t, err)

	// Expect Begin and Commit for transaction
	mock.ExpectBegin()
	mock.ExpectCommit()

	return db, mock
}

// GetTestLogger creates a test logger
func GetTestLogger() *logger.Logger {
	return logger.NewLogger("error")
}
