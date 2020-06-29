package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/beevik/guid"
	"github.com/gocraft/dbr"
)

type Model struct {
	CacheKey  string       `json:"-"`
	Db        *dbr.Session `json:"-"`
	TableName string       `json:"-"`
	Tx        *dbr.Tx      `json:"-"`
}

func (this *Model) GetCount(exps map[string]interface{}) (int, error) {
	var count int

	builder := this.Db.Select("COUNT(0)").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadOne(&count)

	return count, err
}

func (this *Model) GetCountWithCreateTime(exps map[string]interface{}) (int, error) {
	var count int

	builder := this.Db.Select("COUNT(0)").From(this.TableName)

	// 用 CreateTime 来影响 MySQL 优化器的选择，达到选择正确的索引目的
	err := this.SelectWhere(builder, exps).OrderDesc("CreateTime").
		Limit(1).
		LoadOne(&count)
	return count, err
}

func (this *Model) GetIds(exps map[string]interface{}) ([]int64, error) {
	var ids []int64
	builder := this.Db.Select("ID").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).Load(&ids)
	return ids, err
}

func (this *Model) Delete(exps map[string]interface{}) error {
	var builder *dbr.DeleteStmt
	if this.Tx != nil {
		builder = this.Tx.DeleteFrom(this.TableName)
	} else {
		builder = this.Db.DeleteFrom(this.TableName)

	}
	_, err := this.DeleteWhere(builder, exps).Exec()

	return err
}

func (this *Model) Insert(params map[string]interface{}) (int64, error) {
	var builder *dbr.InsertStmt
	if this.Tx != nil {
		builder = this.Tx.InsertInto(this.TableName)
	} else {
		builder = this.Db.InsertInto(this.TableName)
	}
	result, err := this.InsertParams(builder, params).Exec()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (this *Model) BatchInsertTx(colums []string, params []interface{}) error {
	if len(params) > 500 {
		return fmt.Errorf("model: Insert data can't more 500, %v existing records", len(params))
	}
	var builder *dbr.InsertStmt
	if this.Tx != nil {
		builder = this.Tx.InsertInto(this.TableName).Columns(colums...)
	} else {
		builder = this.Db.InsertInto(this.TableName).Columns(colums...)
	}

	for _, v := range params {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Slice {
			return errors.New("model: Insert data must is slice type")
		}

		_data := make([]interface{}, 0)

		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			itemVal := item.Interface()

			_data = append(_data, itemVal)
		}

		builder = builder.Values(_data...)
	}

	_, err := builder.Exec()
	return err
}

func (this *Model) IsExist(exps map[string]interface{}, id int64) (bool, error) {
	var value int64

	builder := this.Db.Select("ID").From(this.TableName)
	err := this.SelectWhere(builder, exps).Limit(1).LoadOne(&value)
	if err != nil {
		if err == dbr.ErrNotFound {
			return false, nil
		}
	} else {
		if value == id {
			return false, nil
		}
	}

	return true, err
}

func (this *Model) Update(params map[string]interface{}, exps map[string]interface{}) error {
	var builder *dbr.UpdateStmt
	if this.Tx != nil {
		builder = this.Tx.Update(this.TableName)
	} else {
		builder = this.Db.Update(this.TableName)
	}
	this.UpdateParams(builder, params)

	_, err := this.UpdateWhere(builder, exps).Exec()
	return err
}

func (this *Model) UpdateWithResult(params map[string]interface{}, exps map[string]interface{}) (int64, error) {
	var builder *dbr.UpdateStmt
	if this.Tx != nil {
		builder = this.Tx.Update(this.TableName)
	} else {
		builder = this.Db.Update(this.TableName)
	}
	this.UpdateParams(builder, params)
	result, err := this.UpdateWhere(builder, exps).Exec()
	count, _ := result.RowsAffected()
	return count, err
}

func (this *Model) Increment(base string, value interface{}, field string, step int) error {
	cmd := fmt.Sprintf(`UPDATE %s SET %s = %s + %d WHERE %s = %v`,
		this.TableName, field, field, step, base, value)

	_, err := this.Db.UpdateBySql(cmd).Exec()
	return err
}

// --------------------------------------------------------------------------------

func (this *Model) InsertParams(builder *dbr.InsertBuilder, exps map[string]interface{}) *dbr.InsertBuilder {
	var column []string
	var value []interface{}
	for k, v := range exps {
		column = append(column, k)
		value = append(value, v)
	}
	builder.Columns(column...)
	builder.Values(value...)
	return builder
}

