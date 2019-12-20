package tools

import (
	"strings"
)

type BrowserClient struct {
	UserAngent string
}

func NewBrowserClient(userAngent string) *BrowserClient {
	return &BrowserClient{
		UserAngent: userAngent,
	}
}

func (b *BrowserClient) DetectOS(userangent string) int32 {
	Win2000 := strings.Contains(b.UserAngent, "Windows NT 5.0")
	if Win2000 {
		return 1
	}

	WinXP := strings.Contains(b.UserAngent, "Windows NT 5.1")
	if WinXP {
		return 2
	}

	Win2003 := strings.Contains(b.UserAngent, "Windows NT 5.2")
	if Win2003 {
		return 3
	}

	return 0
}

func (b *BrowserClient) DetectExplore(userangent string) int32 {
	IE6 := strings.Contains(b.UserAngent, "MSIE 6.0")
	if IE6 {
		return 1
	}

	IE7 := strings.Contains(b.UserAngent, "MSIE 7.0")
	if IE7 {
		return 2
	}

	IE8 := strings.Contains(b.UserAngent, "MSIE 8.0")
	if IE8 {
		return 3
	}

	IE9 := strings.Contains(b.UserAngent, "MSIE 9.0")
	if IE9 {
		return 4
	}

	IE10 := strings.Contains(b.UserAngent, "MSIE 10.0")
	if IE10 {
		return 5
	}

	IE11 := strings.Contains(b.UserAngent, "Trident")
	IE11S := strings.Contains(b.UserAngent, "11.0")
	if IE11 && IE11S {
		return 6
	}
	return 0
}
