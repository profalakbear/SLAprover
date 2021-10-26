package server

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ExcelRow represents each row of sheet in excel file ...
type ExcelRow struct {
	// There 19 columns in total in excel file but we use 14 fields here ...
	ReceiptNumber         uint64 `json:"receipt_number"`
	ReceiptState          string `json:"receipt_state"`
	OperationNumber       uint64 `json:"operation_number"`
	ServiceId             string `json:"service_id"` // uint64
	PointName             string `json:"point_name"`
	PointId               uint64 `json:"point_id"`
	SumIn                 uint64 `json:"sum_in"`
	SumOut                uint64 `json:"sum_out"`
	Fee                   uint64 `json:"fee"`
	Cash                  uint64 `json:"cash"`
	ProviderTransactionId uint64 `json:"provider_transaction_id"`
	PaymentType           string `json:"payment_type"`
	CorrectedOperation    string `json:"corrected_operation"`
	CorrectingOperation   string `json:"correcting_operation"`
}

// This function takes a string array traverse it and map each index to its appropriate field in ExcelRow object ...
func Mapper(arr []string) *ExcelRow {
	receipt_number, _ := strconv.ParseUint(GetElementFromArray(arr, 5), 10, 64)
	receipt_state := GetElementFromArray(arr, 6)
	operation_number, _ := strconv.ParseUint(GetElementFromArray(arr, 7), 10, 64)
	service_id := GetElementFromArray(arr, 8)
	point_name := GetElementFromArray(arr, 9)
	point_id, _ := strconv.ParseUint(GetElementFromArray(arr, 10), 10, 64)
	sum_in, _ := strconv.ParseUint(GetElementFromArray(arr, 11), 10, 64)
	sum_out, _ := strconv.ParseUint(GetElementFromArray(arr, 12), 10, 64)
	fee, _ := strconv.ParseUint(GetElementFromArray(arr, 13), 10, 64)
	cash, _ := strconv.ParseUint(GetElementFromArray(arr, 14), 10, 64)
	provider_transaction_id, _ := strconv.ParseUint(GetElementFromArray(arr, 15), 10, 64)
	payment_type := GetElementFromArray(arr, 16)
	corrected_operation := GetElementFromArray(arr, 17)
	correcting_operation := GetElementFromArray(arr, 18)
	// Mapping ...
	excelRow := &ExcelRow{
		ReceiptNumber:         receipt_number,
		ReceiptState:          receipt_state,
		OperationNumber:       operation_number,
		ServiceId:             service_id,
		PointName:             point_name,
		PointId:               point_id,
		SumIn:                 sum_in,
		SumOut:                sum_out,
		Fee:                   fee,
		Cash:                  cash,
		ProviderTransactionId: provider_transaction_id,
		PaymentType:           payment_type,
		CorrectedOperation:    corrected_operation,
		CorrectingOperation:   correcting_operation,
	}
	return excelRow
}

// This function will get an element from an array with the given index
// and will return if exist, otherwise "nil" ...
// Why did I need this function? Answer is below ....
// Because sometimes, when we read from excel some columns in excel file can be empty. In this case
// that columns will not occur in our array ...
// In such situation we get ArrayIndexOutOfBoundError ...
func GetElementFromArray(arr []string, index int) string {
	// We check here if given index is in bound of array index or not ...
	// If it exists then we return element itself otherwise we return "nil" not to raise error ...
	// Consider example below ...
	// arr[1,2,3,4,5] ==> (len(arr)-1) ==> 4
	// index = 6
	// if 6 > 4 then return "nil"
	if index > (len(arr)-1) || arr[index] == "" {
		return "00000000"
	}
	// Consider example below ...
	// arr[1,2,3,4,5] ==> (len(arr)-1) ==> 4
	// index = 3
	// if 3 > 4 then return "nil" ==> 3 is not greater than 4 so this line is false
	// that is why we return element itself ...
	return arr[index]
}

func GetHeader(excel *excelize.File, sheetName string) ([]string, int, error) {
	var err error
	var rows *excelize.Rows
	rows, err = excel.Rows(sheetName)
	var row_elements []string

	if err != nil {
		return nil, 0, err
	}
	// "nextIterIsTrue" variable gives us boolean value if there is a row after x ...
	nextIterIsTrue := rows.Next()
	if nextIterIsTrue {
		row_elements, err = rows.Columns()
	}
	if err != nil {
		fmt.Println(err)
	}

	return row_elements, len(row_elements), nil
}

