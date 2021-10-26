package server

import (
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	CHUNK_LIMIT int = 100
)

type SLAproverHandler struct {
}

func NewSLAproverHandler() *SLAproverHandler {
	return &SLAproverHandler{}
}

func (h *SLAproverHandler) checkRegistry(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewCustomBadRequestError(err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewCustomInternalServerError(err.Error()))
	}
	// Ensure to close stream ...
	defer src.Close()
	// Getting list of all objects(excel rows converting into ExcelRow struct) from file ...
	list, err := ReadExcelAndConvertArrayOfJson(src)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewCustomInternalServerError(err.Error()))
	}

	arraySize := len(list)
	if arraySize < CHUNK_LIMIT {
		// As if we are sending request to external API
		fmt.Println("\nSending Request To External API \nArraySize lower than or equal to 100 \n", list)
	} else if arraySize > CHUNK_LIMIT {
		// Splitting into chunks here ...
		var times int
		var x float64
		x = float64(arraySize) / float64(100)
		x = math.Ceil(x)
		times = int(x)

		for i := 0; i < times; i++ {
			y := SplitIntoChunks(i+1, list)
			// As if we are sending request to external API
			fmt.Printf("\nSending Request To External API \nArraySize greater than 100 that is why we split it into chunks\nChunk NO %d \n", i+1)
			fmt.Println("Size of chunk", len(y))
			fmt.Println(y)
		}
	}
	return ctx.JSON(200, echo.Map{
		// Return real api response here ...
		"Message": "Check-Registry handler works fine!",
	})
}
