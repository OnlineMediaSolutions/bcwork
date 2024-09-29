package bcdb

import (
	"database/sql/driver"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}

var db *sqlx.DB

// DB Get DB object
func DB() *sqlx.DB {
	return db
}

// Env Get current environment object
func Env() string {
	return currentEnv
}

var currentEnv string

func InitDB(env string) error {

	currentEnv = env
	var err error
	db, err = connect(env)
	if err != nil {
		return err
	}

	return nil
}

// InitTestDB used only for integration tests
func InitTestDB(connString string) error {
	var err error
	db, err = sqlx.Connect("postgres", connString)
	if err != nil {
		return errors.Wrapf(err, "failed to conect to postgres instance")
	}
	return nil
}

func CloseDB() error {
	return db.Close()
}

func dbConfig(env string, key string) string {
	return fmt.Sprintf("db.%s.%s", env, key)
}

func connect(conn string) (*sqlx.DB, error) {

	sslMode := "disable"
	if viper.GetBool(dbConfig(conn, "ssl_mode")) {
		sslMode = "require"
	}

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString(dbConfig(conn, "host")),
		viper.GetInt(dbConfig(conn, "port")),
		viper.GetString(dbConfig(conn, "user")),
		viper.GetString(dbConfig(conn, "password")),
		viper.GetString(dbConfig(conn, "db")), sslMode)
	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to conect to postgres instance")
	}
	return db, nil

}
