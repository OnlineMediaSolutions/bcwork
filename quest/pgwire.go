package quest

import (
	"database/sql/driver"
	"fmt"
	"net"
	"time"

	"github.com/spf13/viper"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/jmoiron/sqlx"
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
	db, err = Connect(env)
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	return db.Close()
}

func dbConfig(env string, key string) string {
	return fmt.Sprintf("quest.%s.%s", env, key)
}

func Connect(conn string) (*sqlx.DB, error) {
	fmt.Println(viper.GetString(dbConfig(conn, "host")))
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
