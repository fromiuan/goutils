package freexl

var errno = map[int]string{
	/* Error codes */
	0: "FREEXL_OK",
	/**< No error success */
	-1: "FREEXL_FILE_NOT_FOUND", /**< .xls file does not exist or is
	not accessible for reading */
	-2: "FREEXL_NULL_HANDLE",         /**< Null xls_handle argument */
	-3: "FREEXL_INVALID_HANDLE",      /**< Invalid xls_handle argument */
	-4: "FREEXL_INSUFFICIENT_MEMORY", /**< some kind of memory allocation
	failure */
	-5: "FREEXL_NULL_ARGUMENT",       /**< an unexpected null argument */
	-6: "FREEXL_INVALID_INFO_ARG",    /**< invalid "what" parameter */
	-7: "FREEXL_INVALID_CFBF_HEADER", /**< the .xls file does not contain a
	valid CFBF header */
	-8: "FREEXL_CFBF_READ_ERROR", /**< Read error. Usually indicates a
	corrupt or invalid .xls file */
	-9: "FREEXL_CFBF_SEEK_ERROR", /**< Seek error. Usually indicates a
	corrupt or invalid .xls file */
	-10: "FREEXL_CFBF_INVALID_SIGNATURE", /**< The .xls file does contain a
	CFBF header: but the header is
	broken or corrupted in some way
	*/
	-11: "FREEXL_CFBF_INVALID_SECTOR_SIZE", /**< The .xls file does contain a
	CFBF header: but the header is
	broken or corrupted in some way
	*/
	-12: "FREEXL_CFBF_EMPTY_FAT_CHAIN", /**< The .xls file does contain a
	CFBF header: but the header is
	broken or corrupted in some way
	*/
	-13: "FREEXL_CFBF_ILLEGAL_FAT_ENTRY", /**< The file contains an invalid
	File Allocation Table record */
	-14: "FREEXL_BIFF_INVALID_BOF", /**< The file contains an invalid
	BIFF format entry */
	-15: "FREEXL_BIFF_INVALID_SST", /**< The file contains an invalid
	Single String Table */
	-16: "FREEXL_BIFF_ILLEGAL_SST_INDEX", /**< The requested Single String
	Table entry is not available */
	-17: "FREEXL_BIFF_WORKBOOK_NOT_FOUND", /**< BIFF does not contain a valid
	workbook */
	-18: "FREEXL_BIFF_ILLEGAL_SHEET_INDEX", /**< The requested worksheet is not
	available in the workbook */
	-19: "FREEXL_BIFF_UNSELECTED_SHEET", /**< There is no currently active
	worksheet. Possibly a forgotten
	call to
	freexl_select_active_worksheet()
	*/
	-20: "FREEXL_INVALID_CHARACTER", /**< Charset conversion detected an
	illegal character (not within
	the declared charset) */
	-21: "FREEXL_UNSUPPORTED_CHARSET", /**< The requested charset
	conversion is not available. */
	-22: "FREEXL_ILLEGAL_CELL_ROW_COL", /**< The requested cell is outside
	the valid range for the sheet*/
	-23: "FREEXL_ILLEGAL_RK_VALUE", /**< Conversion of the RK value
	failed. Possibly a corrupt file
	or a bug in -0: "FreeXL. */
	/*-23: "FREEXL_ILLEGAL_MULRK_VALUE", *< Conversion of the MULRK value
	failed. Possibly a corrupt file
	or a bug in -0: "FreeXL. */
	-24: "FREEXL_INVALID_MINI_STREAM", /**< The MiniFAT stream is invalid.
	Possibly a corrupt file. */
	-25: "FREEXL_CFBF_ILLEGAL_MINI_FAT_ENTRY", /**< The MiniFAT stream
	contains an invalid entry.
	Possibly a corrupt file. */
}

var errstr = map[int]string{
	/* Error codes */
	FREEXL_OK:                       " No error, success ",
	FREEXL_FILE_NOT_FOUND:           " .xls file does not exist or is not accessible for reading ",
	FREEXL_NULL_HANDLE:              " Null xls_handle argument ",
	FREEXL_INVALID_HANDLE:           " Invalid xls_handle argument ",
	FREEXL_INSUFFICIENT_MEMORY:      " some kind of memory allocation failure ",
	FREEXL_NULL_ARGUMENT:            " an unexpected null argument ",
	FREEXL_INVALID_INFO_ARG:         " invalid what parameter ",
	FREEXL_INVALID_CFBF_HEADER:      " the .xls file does not contain a valid CFBF header ",
	FREEXL_CFBF_READ_ERROR:          " Read error. Usually indicates a corrupt or invalid .xls file ",
	FREEXL_CFBF_SEEK_ERROR:          " Seek error. Usually indicates a corrupt or invalid .xls file ",
	FREEXL_CFBF_INVALID_SIGNATURE:   " The .xls file does contain a CFBF header, but the header is broken or corrupted in some way",
	FREEXL_CFBF_INVALID_SECTOR_SIZE: " The .xls file does contain a CFBF header, but the header is broken or corrupted in some way",
	FREEXL_CFBF_EMPTY_FAT_CHAIN:     " The .xls file does contain a CFBF header, but the header is broken or corrupted in some way",
	FREEXL_CFBF_ILLEGAL_FAT_ENTRY:   " The file contains an invalid File Allocation Table record ",
	FREEXL_BIFF_INVALID_BOF:         " The file contains an invalid BIFF format entry ",
	FREEXL_BIFF_INVALID_SST:         " The file contains an invalid Single String Table ",
	FREEXL_BIFF_ILLEGAL_SST_INDEX:   " The requested Single String Table entry is not available ",
	FREEXL_BIFF_WORKBOOK_NOT_FOUND:  " BIFF does not contain a valid workbook ",
	FREEXL_BIFF_ILLEGAL_SHEET_INDEX: " The requested worksheet is not available in the workbook ",
	FREEXL_BIFF_UNSELECTED_SHEET:    " There is no currently active  worksheet. Possibly a forgotten call to  freexl_select_active_worksheet()",
	FREEXL_INVALID_CHARACTER:        " Charset conversion detected an illegal character (not within the declared charset) ",
	FREEXL_UNSUPPORTED_CHARSET:      " The requested charset conversion is not available. ",
	FREEXL_ILLEGAL_CELL_ROW_COL:     " The requested cell is outside the valid range for the sheet",
	FREEXL_ILLEGAL_RK_VALUE:         " Conversion of the RK value failed. Possibly a corrupt file or a bug in FreeXL. ",
	// FREEXL_ILLEGAL_MULRK_VALUE:         " Conversion of the MULRK value failed. Possibly a corrupt file or a bug in FreeXL. ",
	FREEXL_INVALID_MINI_STREAM:         " The MiniFAT stream is invalid. Possibly a corrupt file. ",
	FREEXL_CFBF_ILLEGAL_MINI_FAT_ENTRY: " The MiniFAT stream contains an invalid entry. Possibly a corrupt file. ",
}

func getErrStr(no int32) string {
	if s, ok := errstr[int(no)]; ok {
		return s
	}
	return "FREEXL_ERROR_UNKNOWN"
}

func getErrNo(no int32) string {
	if s, ok := errno[int(no)]; ok {
		return s
	}
	return "FREEXL_ERROR_UNKNOWN"
}
