package testcase

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
)

// 文档地址
// https://medium.com/@akaivdo/golang-how-to-read-and-write-excel-files-fb2120b63f86

func TestReadExcel(t *testing.T) {
	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get value from cell by given name and axis
	cellValue, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cellValue)

	// Get all the rows in the sheet
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Iterate over the rows and print the cell values
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func TestUpdateSingleCell(t *testing.T) {
	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set value of cell D4 to 88
	err = f.SetCellValue("Sheet1", "D4", 88)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save the changes to the file
	err = f.Save()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cell D4 updated successfully.")
}

func TestUpdateMultiCells(t *testing.T) {
	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set values of cells B3, C3, D3, to "Jack", "Physics", 90
	data := []any{"Jack", "Physics", 90}
	err = f.SetSheetRow("Sheet1", "B3", &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save the changes to the file.
	err = f.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Cells B3, C3, D3 updated successfully.")
}
