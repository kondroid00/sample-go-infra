package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	User     string
	Password string
	Protocol string
	DBSchema string
	Option   string
}

var dbs map[string]*gorm.DB

func Init(configs []Config, logWriter gorm.LogWriter) error {
	dbs = make(map[string]*gorm.DB, 8)
	for _, config := range configs {
		connect := config.User + ":" + config.Password + "@" + config.Protocol + "/" + config.DBSchema + "?" + config.Option
		db, err := gorm.Open("mysql", connect)
		if err != nil {
			return err
		}
		if logWriter != nil {
			db.LogMode(true)
			db.SetLogger(&gorm.Logger{
				LogWriter: logWriter,
			})
		}
		dbs[config.DBSchema] = db
	}
	return nil
}

type DB struct {
	db       *gorm.DB
	commited bool
}

type Tx = DB

func GetDB(schema string) *DB {
	return &DB{
		db: dbs[schema].New(),
	}
}

func (s *DB) returnDB(db *gorm.DB) *DB {
	s.db = db
	return s
}

func (s *DB) Where(query interface{}, args ...interface{}) *DB {
	return s.returnDB(s.db.Where(query, args...))
}

func (s *DB) Or(query interface{}, args ...interface{}) *DB {
	return s.returnDB(s.db.Or(query, args...))
}

func (s *DB) Not(query interface{}, args ...interface{}) *DB {
	return s.returnDB(s.db.Not(query, args...))
}

func (s *DB) Limit(limit interface{}) *DB {
	return s.returnDB(s.db.Limit(limit))
}

func (s *DB) Offset(offset interface{}) *DB {
	return s.returnDB(s.db.Offset(offset))
}

func (s *DB) Order(value interface{}, reorder ...bool) *DB {
	return s.returnDB(s.db.Order(value, reorder...))
}

func (s *DB) Select(query interface{}, args ...interface{}) *DB {
	return s.returnDB(s.db.Select(query, args...))
}

func (s *DB) Omit(columns ...string) *DB {
	return s.returnDB(s.db.Omit(columns...))
}

func (s *DB) Group(query string) *DB {
	return s.returnDB(s.db.Group(query))
}

func (s *DB) Having(query interface{}, values ...interface{}) *DB {
	return s.returnDB(s.db.Having(query, values...))
}

func (s *DB) Joins(query string, args ...interface{}) *DB {
	return s.returnDB(s.db.Joins(query, args...))
}

func (s *DB) Scopes(funcs ...func(*DB) *DB) *DB {
	for _, f := range funcs {
		s = f(s)
	}
	return s
}

func (s *DB) Unscoped() *DB {
	return s.returnDB(s.db.Unscoped())
}

func (s *DB) Attrs(attrs ...interface{}) *DB {
	return s.returnDB(s.db.Attrs(attrs...))
}

func (s *DB) Assign(attrs ...interface{}) *DB {
	return s.returnDB(s.db.Assign(attrs...))
}

func (s *DB) First(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.First(out, where...))
}

func (s *DB) Take(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.Take(out, where...))
}

func (s *DB) Last(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.Last(out, where...))
}

func (s *DB) Find(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.Find(out, where...))
}

func (s *DB) Preloads(out interface{}) *DB {
	return s.returnDB(s.db.Preloads(out))
}

func (s *DB) Scan(dest interface{}) *DB {
	return s.returnDB(s.db.Scan(dest))
}

func (s *DB) Pluck(column string, value interface{}) *DB {
	return s.returnDB(s.db.Pluck(column, value))
}

func (s *DB) Count(value interface{}) *DB {
	return s.returnDB(s.db.Count(value))
}

func (s *DB) Related(value interface{}, foreignKeys ...string) *DB {
	return s.returnDB(s.db.Related(value, foreignKeys...))
}

func (s *DB) FirstOrInit(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.FirstOrInit(out, where...))
}

func (s *DB) FirstOrCreate(out interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.FirstOrCreate(out, where...))
}

func (s *DB) Update(attrs ...interface{}) *DB {
	return s.returnDB(s.db.Update(attrs...))
}

func (s *DB) Updates(values interface{}, ignoreProtectedAttrs ...bool) *DB {
	return s.returnDB(s.db.Updates(values, ignoreProtectedAttrs...))
}

func (s *DB) UpdateColumn(attrs ...interface{}) *DB {
	return s.returnDB(s.db.UpdateColumn(attrs...))
}

func (s *DB) UpdateColumns(values interface{}) *DB {
	return s.returnDB(s.db.UpdateColumns(values))
}

func (s *DB) Save(value interface{}) *DB {
	return s.returnDB(s.db.Save(value))
}

func (s *DB) Create(value interface{}) *DB {
	return s.returnDB(s.db.Create(value))
}

func (s *DB) Delete(value interface{}, where ...interface{}) *DB {
	return s.returnDB(s.db.Delete(value, where...))
}

func (s *DB) Raw(sql string, values ...interface{}) *DB {
	return s.returnDB(s.db.Raw(sql, values...))
}

func (s *DB) Exec(sql string, values ...interface{}) *DB {
	return s.returnDB(s.db.Exec(sql, values...))
}

func (s *DB) Model(value interface{}) *DB {
	return s.returnDB(s.db.Model(value))
}

func (s *DB) Table(name string) *DB {
	return s.returnDB(s.db.Table(name))
}

func (s *DB) Debug() *DB {
	return s.returnDB(s.db.Debug())
}

func (s *DB) Begin() *Tx {
	s.commited = false
	return s.returnDB(s.db.Begin())
}

func (s *DB) Commit() error {
	db := s.db.Commit()
	s.commited = db.Error == nil
	return db.Error
}

func (s *DB) RollbackUnlessCommited() error {
	if !s.commited {
		return s.db.Rollback().Error
	}
	return nil
}

func (s *DB) Preload(column string, conditions ...interface{}) *DB {
	return s.returnDB(s.db.Preload(column, conditions...))
}

func (s *DB) Set(name string, value interface{}) *DB {
	return s.returnDB(s.db.Set(name, value))
}

func (s *DB) InstantSet(name string, value interface{}) *DB {
	return s.returnDB(s.db.InstantSet(name, value))
}

func (s *DB) Error() error {
	return s.db.Error
}
