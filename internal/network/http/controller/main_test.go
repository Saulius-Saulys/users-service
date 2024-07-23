package controller_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Saulius-Saulys/users-service/internal/testutils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// Address to broker shared between all tests in this package
var gormConn *gorm.DB

func TestMain(m *testing.M) {
	connectionString, pool, resource, err := testutils.StartDatabase()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	usersDB, err := testutils.PrepareUsersDB(connectionString, true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	gormConn, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: zapgorm2.New(zap.NewNop()), NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "users.",
		SingularTable: false,
	}})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	exitVal := m.Run()
	err = pool.Purge(resource)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = usersDB.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(exitVal)
}
