// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux
// +build cgo

package freexl

/*
#cgo pkg-config: freexl
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <freexl.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// func FreeXLVerison() string {
// 	version := C.freexl_version()
// 	return C.GoString(version)
// }

func FreeXLOpen(path string) (FreeXLHandle, error) {
	var ptr unsafe.Pointer
	result := C.freexl_open(C.CString(path), &ptr)
	if FREEXL_OK != result {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return FreeXLHandle(ptr), nil
}

func FreeXLOpenInfo(path string) (FreeXLHandle, error) {
	var ptr unsafe.Pointer
	result := C.freexl_open_info(C.CString(path), &ptr)
	if FREEXL_OK != result {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return FreeXLHandle(ptr), nil
}

func FreeXLClose(handle FreeXLHandle) error {
	result := C.freexl_close(unsafe.Pointer(handle))
	if FREEXL_OK != result {
		return fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return nil
}

func FreeXLGetInfo(handle FreeXLHandle, what uint16) (uint, error) {
	var info C.uint = FREEXL_UNKNOWN
	cptr := unsafe.Pointer(handle)
	result := C.freexl_get_info(cptr, C.ushort(what), &info)
	if result != C.FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return uint(info), nil
}

func FreeXLGetWorksheetName(handle FreeXLHandle, sheetIndex uint16) (string, error) {
	cptr := unsafe.Pointer(handle)
	var buf *C.char
	result := C.freexl_get_worksheet_name(cptr, C.ushort(sheetIndex), &buf)
	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return C.GoString(buf), nil
}

func FreeXLSelectActiveWorksheet(handle FreeXLHandle, sheetIndex uint16) error {
	cptr := unsafe.Pointer(handle)
	result := C.freexl_select_active_worksheet(cptr, C.ushort(sheetIndex))
	if result != FREEXL_OK {
		return fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return nil
}

func FreeXLGetActiveWorksheet(handle FreeXLHandle) (uint16, error) {
	cptr := unsafe.Pointer(handle)
	var sheetIndex C.ushort

	result := C.freexl_get_active_worksheet(
		cptr,
		&sheetIndex)
	if result != FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}

	return uint16(sheetIndex), nil
}

func FreeXLWorksheetDimensions(handle FreeXLHandle) (rows uint, columns uint16, err error) {
	cptr := unsafe.Pointer(handle)
	var crows C.uint = 0
	var colums C.ushort = 0

	result := C.freexl_worksheet_dimensions(cptr, &crows, &colums)
	if result != FREEXL_OK {
		return 0, 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return uint(crows), uint16(colums), nil
}

func FreeXLGetSSTString(handle FreeXLHandle, stringIndex uint16) (string, error) {
	cptr := unsafe.Pointer(handle)
	var buf *C.char
	result := C.freexl_get_SST_string(cptr, C.ushort(stringIndex), &buf)
	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}
	return C.GoString(buf), nil
}

func FreeXLGetFATEntry(handle FreeXLHandle, sectorIndex uint) (uint, error) {
	cptr := unsafe.Pointer(handle)
	var nextSectorIndexPtr C.uint = 0

	result := C.freexl_get_FAT_entry(cptr, C.uint(sectorIndex), &nextSectorIndexPtr)

	if result != FREEXL_OK {
		return 0, fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
	}

	return uint(nextSectorIndexPtr), nil
}

func FreeXLGetCellValue(handle FreeXLHandle, rows uint, columns uint16) (string, error) {
	var cellValues FreeXLCell
	cptr := unsafe.Pointer(handle)

	result := C.freexl_get_cell_value(
		cptr,
		C.uint(rows),
		C.ushort(columns),
		(*C.struct_FreeXL_CellValue_str)(unsafe.Pointer(&cellValues)))

	if result != FREEXL_OK {
		return "", fmt.Errorf("[%d]%s", result, getErrStr(int32(result)))
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