func (this *Model) UpdateParams(builder *dbr.UpdateBuilder, exps map[string]interface{}) *dbr.UpdateBuilder {
	for k, v := range exps {
		builder.Set(k, v)
	}
	return builder
}

func (this *Model) DeleteWhere(builder *dbr.DeleteBuilder, exps map[string]interface{}) *dbr.DeleteBuilder {
	for k, v := range exps {
		values, ok := v.([]interface{})
		if ok {
			builder.Where(k, values...)
		} else {
			builder.Where(k, v)
		}
	}
	return builder
}

func (this *Model) UpdateWhere(builder *dbr.UpdateBuilder, exps map[string]interface{}) *dbr.UpdateBuilder {
	for k, v := range exps {
		values, ok := v.([]interface{})
		if ok {
			builder.Where(k, values...)
		} else {
			builder.Where(k, v)
		}
	}

	return builder
}

func (this *Model) SelectWhere(builder *dbr.SelectBuilder, exps map[string]interface{}) *dbr.SelectBuilder {
	for k, v := range exps {
		switch v.(type) {
		case dbr.Builder:
			builder.Where(v)
		default:
			values, ok := v.([]interface{})
			if ok {
				builder.Where(k, values...)
			} else {
				builder.Where(k, v)
			}
		}
	}

	return builder
}

// --------------------------------------------------------------------------------

func FormatInt(id int64) string {
	return strconv.FormatInt(id, 10)
}

func ConvertInt(n interface{}) int {
	var result int
	switch n.(type) {
	case int:
		result = n.(int)
	case int64:
		v, _ := n.(int64)
		result = int(v)
	case float64:
		v, _ := n.(float64)
		result = int(v)
	case string:
		result, _ = strconv.Atoi(n.(string))
	}
	return result
}

func ConvertInt64(n interface{}) int64 {
	var result int64
	switch n.(type) {
	case int64:
		result = n.(int64)
	case float64:
		v, _ := n.(float64)
		result = int64(v)
	case int:
		i, _ := n.(int)
		result = int64(i)
	case string:
		result, _ = strconv.ParseInt(n.(string), 10, 64)
	default:
	}
	return result
}

func ConvertUint64(n interface{}) uint64 {
	var result uint64
	switch n.(type) {
	case int64:
		result = n.(uint64)
	case float64:
		v, _ := n.(float64)
		result = uint64(v)
	case int:
		i, _ := n.(int)
		result = uint64(i)
	case string:
		result, _ = strconv.ParseUint(n.(string), 10, 64)
	default:
	}
	return result
}

// 字符串（1,2,3,4形式）转数组，baseSize是类型
func StrToIntArray(str string, baseSize int) ([]int8, []int32, []int64) {
	var ids8 []int8
	var ids32 []int32
	var ids64 []int64

	if len(str) > 0 {
		idsArr := strings.Split(str, ",")
		for _, idStr := range idsArr {
			idInt, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				idInt = -100 // 为了兼容特殊字符输入默认给一个不存在的值
			}
			switch baseSize {
			case 8:
				ids8 = append(ids8, int8(idInt))
				break
			case 32:
				ids32 = append(ids32, int32(idInt))
				break
			case 64:
				ids64 = append(ids64, idInt)
				break
			}
		}
	}
	return ids8, ids32, ids64
}

func NewUUID() string {
	return guid.New().String()
}

func NewUniqueID() string {
	return strings.Replace(guid.New().String(), "-", "", -1)
}

func New6AuthCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%v", r.Intn(899999)+100000)
}

func New4BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%04x", r.Intn(65535))
}

func New6BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06x", r.Intn(16777215))
}

func New8BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%08x", r.Intn(4294967295))
}

func NewBatchNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("3%v%v", time.Now().Format("060102150405"), r.Intn(89999)+10000)
}

// --------------------------------------------------------------------------------

func StructToMap(dst interface{}) map[string]interface{} {
	v := reflect.ValueOf(dst)
	t := reflect.Indirect(v).Type()

	var elem reflect.Value

	if reflect.ValueOf(dst).Kind() == reflect.Ptr || reflect.ValueOf(dst).Kind() == reflect.Interface {
		elem = v.Elem()
	} else {
		elem = v
	}

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		if tag == "-" {
			// ignore
			continue
		}

		if elem.Field(i).Type().Kind() == reflect.Struct {
			data[field.Name] = elem.Field(i).Field(0).Field(0).Interface()
			continue
		}
		data[field.Name] = elem.Field(i).Interface()
	}
	return data
}

