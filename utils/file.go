package utils

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var FileURL string

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type item struct {
	title, desc string
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func GetFiles(dir string) ([]list.Item, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	items := make([]list.Item, len(files))
	for i, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting file info: %w", err)
		}

		size := formatSize(fileInfo.Size())
		var desc string
		if fileInfo.IsDir() {
			innerFiles, err := ioutil.ReadDir(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("error reading directory: %w", err)
			}
			desc = fmt.Sprintf("Directory - %d items", len(innerFiles))
		} else {
			desc = fileInfo.Mode().String() + " " + size
		}
		items[i] = item{title: file.Name(), desc: desc}
	}
	return items, nil
}

func UploadFile(filePath string) error {
	conf, err := GetConf()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	req, err := http.NewRequest("PUT", "http://v0.api.upyun.com/"+conf.Bucket+"/"+filePath, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.SetBasicAuth(conf.Operator, conf.Secret)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	FileURL = conf.Bucketurl + "/" + filePath
	if err := clipboard.WriteAll(FileURL); err != nil {
		return fmt.Errorf("error copying file URL to clipboard: %w", err)
	}

	return resp.Body.Close()
}
