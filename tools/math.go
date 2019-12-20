package tools

import (
	"math"
	"strconv"
	"strings"
)

// 普通四舍五入
func MathRound(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Round(f*n10) / n10
}

// 普通不四舍五入
func MathTruncFloat(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc(f*n10) / n10
}

//普通四舍五入
func MathBeanRound(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Round(f*n10) / n10
}

//银行家四舍五入
func MathRoundToEven(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.RoundToEven(f*n10) / n10
}

// 保留2位小数 不进行四舍五入
func MathCutFloat(f float64) (float64, error) {
	s := strconv.FormatFloat(f, 'f', 5, 64)
	sArr := strings.Split(s, ".")
	if len(sArr) == 2 {
		if len(sArr[1]) > 2 {
			sArr[1] = sArr[1][:2]
		}
	}
	news := strings.Join(sArr, ".")
	newf, err := strconv.ParseFloat(news, 64)
	if err != nil {
		return f, err
	}
	return newf, nil
}

func MustInt(args interface{}) int {
	if s, ok := args.(int); ok {
		return s
	} else {
		return 0
	}
}

func MustString(args interface{}) string {
	if s, ok := args.(string); ok {
		return s
	} else {
		return ""
	}

}

func MustInt64(args interface{}) int64 {
	if s, ok := args.(int64); ok {
		return s
	} else {
		return 0
	}
}

func MustFloat64(args interface{}) float64 {
	if s, ok := args.(float64); ok {
		return s
	} else {
		return 0
	}
}