func MapToStruct(dst interface{}, src map[string]string) error {
	for k, v := range src {
		err := setValue(dst, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------------------

// 生成参数用于数值范围筛选
// 例如：&balance_exp=100,500
func MakeParamWithNumRange(field string, param string, exps map[string]interface{}) {
	if len(strings.TrimSpace(param)) == 0 {
		return
	}

	params := strings.Split(param, ",")
	if len(params) > 0 {
		if len(params[0]) > 0 {
			startParamInt, err := strconv.ParseInt(params[0], 10, 64)
			if err == nil {
				exps[field+">=?"] = startParamInt
			}
		}
		if len(params) == 2 && len(params[1]) > 0 {
			endParamInt, err := strconv.ParseInt(params[1], 10, 64)
			if err == nil {
				exps[field+"<=?"] = endParamInt
			}
		}
	}
}

// 生成参数用于ID数组范围筛选
// 例如：&type=1,2,3,7
func MakeParamWithIds(field string, ids string, exps map[string]interface{}) {
	if len(strings.TrimSpace(ids)) == 0 {
		return
	}
	_, _, arr := StrToIntArray(ids, 64)
	if len(arr) > 0 {
		if len(arr) == 1 {
			exps[field+"=?"] = arr[0]
		} else {
			exps[field+" IN ?"] = arr
		}
	}
}

// 生成参数用于多组筛选条件
// 例如：filter=mobile:18600001234
//      filter=mobile:18600001234,name:王老二...
func MakeParamWithFilter(rec map[string]string, param string, exps map[string]interface{}) error {
	params := strings.Split(param, ",")
	for _, v := range params {
		arr := strings.Split(v, ":")
		if len(arr) == 2 {
			_field := arr[0]
			_value := arr[1]
			if len(_value) > 0 {
				if value, ok := rec[_field]; ok {
					if _field == "id" || _field == "event_id" {
						exps[value] = ConvertInt64(_value)
					} else {
						exps[value] = _value
					}
				}
			}
		} else {
			return errors.New("make param: Parameter format does not match")
		}
	}
	return nil
}

// judge Int param
func ParamInt(param string, defaultValue int) (int, error) {
	if len(param) == 0 {
		return defaultValue, nil
	}
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return paramInt, nil
}

// --------------------------------------------------------------------------------

func setValue(obj interface{}, name string, value string) error {
	elem := reflect.ValueOf(obj).Elem()
	fieldName := elem.FieldByName(name)

	if !fieldName.IsValid() {
		// return fmt.Errorf("models: No such field %v", name)
		return nil
	}

	if !fieldName.CanSet() {
		return fmt.Errorf("models: Cannot set %v field value", name)
	}

	fieldType := fieldName.Type()
	v := reflect.ValueOf(value)
	if fieldType != v.Type() {
		var typeValue interface{}
		switch fieldType.Kind() {
		case reflect.Int:
			typeValue, _ = strconv.Atoi(value)
		case reflect.Int32:
			temp, _ := strconv.ParseInt(value, 10, 64)
			typeValue = int32(temp)
		case reflect.Int64:
			typeValue, _ = strconv.ParseInt(value, 10, 64)
		case reflect.Float32:
			temp, _ := strconv.ParseFloat(value, 64)
			typeValue = float32(temp)
		case reflect.Float64:
			typeValue, _ = strconv.ParseFloat(value, 64)
		case reflect.Struct:
			switch fieldName.Interface().(type) {
			case dbr.NullString:
				typeValue = dbr.NullString{sql.NullString{value, true}}
			case dbr.NullInt64:
				nullInt64, _ := strconv.ParseInt(value, 10, 64)
				typeValue = dbr.NullInt64{sql.NullInt64{nullInt64, true}}
			default:
				fieldName.Field(0).Set(reflect.ValueOf(sql.NullString{value, true}))
				return nil
			}
		default:
			return fmt.Errorf("models: Type didn't match %v", name)
		}
		v = reflect.ValueOf(typeValue)
	}
	fieldName.Set(v)
	return nil
}
