package install

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
)

type SqlxInstall struct {
	Operator Operator
	isDrop   bool
	db       *sql.DB
}

type SqlxOperator struct {
	db *sql.DB
}

func NewInstall(host string, port int, user, password, database string) (*SqlxInstall, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, database)
	return NewSourceInstall(dataSourceName)
}

func NewSourceInstall(source string) (*SqlxInstall, error) {
	DB, err := sql.Open("mysql", source)
	if err != nil {
		fmt.Println("connect mysql error", err, DB)
		return nil, err
	}
	ins := &SqlxInstall{
		isDrop: false,
		db:     DB,
	}
	ins.Operator = newSqlxOperator(ins.isDrop, ins.db)
	return ins, nil
}

func (sqx *SqlxInstall) SetOperator(operator Operator) {
	sqx.Operator = operator
}

func (sqx *SqlxInstall) SetDrop(b bool) {
	sqx.isDrop = b
}

func (sqx *SqlxInstall) Exec(fileName string) error {
	defer sqx.db.Close()
	if f, err := os.Stat(fileName); err != nil && f.Size() <= 0 {
		return err
	}
	fi, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	return sqx.ExecReder(fi)
}

func (sqx *SqlxInstall) ExecReder(r io.Reader) error {
	tokens := sqlparser.NewTokenizer(r)
	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		// exec
		if err := sqx.exec(stmt); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (sqx *SqlxInstall) exec(stmt sqlparser.Statement) error {
	switch stmt.(type) {
	case *sqlparser.Select:
		if err := sqx.Operator.Select(sqlparser.String(stmt)); err != nil {
			return fmt.Errorf("exec select error:%v", err)
		}
	case *sqlparser.Update:
		if err := sqx.Operator.Update(sqlparser.String(stmt)); err != nil {
			return fmt.Errorf("exec update error:%v", err)
		}
	case *sqlparser.Insert:
		if err := sqx.Operator.Insert(sqlparser.String(stmt)); err != nil {
			return fmt.Errorf("exec insert error:%v", err)
		}
	case *sqlparser.DDL:
		ddl := stmt.(*sqlparser.DDL)
		if ddl.Action == "create" {
			if err := sqx.Operator.Create(sqlparser.String(stmt)); err != nil {
				return fmt.Errorf("exec create error:%v", err)
			}
		} else if ddl.Action == "drop" {
			if sqx.isDrop == true {
				if err := sqx.Operator.Drop(sqlparser.String(stmt)); err != nil {
					return fmt.Errorf("exec drop error:%v", err)
				}
			}
		} else {
			log.Println("action:", ddl.Action)
		}
	default:
	}
	return nil
}

func newSqlxOperator(isDrop bool, db *sql.DB) SqlxOperator {
	return SqlxOperator{
		db: db,
	}
}

func (f SqlxOperator) Insert(data string) error {
	if _, err := f.db.Exec(data); err != nil {
		return err
	}
	return nil
}

func (f SqlxOperator) Update(data string) error {
	if _, err := f.db.Exec(data); err != nil {
		return err
	}
	return nil
}

func (f SqlxOperator) Create(data string) error {
	if _, err := f.db.Exec(data); err != nil {
		return err
	}
	return nil
}

func (f SqlxOperator) Select(data string) error {
	if _, err := f.db.Exec(data); err != nil {
		return err
	}
	return nil
}

func (f SqlxOperator) Drop(data string) error {
	if _, err := f.db.Exec(data); err != nil {
		return err
	}
	return nil
}
