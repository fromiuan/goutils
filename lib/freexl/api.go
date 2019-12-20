package freexl

import (
	"unsafe"
)

const (

	/* constants */
	/** query is not applicable, or information is not available */
	FREEXL_UNKNOWN = 0

	/* CFBF constants */
	/** CFBF file is version 3 */
	FREEXL_CFBF_VER_3 = 3
	/** CFBF file is version 4 */
	FREEXL_CFBF_VER_4 = 4

	/** CFBF file uses 512 byte sectors */
	FREEXL_CFBF_SECTOR_512 = 512
	/** CFBF file uses 4096 (4K) sectors */
	FREEXL_CFBF_SECTOR_4096 = 4096

	/* BIFF versions */
	/** BIFF file is version 2 */
	FREEXL_BIFF_VER_2 = 2
	/** BIFF file is version 3 */
	FREEXL_BIFF_VER_3 = 3
	/** BIFF file is version 4 */
	FREEXL_BIFF_VER_4 = 4
	/** BIFF file is version 5 */
	FREEXL_BIFF_VER_5 = 5
	/** BIFF file is version 9 */
	FREEXL_BIFF_VER_8 = 8

	/* BIFF MaxRecordSize */
	/** Maximum BIFF record size is 2080 bytes */
	FREEXL_BIFF_MAX_RECSZ_2080 = 2080
	/** Maximum BIFF record size is 8224 bytes */
	FREEXL_BIFF_MAX_RECSZ_8224 = 8224

	/* BIFF DateMode */
	/** BIFF date mode starts at 1 Jan 1900 */
	FREEXL_BIFF_DATEMODE_1900 = 1900
	/** BIFF date mode starts at 2 Jan 1904 */
	FREEXL_BIFF_DATEMODE_1904 = 1904

	/* BIFF Obsfuscated */
	/** BIFF file is password protected */
	FREEXL_BIFF_OBFUSCATED = 3003
	/** BIFF file is not password protected */
	FREEXL_BIFF_PLAIN = 3004

	/* BIFF CodePage */
	/** BIFF file uses plain ASCII encoding */
	FREEXL_BIFF_ASCII = 0x016F
	/** BIFF file uses CP437 (OEM US format) encoding */
	FREEXL_BIFF_CP437 = 0x01B5
	/** BIFF file uses CP720 (Arabic DOS format) encoding */
	FREEXL_BIFF_CP720 = 0x02D0
	/** BIFF file uses CP737 (Greek DOS format) encoding */
	FREEXL_BIFF_CP737 = 0x02E1
	/** BIFF file uses CP775 (Baltic DOS format) encoding */
	FREEXL_BIFF_CP775 = 0x0307
	/** BIFF file uses CP850 (Western Europe DOS format) encoding */
	FREEXL_BIFF_CP850 = 0x0352
	/** BIFF file uses CP852 (Central Europe DOS format) encoding */
	FREEXL_BIFF_CP852 = 0x0354
	/** BIFF file uses CP855 (OEM Cyrillic format) encoding */
	FREEXL_BIFF_CP855 = 0x0357
	/** BIFF file uses CP857 (Turkish DOS format) encoding */
	FREEXL_BIFF_CP857 = 0x0359
	/** BIFF file uses CP858 (OEM Multiligual Latin 1 format) encoding */
	FREEXL_BIFF_CP858 = 0x035A
	/** BIFF file uses CP860 (Portuguese DOS format) encoding */
	FREEXL_BIFF_CP860 = 0x035C
	/** BIFF file uses CP861 (Icelandic DOS format) encoding */
	FREEXL_BIFF_CP861 = 0x035D
	/** BIFF file uses CP862 (Hebrew DOS format) encoding */
	FREEXL_BIFF_CP862 = 0x035E
	/** BIFF file uses CP863 (French Canadian DOS format) encoding */
	FREEXL_BIFF_CP863 = 0x035F
	/** BIFF file uses CP864 (Arabic DOS format) encoding */
	FREEXL_BIFF_CP864 = 0x0360
	/** BIFF file uses CP865 (Nordic DOS format) encoding */
	FREEXL_BIFF_CP865 = 0x0361
	/** BIFF file uses CP866 (Cyrillic DOS format) encoding */
	FREEXL_BIFF_CP866 = 0x0362
	/** BIFF file uses CP869 (Modern Greek DOS format) encoding */
	FREEXL_BIFF_CP869 = 0x0365
	/** BIFF file uses CP874 (Thai Windows format) encoding */
	FREEXL_BIFF_CP874 = 0x036A
	/** BIFF file uses CP932 (Shift JIS format) encoding */
	FREEXL_BIFF_CP932 = 0x03A4
	/** BIFF file uses CP936 (Simplified Chinese GB2312 format) encoding */
	FREEXL_BIFF_CP936 = 0x03A8
	/** BIFF file uses CP949 (Korean) encoding */
	FREEXL_BIFF_CP949 = 0x03B5
	/** BIFF file uses CP950 (Traditional Chinese Big5 format) encoding */
	FREEXL_BIFF_CP950 = 0x03B6
	/** BIFF file uses Unicode (UTF-16LE format) encoding */
	FREEXL_BIFF_UTF16LE = 0x04B0
	/** BIFF file uses CP1250 (Central Europe Windows) encoding */
	FREEXL_BIFF_CP1250 = 0x04E2
	/** BIFF file uses CP1251 (Cyrillic Windows) encoding */
	FREEXL_BIFF_CP1251 = 0x04E3
	/** BIFF file uses CP1252 (Windows Latin 1) encoding */
	FREEXL_BIFF_CP1252 = 0x04E4
	/** BIFF file uses CP1252 (Windows Greek) encoding */
	FREEXL_BIFF_CP1253 = 0x04E5
	/** BIFF file uses CP1254 (Windows Turkish) encoding */
	FREEXL_BIFF_CP1254 = 0x04E6
	/** BIFF file uses CP1255 (Windows Hebrew) encoding */
	FREEXL_BIFF_CP1255 = 0x04E7
	/** BIFF file uses CP1256 (Windows Arabic) encoding */
	FREEXL_BIFF_CP1256 = 0x04E8
	/** BIFF file uses CP1257 (Windows Baltic) encoding */
	FREEXL_BIFF_CP1257 = 0x04E9
	/** BIFF file uses CP1258 (Windows Vietnamese) encoding */
	FREEXL_BIFF_CP1258 = 0x04EA
	/** BIFF file uses CP1361 (Korean Johab) encoding */
	FREEXL_BIFF_CP1361 = 0x0551
	/** BIFF file uses Mac Roman encoding */
	FREEXL_BIFF_MACROMAN = 0x2710

	/* CELL VALUE Types */
	/** Cell has no value (empty cell) */
	FREEXL_CELL_NULL = 101
	/** Cell contains an integer value */
	FREEXL_CELL_INT = 102
	/** Cell contains a floating point number */
	FREEXL_CELL_DOUBLE = 103
	/** Cell contains a text value */
	FREEXL_CELL_TEXT = 104
	/** Cell contains a reference to a Single String Table entry (BIFF8) */
	FREEXL_CELL_SST_TEXT = 105
	/** Cell contains a number intended to represent a date */
	FREEXL_CELL_DATE = 106
	/** Cell contains a number intended to represent a date and time */
	FREEXL_CELL_DATETIME = 107
	/** Cell contains a number intended to represent a time */
	FREEXL_CELL_TIME = 108

	/* INFO params */
	/** Information query for CFBF version */
	FREEXL_CFBF_VERSION = 32001
	/** Information query for CFBF sector size */
	FREEXL_CFBF_SECTOR_SIZE = 32002
	/** Information query for CFBF FAT entry count */
	FREEXL_CFBF_FAT_COUNT = 32003
	/** Information query for BIFF version */
	FREEXL_BIFF_VERSION = 32005
	/** Information query for BIFF maximum record size */
	FREEXL_BIFF_MAX_RECSIZE = 32006
	/** Information query for BIFF date mode */
	FREEXL_BIFF_DATEMODE = 32007
	/** Information query for BIFF password protection state */
	FREEXL_BIFF_PASSWORD = 32008
	/** Information query for BIFF character encoding */
	FREEXL_BIFF_CODEPAGE = 32009
	/** Information query for BIFF sheet count */
	FREEXL_BIFF_SHEET_COUNT = 32010
	/** Information query for BIFF Single String Table entry count (BIFF8) */
	FREEXL_BIFF_STRING_COUNT = 32011
	/** Information query for BIFF format count */
	FREEXL_BIFF_FORMAT_COUNT = 32012
	/** Information query for BIFF extended format count */
	FREEXL_BIFF_XF_COUNT = 32013

	/* Error codes */
	FREEXL_OK             = 0  /**< No error, success */
	FREEXL_FILE_NOT_FOUND = -1 /**< .xls file does not exist or is
	not accessible for reading */
	FREEXL_NULL_HANDLE         = -2 /**< Null xls_handle argument */
	FREEXL_INVALID_HANDLE      = -3 /**< Invalid xls_handle argument */
	FREEXL_INSUFFICIENT_MEMORY = -4 /**< some kind of memory allocation
	  failure */
	FREEXL_NULL_ARGUMENT       = -5 /**< an unexpected null argument */
	FREEXL_INVALID_INFO_ARG    = -6 /**< invalid "what" parameter */
	FREEXL_INVALID_CFBF_HEADER = -7 /**< the .xls file does not contain a
	  valid CFBF header */
	FREEXL_CFBF_READ_ERROR = -8 /**< Read error. Usually indicates a
	  corrupt or invalid .xls file */
	FREEXL_CFBF_SEEK_ERROR = -9 /**< Seek error. Usually indicates a
	  corrupt or invalid .xls file */
	FREEXL_CFBF_INVALID_SIGNATURE = -10 /**< The .xls file does contain a
	  CFBF header, but the header is
	  broken or corrupted in some way
	*/
	FREEXL_CFBF_INVALID_SECTOR_SIZE = -11 /**< The .xls file does contain a
	  CFBF header, but the header is
	  broken or corrupted in some way
	*/
	FREEXL_CFBF_EMPTY_FAT_CHAIN = -12 /**< The .xls file does contain a
	  CFBF header, but the header is
	  broken or corrupted in some way
	*/
	FREEXL_CFBF_ILLEGAL_FAT_ENTRY = -13 /**< The file contains an invalid
	  File Allocation Table record */
	FREEXL_BIFF_INVALID_BOF = -14 /**< The file contains an invalid
	  BIFF format entry */
	FREEXL_BIFF_INVALID_SST = -15 /**< The file contains an invalid
	  Single String Table */
	FREEXL_BIFF_ILLEGAL_SST_INDEX = -16 /**< The requested Single String
	  Table entry is not available */
	FREEXL_BIFF_WORKBOOK_NOT_FOUND = -17 /**< BIFF does not contain a valid
	  workbook */
	FREEXL_BIFF_ILLEGAL_SHEET_INDEX = -18 /**< The requested worksheet is not
	  available in the workbook */
	FREEXL_BIFF_UNSELECTED_SHEET = -19 /**< There is no currently active
	  worksheet. Possibly a forgotten
	  call to
	  freexl_select_active_worksheet()
	*/
	FREEXL_INVALID_CHARACTER = -20 /**< Charset conversion detected an
	  illegal character (not within
	  the declared charset) */
	FREEXL_UNSUPPORTED_CHARSET = -21 /**< The requested charset
	  conversion is not available. */
	FREEXL_ILLEGAL_CELL_ROW_COL = -22 /**< The requested cell is outside
	  the valid range for the sheet*/
	FREEXL_ILLEGAL_RK_VALUE = -23 /**< Conversion of the RK value
	  failed. Possibly a corrupt file
	  or a bug in FreeXL. */
	FREEXL_ILLEGAL_MULRK_VALUE = -23 /**< Conversion of the MULRK value
	  failed. Possibly a corrupt file
	  or a bug in FreeXL. */
	FREEXL_INVALID_MINI_STREAM = -24 /**< The MiniFAT stream is invalid.
	  Possibly a corrupt file. */
	FREEXL_CFBF_ILLEGAL_MINI_FAT_ENTRY = -25 /**< The MiniFAT stream
	  contains an invalid entry.
	  Possibly a corrupt file. */
)

type FreeXLHandle uintptr
type FreeXLCell struct {
	types byte
	value uintptr
}

func AnsiToString(ptr uintptr) string {
	buf := make([]byte, 0, 0)
	for i := 0; ; i++ {
		ptr := unsafe.Pointer(ptr + uintptr(i))
		b := *((*byte)(ptr))
		if b == 0x00 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}
