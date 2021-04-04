package install

import (
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/xwb1989/sqlparser"
)

type GormInstall struct {
	Operator Operator
	isDrop   bool
	db       *gorm.DB
}

type GormOperator struct {
	db *gorm.DB
}

func NewGormInstall(host string, port int, user, password, database string) (*GormInstall, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, database)
	return NewGormSourceInstall(dataSourceName)
}

func NewGormSourceInstall(source string) (*GormInstall, error) {
	db, err := gorm.Open("mysql", source)
	if err != nil {
		fmt.Println("connect mysql error", err, db)
		return nil, err
	}
	ins := &GormInstall{
		isDrop: false,
		db:     db,
	}
	ins.Operator = newGormOperator(ins.isDrop, ins.db)
	return ins, nil
}

func (sqx *GormInstall) SetOperator(operator Operator) {
	sqx.Operator = operator
}

func (sqx *GormInstall) SetDrop(b bool) {
	sqx.isDrop = b
}

func (sqx *GormInstall) Exec(fileName string) error {
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

func (sqx *GormInstall) ExecReder(r io.Reader) error {
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

func (sqx *GormInstall) exec(stmt sqlparser.Statement) error {
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

func newGormOperator(isDrop bool, db *gorm.DB) GormOperator {
	return GormOperator{
		db: db,
	}
}

func (f GormOperator) Insert(data string) error {
	if err := f.db.Exec(data).Error; err != nil {
		return err
	}
	return nil
}

func (f GormOperator) Update(data string) error {
	if err := f.db.Exec(data).Error; err != nil {
		return err
	}
	return nil
}

func (f GormOperator) Create(data string) error {
	if err := f.db.Exec(data).Error; err != nil {
		return err
	}
	return nil
}

func (f GormOperator) Select(data string) error {
	if err := f.db.Exec(data).Error; err != nil {
		return err
	}
	return nil
}

func (f GormOperator) Drop(data string) error {
	if err := f.db.Exec(data).Error; err != nil {
		return err
	}
	return nil
}
