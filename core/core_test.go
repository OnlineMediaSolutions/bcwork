package core

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initTest()

	code := m.Run()

	os.Exit(code)

}

//
//func initTest() (*dockertest.Resource, context.Context) {
//	pool = testutils.SetupDockerTestPool()
//	pg := testutils.SetupDB(pool)
//
//	err := bcdb.DB().Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background()
//	createTables(bcdb.DB())
//	return pg, ctx
//}
