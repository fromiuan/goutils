package tools

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/micro/go-micro/metadata"
)

// environment models
// --------------------------------------------------------------------------------
const (
	ENVIRONMENT_MODEL_DEVELOPMENT = "development" // 开发模式
	ENVIRONMENT_MODEL_PRE         = "pre-release" // 预发布模式
	ENVIRONMENT_MODEL_PRODUCTION  = "production"  // 线上模式
	ENVIRONMENT_MODEL_TESTING     = "testing"     // 测试模式

)

// 分转换元格式 保留2位小数
func CentToYuan(cent int64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(cent)*0.01), 64)
	return value
}

// f:表达式值 n:保留精度位数
func RoundYuanToCent(f float64, n int) int64 {
	n10 := math.Pow10(n)
	floatValue := math.Trunc((f+0.5/n10)*n10) / n10
	return int64(floatValue * float64(100))
}

func Md5(src string, short bool) string {
	h := md5.New()
	h.Write([]byte(src))
	data := hex.EncodeToString(h.Sum(nil))
	if short {
		return data[8:24]
	}
	return data
}

func GetArrayWithCSV(str [][]string) map[string]string {
	if str == nil || len(str) <= 1 {
		return nil
	}
	result := make(map[string]string)
	columns := strings.Split(str[0][0], "|")
	for i := 1; i < len(str); i++ {
		items := strings.Split(str[i][0], "|")
		if len(items) == len(columns) {
			for n := 0; n < len(items); n++ {
				result[columns[n]] = items[n]
			}
		}
	}
	return result
}

func GetLocalIPAddr() (string, error) {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

func FiltrationHTML(str string) string {
	if len(str) == 0 {
		return str
	}
	re, _ := regexp.Compile(`style=\\"(.[^"]*)\\"`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`style='(.[^"]*)'`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`class=\\"(.[^"]*)\\"`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`class='(.[^"]*)'`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile("\\<br[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "\n")

	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "")

	return str
}

func TrimIPAddr(addr string) string {
	return strings.Split(addr, ":")[0]
}

func ClientIP(ctx context.Context) string {
	meta, ok := metadata.FromContext(ctx)
	if !ok {
		return ""
	}

	if addr, ok := meta["X-Forwarded-For"]; ok && addr != "" {
		return strings.TrimSpace(strings.Split(addr, ",")[0])
	}

	if addr, ok := meta["X-Real-Ip"]; ok && addr != "" {
		return addr
	}

	if addr, ok := meta["Remote"]; ok && addr != "" {
		return strings.TrimSpace(strings.Split(addr, ":")[0])
	}

	if addr, ok := meta["X-Appengine-Remote-Addr"]; ok && addr != "" {
		return strings.TrimSpace(strings.Split(addr, ":")[0])
	}

	return ""
}

func IsEmail(email string) bool {
	if len(email) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`^[a-zA-Z0-9_\.]+@\w{2,20}\.\w{2,3}$`)
		return reg.MatchString(email)
	}
}

func IsMobile(number string) bool {
	if len(number) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`0?1[3456789]{1}[0-9]{9}`)
		return reg.MatchString(number)
	}
}

func IsTelephone(number string) bool {
	if len(number) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`(0\d{2})?\d{8}|(0\d{3})?\d{8}|(0\d{3})?\d{7}`)
		return reg.MatchString(number)
	}
}

func IsIDCard(id string) bool {
	id = strings.ToUpper(id)
	if len(id) != 15 && len(id) != 18 {
		return false
	}
	r := regexp.MustCompile("(\\d{15})|(\\d{17}([0-9]|X))")
	if !r.MatchString(id) {
		return false
	}
	if len(id) == 15 {
		tm2, _ := time.Parse("01/02/2006", string([]byte(id)[8:10])+"/"+string([]byte(id)[10:12])+"/"+"19"+string([]byte(id)[6:8]))
		if tm2.Unix() <= 0 {
			return false
		}
		return true
	} else {
		tm2, _ := time.Parse("01/02/2006", string([]byte(id)[10:12])+"/"+string([]byte(id)[12:14])+"/"+string([]byte(id)[6:10]))
		if tm2.Unix() <= 0 {
			return false
		}
		// 检验18位身份证的校验码是否正确。
		// 校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。
		arrInt := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		arrCh := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
		sign := 0
		for k, v := range arrInt {
			intTemp, _ := strconv.Atoi(string([]byte(id)[k : k+1]))
			sign += intTemp * v
		}
		n := sign % 11
		valNum := arrCh[n]
		if valNum != string([]byte(id)[17:18]) {
			return false
		}
		return true
	}
}

// 是否是中文
func IsChinese(str string) bool {
	var hzRegexp = regexp.MustCompile("^[\u4E00-\u9FA5A]+$")
	return hzRegexp.MatchString(str)
}

// 是否含中文
func IsContainChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts[`Han`], r) {
			return true
		}
	}
	return false
}

func getCurrPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func GetPath(name string) string {
	return getCurrPath() + "/" + name
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func IsFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 检查目录是否存在，否则并自动创建
func CheckDir(path string) bool {
	if IsDirExists(path) == false {
		var npath string
		var i int

		paths := strings.Split(path, "/")
		for i = 0; i < len(paths)-2; i++ {
			npath += paths[i] + "/"
		}
		npath += paths[i]
		CheckDir(npath)
		os.Mkdir(path, 0755)
	}
	return true
}

// 检查颜色值是否正确
func CheckColor(color string) bool {
	if len(color) != 7 || color[:1] != "#" {
		return false
	}
	return true
}

func DeletePath(path string) error {
	if IsDirExists(path) == false {
		return nil
	}
	return os.RemoveAll(path)
}

func UnsignedParams(key, timestamp, version string, queryParams url.Values) url.Values {
	params := url.Values{
		"auth_key":       {key},
		"auth_timestamp": {timestamp},
		"auth_version":   {version},
	}

	if queryParams != nil {
		for k, v := range queryParams {
			params[k] = v
		}
	}
	return params
}

func UnescapeUrl(_url url.Values) string {
	unesc, _ := url.QueryUnescape(_url.Encode())
	return unesc
}

func ValidateBinary(limits int, action int) bool {
	binaryList := []int{1, 2, 4, 8, 16, 32}
	if limits == -1 {
		return true
	}
	if action < 1 || action > len(binaryList)-1 {
		return false
	}
	if action == len(binaryList)-1 {
		return limits == binaryList[action]
	}
	return (limits & binaryList[action-1]) >= binaryList[action-1]
}

func InArray(array []string, str string) bool {
	for _, v := range array {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

func InetNtoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func InetAton(ipnr string) int64 {
	bits := strings.Split(strings.Split(ipnr, ":")[0], ".")
	if len(bits) < 4 {
		return 0
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// 字符串以分隔符转化为切片，并且每个值都加上前缀
func StringToSlice(s string, sep string, pre string) []string {
	s = strings.TrimSpace(s)

	arr := make([]string, 0)
	if s != "" {
		arr = strings.Split(s, sep)
		for i, v := range arr {
			arr[i] = pre + v
		}
	}

	return arr
}

// 四舍五入
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Round(f*n10) / n10
}

// 拼接字符串
func MergeString(args ...string) string {
	buf := bytes.Buffer{}
	argsLen := len(args)
	for i := 0; i < argsLen; i++ {
		buf.WriteString(args[i])
	}
	return buf.String()
}

func MapToXml(params map[string]string) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

func MapToXmlNotNil(params map[string]string, isCdata bool) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		if len(v) > 0 {
			buf.WriteString(`<`)
			buf.WriteString(k)
			buf.WriteString(`>`)
			if isCdata {
				buf.WriteString(`<![CDATA[`)
			}
			buf.WriteString(v)
			if isCdata {
				buf.WriteString(`]]>`)
			}
			buf.WriteString(`</`)
			buf.WriteString(k)
			buf.WriteString(`>`)
		}
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

func MapToXmlIsCDATA(params map[string]string, isCdata bool) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`>`)
		if isCdata {
			buf.WriteString(`<![CDATA[`)
		}
		buf.WriteString(v)
		if isCdata {
			buf.WriteString(`]]>`)
		}
		buf.WriteString(`</`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

func XmlToMap(xmlStr string) map[string]string {

	params := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params[key] = value
			}
		}
	}

	return params
}

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

func ConvertInt32(n interface{}) int32 {
	var result int32
	switch n.(type) {
	case int:
		result = n.(int32)
	case int64:
		v, _ := n.(int64)
		result = int32(v)
	case float64:
		v, _ := n.(float64)
		result = int32(v)
	case string:
		v, _ := strconv.Atoi(n.(string))
		result = int32(v)
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
			return errors.New(fmt.Sprintf("Map to struct: %v", err))
		}
	}
	return nil
}

func setValue(obj interface{}, name string, value string) error {
	elem := reflect.ValueOf(obj).Elem()
	fieldName := elem.FieldByName(name)

	if !fieldName.IsValid() {
		//return fmt.Errorf("No such field %v", name)

		return nil
	}

	if !fieldName.CanSet() {
		return fmt.Errorf("Cannot set %v field value", name)
	}

	fieldType := fieldName.Type()
	v := reflect.ValueOf(value)
	if fieldType != v.Type() {
		var typeValue interface{}
		switch fieldType.Kind() {
		case reflect.Int:
			typeValue, _ = strconv.Atoi(value)
		case reflect.Int64:
			typeValue, _ = strconv.ParseInt(value, 10, 64)
		case reflect.Float64:
			typeValue, _ = strconv.ParseFloat(value, 64)
		default:
			return fmt.Errorf("Type didn't match %v", name)
		}
		v = reflect.ValueOf(typeValue)
	}
	fieldName.Set(v)
	return nil
}
