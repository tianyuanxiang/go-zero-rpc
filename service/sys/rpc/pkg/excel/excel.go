// Package excel 提供Excel文件的读写工具函数。
// 基于 xuri/excelize/v2 库实现，支持Excel导出和导入功能。
package excel

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel 将数据导出为Excel文件，返回文件的字节缓冲。
//
// 生成的Excel文件包含一个名为"Sheet1"的工作表，
// 第一行为表头，从第二行开始写入数据。
//
// 参数：
//   - headers : 表头列表，如：[]string{"用户名", "昵称", "邮箱"}
//   - data    : 数据行列表，每行为interface{}切片，与headers顺序对应
//
// 返回：
//   - *bytes.Buffer : 包含Excel文件内容的字节缓冲，可直接写入HTTP响应
//   - error         : 生成失败时的错误信息
func ExportToExcel(headers []string, data [][]interface{}) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	// 写入表头行（加粗样式）
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9E1F2"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#BFBFBF", Style: 1},
			{Type: "right", Color: "#BFBFBF", Style: 1},
			{Type: "top", Color: "#BFBFBF", Style: 1},
			{Type: "bottom", Color: "#BFBFBF", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建表头样式失败: %w", err)
	}

	for colIdx, header := range headers {
		// excelize使用列字母索引，从A开始
		colName, err := excelize.ColumnNumberToName(colIdx + 1)
		if err != nil {
			return nil, fmt.Errorf("转换列号失败（列索引：%d）: %w", colIdx, err)
		}
		cell := colName + "1"

		if err = f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("写入表头单元格[%s]失败: %w", cell, err)
		}
		if err = f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return nil, fmt.Errorf("设置表头样式失败: %w", err)
		}
	}

	// 写入数据行
	for rowIdx, row := range data {
		rowNum := rowIdx + 2 // 数据从第2行开始
		for colIdx, cellValue := range row {
			colName, err := excelize.ColumnNumberToName(colIdx + 1)
			if err != nil {
				return nil, fmt.Errorf("转换列号失败: %w", err)
			}
			cell := fmt.Sprintf("%s%d", colName, rowNum)
			if err = f.SetCellValue(sheetName, cell, cellValue); err != nil {
				return nil, fmt.Errorf("写入数据单元格[%s]失败: %w", cell, err)
			}
		}
	}

	// 自动调整列宽（基于内容长度）
	if err = f.SetColWidth(sheetName, "A", "Z", 15); err != nil {
		// 列宽调整失败不影响主流程，仅记录
		_ = err
	}

	// 将文件写入内存缓冲
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("生成Excel文件缓冲失败: %w", err)
	}

	return buf, nil
}

// ImportFromExcel 从上传的Excel文件中读取数据。
//
// 读取第一个工作表的所有数据，第一行通常为表头。
// 返回的二维切片包含所有行的原始字符串值。
//
// 参数：
//   - file : multipart上传的文件句柄
//
// 返回：
//   - [][]string : 二维字符串切片，第一行为表头，其余为数据行
//   - error      : 解析失败时的错误信息
func ImportFromExcel(file multipart.File) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败（请确认文件格式为.xlsx）: %w", err)
	}
	defer f.Close()

	// 获取第一个工作表名称
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel文件中没有工作表")
	}
	sheetName := sheets[0]

	// 读取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表[%s]数据失败: %w", sheetName, err)
	}

	return rows, nil
}

// ReadExcelSheetNames 读取Excel文件中所有工作表名称。
//
// 参数：
//   - file : multipart上传的文件句柄
//
// 返回：
//   - []string : 工作表名称列表
//   - error    : 读取失败时的错误信息
func ReadExcelSheetNames(file multipart.File) ([]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	return f.GetSheetList(), nil
}
