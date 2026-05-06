// Package upload 提供文件上传处理工具函数。
// 支持multipart/form-data格式的文件接收和本地存储。
package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 文件上传限制常量
const (
	// MaxUploadSize 最大上传文件大小：50MB
	MaxUploadSize = 50 << 20
	// FormFieldName 表单中文件字段的名称
	FormFieldName = "file"
)

// FileInfo 文件上传结果信息结构体。
type FileInfo struct {
	// Filename 存储到磁盘的文件名（UUID生成）
	Filename string
	// OriginName 原始文件名（用户上传时的文件名）
	OriginName string
	// FilePath 文件在服务器上的完整存储路径
	FilePath string
	// FileURL 文件的访问URL路径（相对路径，如：/uploads/2024/01/xxx.xlsx）
	FileURL string
	// FileSize 文件大小（字节）
	FileSize int64
	// FileType 文件MIME类型（如：application/vnd.openxmlformats-officedocument.spreadsheetml.sheet）
	FileType string
}

// SaveFile 从HTTP请求中读取上传的文件并保存到本地磁盘。
//
// 文件存储路径规则：uploadPath/年/月/UUID.扩展名
// 例如：uploads/2024/01/a1b2c3d4-e5f6-7890-abcd-ef1234567890.xlsx
//
// 参数：
//   - r          : HTTP请求对象
//   - uploadPath : 文件上传根目录路径
//
// 返回：
//   - *FileInfo : 上传成功后的文件信息
//   - error     : 上传失败时的错误信息
func SaveFile(r *http.Request, uploadPath string) (*FileInfo, error) {
	// 限制请求体大小，防止内存溢出
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		return nil, fmt.Errorf("解析multipart表单失败（文件大小可能超过50MB）: %w", err)
	}

	// 从表单中获取上传的文件
	file, header, err := r.FormFile(FormFieldName)
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败（字段名必须为'file'）: %w", err)
	}
	defer file.Close()

	// 提取原始文件名和扩展名
	originName := header.Filename
	ext := strings.ToLower(filepath.Ext(originName))

	// 生成UUID文件名，避免文件名冲突
	newFileName := uuid.New().String() + ext

	// 按年月构建目录结构，便于管理
	now := time.Now()
	subDir := filepath.Join(fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	saveDirPath := filepath.Join(uploadPath, subDir)

	// 确保目录存在
	if err = os.MkdirAll(saveDirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 完整文件存储路径
	saveFilePath := filepath.Join(saveDirPath, newFileName)

	// 创建目标文件
	dst, err := os.Create(saveFilePath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 将上传文件内容复制到目标文件
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("保存文件内容失败: %w", err)
	}

	// 获取文件MIME类型（从Header中读取）
	fileType := header.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// 构建访问URL（使用正斜杠，跨平台兼容）
	fileURL := "/" + strings.ReplaceAll(filepath.Join(uploadPath, subDir, newFileName), "\\", "/")

	return &FileInfo{
		Filename:   newFileName,
		OriginName: originName,
		FilePath:   saveFilePath,
		FileURL:    fileURL,
		FileSize:   fileSize,
		FileType:   fileType,
	}, nil
}