func ReadExcelAndConvertArrayOfJson(reader io.Reader) ([]ExcelRow, error) {
	var excelList []ExcelRow
	var statement bool
	var count int = 0
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	sheetNames := GetSheetNameList(f)

	for _, elem := range sheetNames {
		_, length, err := GetHeader(f, elem)
		if err != nil || length < 19 || length > 19 {
			fmt.Println(err)
			return nil, err
		}
	}
	// Looping through all sheet names in an excel sheet ...
	for _, elem := range sheetNames {
		// Getting rows in "elem" sheet ...
		rows, err := f.Rows(elem)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// Getting 2D array of each excel sheet
		// and it's row count(total number of rows in one excel sheet) ==> len(rowNumberInSheet) ...
		totalRowNumberInSheet, err := f.GetRows(elem)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for rows.Next() {
			// Loop Count ==> If loop count greater than (totalRowNumberInSheet - 3) in a sheet then we break ...
			// Consider example below ...
			// Let's say we have 265 rows in total (firstrow(header) and lastrow(Итог) included) ...
			// We are omitting first row in our code because we handle first line in our GetHeader function ...
			// That is why we get 264 rows in total after omitting first line ...
			// And finally we extract firstrow(header) and lastrow(Итог) count from total ...
			// In total we are extracting 3 from 265 that is why minus(-) 3 ...
			// (totalRowNumberInSheet - 3) equals to EOF(end of file) ...
			if count > (len(totalRowNumberInSheet) - 3) {
				break
			}
			row, _ := rows.Columns()
			// Omitting first line here ...
			if !statement {
				statement = true
				continue
			}
			// Mapper function maps our excel row to the ExcelRow struct ...
			xls := Mapper(row)
			// Takes account of "Loop Count"
			count = count + 1
			// Appendig element to our ExcelList ...
			excelList = append(excelList, *xls)
		}
	}
	return excelList, nil
}

func SplitIntoChunks(requestNumber int, list []ExcelRow) []ExcelRow {
	// Total data in list
	total := len(list)
	// Data number per request.
	// We want 100 data per reqeust...
	perRequest := CHUNK_LIMIT
	// Pay attention to the variable(requestNumber) in the function parameters ...
	// requestNumber variable holds "Which request are we sending?" , "Which number of request?" ...
	// In each request, "requestNumber" parameter is changing ...
	// In first request: requestNumber=1 (also staring value)
	// In first request: requestNumber=2 and so on ...
	var dataToReturn []ExcelRow
	// Consider example below ...
	// Assume that we have variables like below ...
	// total = 423; requestNumber = 1; perRequest = 100
	// We will do 5 request in total because we have 423 data. That is how we get this: (Attention below)
	// Ceiling(total/perRequest) ==> Ceiling(423/100) = 5
	// if 423 <= (1*100) && 423 (1-1)*100 then
	if (total <= requestNumber*perRequest) && total >= ((requestNumber-1)*perRequest) {
		// do this ==> list[(1-1)*100 : 423]
		dataToReturn = list[(requestNumber-1)*perRequest : total]
		return dataToReturn
		// If 423 >= 1*100 then
	} else if total >= requestNumber*perRequest {
		// do this ==> list[(1-1)*100 : 1*100] ==> list[0 : 100]...
		// In first request(requestNumber=1), We are slicing from 0 to 100(excluding) ...
		// In first request(requestNumber=2), We are slicing from 100 to 200(excluding)  ==> list[100 : 200] ...
		// In first request(requestNumber=3), We are slicing from 200 to 300(excluding)  ==> list[200 : 300] ...
		// In first request(requestNumber=4), We are slicing from 300 to 400(excluding)  ==> list[300 : 400] ...
		// In first request(requestNumber=5), We are slicing from 400 to 423(to the end) ==> list[400 : ] ...
		dataToReturn = list[(requestNumber-1)*perRequest : requestNumber*perRequest]
		return dataToReturn
	}
	// If request does not meet conditions above we return empty array ....
	dataToReturn = []ExcelRow{}
	return dataToReturn
}

func GetSheetNameList(excel *excelize.File) []string {
	sheetNames := excel.GetSheetList()
	return sheetNames
}

func RequestToExternalApi(list []ExcelRow) bool {
	// Send rquest to the external api here ...
	return true
}

type Error struct {
	Code    string `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewCustomBadRequestError(message string) *Error {
	return &Error{
		Code:    "400",
		Type:    "Bad Request",
		Message: message,
	}
}

func NewCustomInternalServerError(message string) *Error {
	return &Error{
		Code:    "500",
		Type:    "Server Error",
		Message: message,
	}
}
