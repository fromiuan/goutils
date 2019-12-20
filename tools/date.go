package tools

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	YYYY_MM_DD          = "2006-01-02"
	YYYYMMDD            = "20060102"
	YYYYMM              = "200601"
	YYYY                = "2006"
	YYYY_MM_DD_HH_MI_SS = "2006-01-02 15:04:05"
	YYYYMMDDTHHMISSZ    = "20060102T150405Z"
	YYYYMMDDHHMISS      = "20060102150405"
)

// 获取指定日期开始,结束时间戳,day 格式：例如:2019-01-01
func GetDayStartAndEndUnix(day string) (int64, int64) {
	dayStart := day + " 00:00:00"
	dayEnd := day + " 23:59:59"
	parseStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dayStart, time.Local)
	parseEndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dayEnd, time.Local)
	return parseStartTime.Unix(), parseEndTime.Unix()
}

// 获取指定日期开始,结束时间戳
func GetDayStartAndEndUnixByTime(date time.Time) (int64, int64) {
	year, month, day := date.Date()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return d.Unix(), d.Unix() + 24*60*60 - 1
}

// 时间戳按格式转时间字符串
func TimestampFormat(timestamp int64, format string) string {
	if len(strings.TrimSpace(format)) == 0 {
		return DefaultString(time.Unix(timestamp, -1))
	}
	return time.Unix(timestamp, 0).Format(format)
}

// 时间按格式转字符串
func TimeFormat(t time.Time, format string) string {
	if len(strings.TrimSpace(format)) == 0 {
		return DefaultString(t)
	}
	return t.Format(format)
}

// 时间转字符串
func DefaultString(t time.Time) string {
	return t.Format(YYYY_MM_DD_HH_MI_SS)
}

func StringToTimestamp(str string, format string) int64 {
	date, err := time.ParseInLocation(format, str, time.Local)
	if err != nil {
		return 0
	}
	return date.Unix()
}

// 时间文字
func TimeString(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour < 6:
		return "凌晨"
	case hour < 9:
		return "早上"
	case hour < 12:
		return "上午"
	case hour < 14:
		return "中午"
	case hour < 17:
		return "下午"
	case hour < 19:
		return "傍晚"
	case hour < 22:
		return "晚上"
	default:
		return "夜里"
	}
}

// 相对时间文字，例如：xx分钟前
func RelativeTimeString(stamp int) string {
	diff := time.Now().Unix() - int64(stamp)
	if diff == 0 {
		return "刚刚"
	}

	unit := []string{"年", "天", "小时", "分钟", "秒钟"}
	byTime := []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}

	sliceLen := len(byTime)
	diffAbs := int64(math.Abs(float64(diff)))
	buf := bytes.Buffer{}

	for i := 0; i < sliceLen; i++ {
		if diffAbs >= byTime[i] {
			temp := math.Floor(float64(diffAbs / byTime[i]))
			buf.WriteString(strconv.FormatFloat(temp, 'f', -1, 64))
			buf.WriteString(unit[i])

			if i == 1 { // XX天XX时
				temp2 := math.Floor(float64((diffAbs - int64(temp)*byTime[i]) / byTime[i+1]))
				buf.WriteString(strconv.FormatFloat(temp2, 'f', -1, 64))
				buf.WriteString(unit[i+1])
			}

			if i == 2 { // XX时XX分
				temp2 := math.Floor(float64((diffAbs - int64(temp)*byTime[i]) / byTime[i+1]))
				buf.WriteString(strconv.FormatFloat(temp2, 'f', -1, 64))
				buf.WriteString(unit[i+1])
			}

			break
		}
	}

	switch {
	case diff > 0:
		buf.WriteString("前")
		return buf.String()
	case diff < 0:
		buf.WriteString("后")
		return buf.String()
	default:
		return ""
	}
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取传入的时间所在年份的第一天，即某年元旦0点
func GetFirstDateOfYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
}

// 获取传入的开始时间[格式200601]到结束时间[格式200601]为止所有的月份
func GetMonths(startMonth string, endMonth string) []string {
	months := make([]string, 0)
	y1 := time.Unix(StringToTimestamp(startMonth, YYYYMM), 0).Year()
	m1 := time.Unix(StringToTimestamp(startMonth, YYYYMM), 0).Month()
	y2 := time.Unix(StringToTimestamp(endMonth, YYYYMM), 0).Year()
	m2 := time.Unix(StringToTimestamp(endMonth, YYYYMM), 0).Month()

	if y2-y1 == 0 {
		for i := m1; i <= m2; i++ {
			months = append(months, fmt.Sprintf("%d%02d", y2, i))
		}
	} else {
		for i := m1; i <= 12; i++ {
			months = append(months, fmt.Sprintf("%d%02d", y1, i))
		}
		for i := 1; i <= int(m2); i++ {
			months = append(months, fmt.Sprintf("%d%02d", y2, i))
		}
		for i := 1; i < y2-y1; i++ {
			for j := 1; j <= 12; j++ {
				months = append(months, fmt.Sprintf("%d%02d", y1+i, j))
			}
		}
	}
	return months
}

type hourRange struct {
	StartTime int64
	EndTime   int64
}

// 获取传入的时间所在天的24小时
func GetHourOfDay(d time.Time) []hourRange {
	var allHour []hourRange
	zeroTime := GetZeroTime(d)
	theTime := zeroTime.Unix()
	for i := 0; i < 24; i++ {
		allHour = append(allHour, hourRange{
			StartTime: theTime,
			EndTime:   theTime + 3599,
		})
		theTime = theTime + 3600
	}
	return allHour
}

// 两个时间戳之间的所有日期
func GetDaysOfTheSection(startTime, endTime int64) ([]int64, []string) {
	var (
		dayListInt64  []int64
		dayListString []string
	)

	start := time.Unix(startTime, 0)
	dayListInt64 = append(dayListInt64, ConvertInt64(start.Format(YYYYMMDD)))
	dayListString = append(dayListString, start.Format(YYYY_MM_DD))

	for {
		nextDate := start.AddDate(0, 0, 1)

		if nextDate.Unix() > endTime {
			break
		} else {
			dayListInt64 = append(dayListInt64, ConvertInt64(nextDate.Format(YYYYMMDD)))
			dayListString = append(dayListString, nextDate.Format(YYYY_MM_DD))
			start = nextDate
		}
	}

	return dayListInt64, dayListString
}
