package freexl

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	mododbc32                      = syscall.NewLazyDLL(GetDllName())
	procFreeXLVersion              = mododbc32.NewProc("freexl_version")
	procFreeXLOpen                 = mododbc32.NewProc("freexl_open")
	procFreeXLOpenInfo             = mododbc32.NewProc("freexl_open_info")
	procFreeXLClose                = mododbc32.NewProc("freexl_close")
	procFreeXLGetInfo              = mododbc32.NewProc("freexl_get_info")
	procFreeXLGetWorksheetName     = mododbc32.NewProc("freexl_get_worksheet_name")
	proFreeXLSelectActiveWorksheet = mododbc32.NewProc("freexl_select_active_worksheet")
	procFreeXLGetActiveWorksheet   = mododbc32.NewProc("freexl_get_active_worksheet")
	procFreeXLWorksheetDimensions  = mododbc32.NewProc("freexl_worksheet_dimensions")
	procFreeXLGetSSTString         = mododbc32.NewProc("freexl_get_SST_string")
	procFreeXLGetFATEntry          = mododbc32.NewProc("freexl_get_FAT_entry")
	procFreeXLGetCellValue         = mododbc32.NewProc("freexl_get_cell_value")
)

func GetDllName() string {
	if winArch := os.Getenv("PROCESSOR_ARCHITECTURE"); winArch == "x86" {
		return "freexl.dll"
	} else {
		return "freexl64.dll"
	}
}

func FreeXLVerison() string {
	r0, _, _ := procFreeXLVersion.Call()
	return AnsiToString(r0)
}

func FreeXLOpen(path string) (FreeXLHandle, error) {
	var handle FreeXLHandle = 0
	in := syscall.StringBytePtr(path)
	r0, _, _ := procFreeXLOpen.Call(uintptr(unsafe.Pointer(in)), uintptr(unsafe.Pointer(&handle)))

	result := int32(r0)
	if FREEXL_OK != result {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return handle, nil
}

func FreeXLOpenInfo(path string) (FreeXLHandle, error) {
	var handle FreeXLHandle
	in := syscall.StringBytePtr(path)
	r0, _, _ := procFreeXLOpenInfo.Call(uintptr(unsafe.Pointer(in)), uintptr(unsafe.Pointer(&handle)))

	result := int32(r0)
	if result != FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return handle, nil
}

func FreeXLClose(handle FreeXLHandle) error {
	r0, _, _ := syscall.Syscall(procFreeXLClose.Addr(),
		1,
		uintptr(handle),
		0,
		0)

	result := int32(r0)
	if FREEXL_OK != result {
		return fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return nil
}

func FreeXLGetInfo(handle FreeXLHandle, what uint16) (uint, error) {
	var info uint = FREEXL_UNKNOWN
	r0, _, _ := procFreeXLGetInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info)))

	result := int32(r0)
	if FREEXL_OK != result {
		return FREEXL_UNKNOWN, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return info, nil
}

func FreeXLGetWorksheetName(handle FreeXLHandle, sheetIndex uint16) (string, error) {
	var s uintptr = 0
	var index uint16 = sheetIndex
	r0, _, _ := procFreeXLGetWorksheetName.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&index)),
		uintptr(unsafe.Pointer(&s)))

	result := int32(r0)
	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return AnsiToString(s), nil
}

func FreeXLSelectActiveWorksheet(handle FreeXLHandle, sheetIndex uint16) error {
	r0, _, _ := proFreeXLSelectActiveWorksheet.Call(
		uintptr(handle),
		uintptr(sheetIndex))

	result := int32(r0)
	if result != FREEXL_OK {
		return fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return nil
}

func FreeXLGetActiveWorksheet(handle FreeXLHandle) (uint16, error) {
	var sheetIndex uintptr = 0
	r0, _, _ := procFreeXLGetActiveWorksheet.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&sheetIndex)))

	result := int32(r0)
	if result != FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return uint16(sheetIndex), nil
}

func FreeXLWorksheetDimensions(handle FreeXLHandle) (rows uint, columns uint16, err error) {
	var rowsPtr uintptr = 0
	var columnsPtr uintptr = 0
	r0, _, _ := procFreeXLWorksheetDimensions.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&rowsPtr)),
		uintptr(unsafe.Pointer(&columnsPtr)))

	result := int32(r0)
	if result != FREEXL_OK {
		return 0, 0, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	rows = *(*uint)(unsafe.Pointer(&rowsPtr))
	columns = *(*uint16)(unsafe.Pointer(&columnsPtr))

	return rows, columns, nil
}

func FreeXLGetSSTString(handle FreeXLHandle, stringIndex uint16) (string, error) {
	var s uintptr = 0
	r0, _, _ := procFreeXLGetSSTString.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&stringIndex)),
		uintptr(unsafe.Pointer(&s)))

	result := int32(r0)
	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(result))
	}
	return AnsiToString(s), nil
}

func FreeXLGetFATEntry(handle FreeXLHandle, sectorIndex uint) (uint, error) {
	var nextSectorIndexPtr uintptr = 0
	r0, _, _ := procFreeXLGetFATEntry.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&nextSectorIndexPtr)))

	result := int32(r0)
	if result != FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(result))
	}

	nextSectorIndex := *(*uint)(unsafe.Pointer(&nextSectorIndexPtr))

	return nextSectorIndex, nil
}

func FreeXLGetCellValue(handle FreeXLHandle, rows uint, columns uint16) (string, error) {
	var cellValues FreeXLCell
	r0, _, _ := procFreeXLGetCellValue.Call(
		uintptr(handle),
		uintptr(rows),
		uintptr(columns),
		uintptr(unsafe.Pointer(&cellValues)))

	result := int32(r0)
	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(result))
	}

	value := ""
	switch cellValues.types {
	case FREEXL_CELL_INT:
		value = fmt.Sprintf("%v", *((*int32)(unsafe.Pointer(&cellValues.value))))
		break
	case FREEXL_CELL_DOUBLE:
		value = fmt.Sprintf("%v", *((*float64)(unsafe.Pointer(&cellValues.value))))
		break
	case FREEXL_CELL_TEXT:
		fallthrough
	case FREEXL_CELL_SST_TEXT:
		fallthrough
	case FREEXL_CELL_DATE:
		fallthrough
	case FREEXL_CELL_DATETIME:
		fallthrough
	case FREEXL_CELL_TIME:
		value = fmt.Sprintf("%v", AnsiToString(cellValues.value))
		break
	case FREEXL_CELL_NULL:
		break
	default:
		return "", fmt.Errorf("Invalid data-type")
	}
	return value, nil
}
