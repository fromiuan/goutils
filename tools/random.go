package tools

import (
	_ "errors"
	"fmt"
	_ "fmt"
	"math"
	"math/rand"
	_ "strconv"
	_ "strings"
	"time"
)

var (
	defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandomSpec0(count uint, start, end int, letters, numbers bool,
	chars []rune, rand *rand.Rand) string {
	if count == 0 {
		return ""
	}
	if start == 0 && end == 0 {
		end = 'z' + 1
		start = ' '
		if !letters && !numbers {
			start = 0
			end = math.MaxInt32
		}
	}
	buffer := make([]rune, count)
	gap := end - start
	for count != 0 {
		count--
		var ch rune
		if len(chars) == 0 {
			ch = rune(rand.Intn(gap) + start)
		} else {
			ch = chars[rand.Intn(gap)+start]
		}
		if letters && ((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) ||
			numbers && (ch >= '0' && ch <= '9') ||
			(!letters && !numbers) {
			if ch >= rune(56320) && ch <= rune(57343) {
				if count == 0 {
					count++
				} else {
					buffer[count] = ch
					count--
					buffer[count] = rune(55296 + rand.Intn(128))
				}
			} else if ch >= rune(55296) && ch <= rune(56191) {
				if count == 0 {
					count++
				} else {
					// high surrogate, insert low surrogate before putting it in
					buffer[count] = rune(56320 + rand.Intn(128))
					count--
					buffer[count] = ch
				}
			} else if ch >= rune(56192) && ch <= rune(56319) {
				// private high surrogate, no effing clue, so skip it
				count++
			} else {
				buffer[count] = ch
			}
		} else {
			count++
		}
	}
	return string(buffer)
}

func RandomSpec1(count uint, start, end int, letters, numbers bool) string {
	return RandomSpec0(count, start, end, letters, numbers, nil, defaultRand)
}

func RandomAlphaOrNumeric(count uint, letters, numbers bool) string {
	return RandomSpec1(count, 0, 0, letters, numbers)
}

func RandomString(count uint) string {
	return RandomAlphaOrNumeric(count, false, false)
}

func RandomStringSpec0(count uint, set []rune) string {
	return RandomSpec0(count, 0, len(set)-1, false, false, set, defaultRand)
}

func RandomStringSpec1(count uint, set string) string {
	return RandomStringSpec0(count, []rune(set))
}

func RandomAscii(count uint) string {
	return RandomSpec1(count, 32, 127, false, false)
}

func RandomAlphabetic(count uint) string {
	return RandomAlphaOrNumeric(count, true, false)
}

func RandomAlphanumeric(count uint) string {
	return RandomAlphaOrNumeric(count, true, true)
}

func RandomNumeric(count uint) string {
	return RandomAlphaOrNumeric(count, false, true)
}

func RandomMac() string {
	return fmt.Sprintf(`%02x:%02x:%02x:%02x:%02x:%02x`,
		rand.Intn(0xff),
		rand.Intn(0xff),
		rand.Intn(0xff),
		rand.Intn(0xff),
		rand.Intn(0xff),
		rand.Intn(0xff))
}
