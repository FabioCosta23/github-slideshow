package p2k

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	pgx "github.com/jackc/pgx/v4/stdlib"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

const (
	MOVEMENT_TYPE_TRANSFER           string = "T"
	STOCK_TRANSIT_TOLERANCE_INTERVAL int    = 15
)

var dbInstance *sql.DB

func Setup() (err error) {
	if dbInstance != nil {
		return nil
	}

	dbInstance, err = setupDbInstance()
	return
}

func setupDbInstance() (*sql.DB, error) {
	sqltrace.Register("pgx", &pgx.Driver{}, sqltrace.WithServiceName("stock-worker-db"))

	fmt.Println("CONN STRING: ", getConnectionString())

	db, err := sqltrace.Open("pgx", getConnectionString())
	if err != nil {
		return nil, err
	} else {
		fmt.Println("CONCETADO STOCK")
	}

	db.SetMaxIdleConns(GetEnvInt("DB_MAX_IDLE_CONNS"))
	db.SetMaxOpenConns(GetEnvInt("DB_MAX_OPEN_CONNS"))
	db.SetConnMaxLifetime(time.Duration(GetEnvInt("DB_CONN_MAX_LIFETIME")) * time.Second)

	// Open doesn't open a connection. Validate DSN data with:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getConnectionString() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_NAME"))
}

type DatabaseCommons interface {
	BeginTransaction() (Transaction, error)
}

type DatabaseCommonsImpl struct{}

func (bt *DatabaseCommonsImpl) BeginTransaction() (Transaction, error) {
	tx, err := dbInstance.Begin()
	if err != nil {
		return nil, err
	}

	return TransactionImpl{
		tx: tx,
	}, err
}

type Transaction interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

type TransactionImpl struct {
	tx *sql.Tx
}

func (t TransactionImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

func (t TransactionImpl) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t TransactionImpl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t TransactionImpl) Commit() error {
	return t.tx.Commit()
}

func (t TransactionImpl) Rollback() error {
	return t.tx.Rollback()
}

func GetEnvInt(configKey string) int {
	i := 0
	i, _ = strconv.Atoi(GetEnv(configKey))
	return i
}
