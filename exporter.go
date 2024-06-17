package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/peteradeojo/lamp-logger/internal/database"
	"github.com/sqlc-dev/pqtype"
	"github.com/xuri/excelize/v2"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (apiCfg *ApiConfig) generateLogExport(cx context.Context, appToken, saveName string) {
	job := ExportJob{
		AppToken: appToken,
		Status:   "pending",
		Path:     "",
	}

	_, err := apiCfg.saveExportJobStatus(cx, appToken, job)
	if err != nil {
		return
	}

	logs, err := apiCfg.DB.ExportLogs(cx, appToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(logs)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newLogs []Log
	err = json.Unmarshal(data, &newLogs)
	if err != nil {
		fmt.Println(err)
		return
	}

	f := excelize.NewFile()

	defer func() {
		dir := path.Dir(saveName)

		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm|os.ModeDir)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		if err := f.SaveAs(saveName); err != nil {
			fmt.Println(err)
		}

		if err := f.Close(); err != nil {
			fmt.Println(err)
		}

		job.Path = saveName
		job.Status = "done"

		apiCfg.saveExportJobStatus(cx, appToken, job)

		go apiCfg.uploadExport(job)

		fmt.Println("Export Done")
	}()

	sheetName := "Export"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	headers := []string{"SN", "Date", "Level", "Text"}
	columns := []string{"A", "B", "C", "D"}
	for i, h := range headers {
		f.SetCellStr(sheetName, fmt.Sprintf("%s1", columns[i]), h)
	}

	for i, log := range newLogs {
		f.SetCellValue(sheetName, fmt.Sprintf("A%v", i+2), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%v", i+2), log.Createdat.Time)
		f.SetCellValue(sheetName, fmt.Sprintf("C%v", i+2), log.Level)
		f.SetCellValue(sheetName, fmt.Sprintf("D%v", i+2), log.Text)
	}
}

func (apiCfg *ApiConfig) saveExportJobStatus(cx context.Context, key string, job ExportJob) ([]byte, error) {
	jobData, _ := json.Marshal(job)
	_, err := apiCfg.redisClient.HSet(cx, "export_jobs", key, jobData).Result()

	if err != nil {
		fmt.Println("Unable to register job:", err)
		apiCfg.DB.CreateSystemLog(cx, database.CreateSystemLogParams{
			Text:  fmt.Sprintf("Unable to register job: %v", err),
			Level: "error",
			Context: pqtype.NullRawMessage{
				RawMessage: jobData,
				Valid:      true,
			},
		})
	}

	return jobData, err
}

func (apiCfg *ApiConfig) uploadExport(job ExportJob) {
	ctx := context.Background()
	filePath := job.Path
	data, _ := apiCfg.saveExportJobStatus(ctx, job.AppToken, job)
	cld, _ := loadUploader()

	result, err := cld.Upload.Upload(ctx, job.Path, uploader.UploadParams{
		Folder:      path.Dir(filePath),
		UseFilename: api.Bool(true),
	})

	if err != nil {
		reportError(ctx, err, pqtype.NullRawMessage{RawMessage: data})
		return
	}

	job.Path = result.SecureURL
	job.Status = "completed"

	apiCfg.saveExportJobStatus(ctx, job.AppToken, job)

	err = os.Remove(filePath)
	if err != nil {
		reportError(ctx, err, pqtype.NullRawMessage{})
		return
	}

	job.Status = "pruned"
	apiCfg.saveExportJobStatus(ctx, job.AppToken, job)
}

func loadUploader() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
}
