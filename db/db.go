package db

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"time"
)

type Config struct {
	User            string
	Password        string
	Protocol        string
	DBSchema        string
	Option          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

var dbs map[string]*sql.DB

func Init(configs []Config, debugMode bool) error {
	dbs = make(map[string]*sql.DB, 8)
	for _, config := range configs {
		connect := config.User + ":" + config.Password + "@" + config.Protocol + "/" + config.DBSchema + "?" + config.Option
		db, err := sql.Open("mysql", connect)
		if err != nil {
			return err
		}
		db.SetMaxIdleConns(config.MaxIdleConns)
		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetConnMaxLifetime(config.ConnMaxLifetime)

		boil.SetDB(db)
		dbs[config.DBSchema] = db
	}
	if debugMode {
		boil.DebugMode = true
		boil.DebugWriter = &Logger{}
	}
	return nil
}

func Close() error {
	for _, v := range dbs {
		err := v.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

type (
	DB struct {
		cx boil.ContextExecutor
	}

	Tx = DB
)

func GetDB(schema string) *DB {
	return &DB{
		cx: dbs[schema],
	}
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.cx.Exec(query, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.cx.Query(query, args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.QueryRow(query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.cx.ExecContext(ctx, query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.cx.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.cx.QueryRowContext(ctx, query, args...)
}

func (d *DB) Begin(ctx context.Context, options *sql.TxOptions) (*Tx, error) {
	db, ok := d.cx.(*sql.DB)
	if !ok || db == nil {
		return nil, errors.New("cannot get *sql.DB")
	}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return &Tx{
		cx: tx,
	}, nil
}

func (d *DB) Commit() error {
	if tx, ok := d.cx.(boil.ContextTransactor); ok && tx != nil {
		return tx.Commit()
	}
	return errors.New("cannot commit")
}

func (d *DB) Rollback() error {
	if tx, ok := d.cx.(boil.ContextTransactor); ok && tx != nil {
		return tx.Rollback()
	}
	return nil
}

func (d *DB) RollbackUnlessCommitted() error {
	err := d.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}
