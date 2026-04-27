package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func Backup[T any](t doBackupRecover[T]) {
	dir := "./backup/" + t.Dir()
	if err := os.MkdirAll(dir, 0644); err != nil {
		mlog.Error("mkdir failed", mlog.Field("error", err))
		return
	}

	// if it has data need backup
	var count int64
	if err := dal.DB().Unscoped().Model(t.Model()).Where(t.Condition()).Count(&count).Error; err != nil {
		mlog.Error("db count failed", mlog.Field("error", err))
		return
	}
	if count < 1 { // no data need backup
		return
	}

	// do backup in page
	pageSize := 100
	timestamp := time.Now().UnixMilli()
	for range int(count)/pageSize + 1 {
		dbRecords := t.EmptySlice()
		if err := dal.DB().Unscoped().Model(t.Model()).Where(t.Condition()).Limit(pageSize).Find(&dbRecords).Error; err != nil {
			mlog.Error("get data need to backup failed", mlog.Field("error", err))
			return
		}

		for _, record := range dbRecords {
			// gen file path
			index := uuidToIndex(t.ID(record), 16) // hard code: max 16 files for each table
			filePath := fmt.Sprintf("%s%d.json", dir, index)

			// read data from file
			fileData := t.EmptySlice()
			fileBytes, err := os.ReadFile(filePath) // error if file not exist
			if err == nil {
				if err = json.Unmarshal(fileBytes, &fileData); err != nil {
					mlog.Error("unmarshal file failed", mlog.Field("error", err))
					return
				}
			}

			// set 'backupAt' and update file data
			// 因为这里更改了查询条件涉及的列（备份时间），所以每次分页查询均查询第一页
			t.Update(record, timestamp)

			isExist := false
			for i := range fileData {
				if t.ID(fileData[i]) == t.ID(record) {
					fileData[i] = record
					isExist = true
					break
				}
			}
			if !isExist { // new data which is first time do backup
				fileData = append(fileData, record)
			}

			// write file and update db record
			// 检查：写文件成功但是写数据库失败，下一次会重新尝试备份，而备份函数具有幂等性，所以可以不写在一个事务里
			fileBytes, err = json.Marshal(fileData)
			if err != nil {
				mlog.Error("marshal file failed", mlog.Field("error", err))
				return
			}

			if err := os.WriteFile(filePath, fileBytes, 0644); err != nil { // implicit create file at first time
				mlog.Error("write file failed", mlog.Field("error", err))
				return
			}

			if err := dal.DB().Unscoped().Model(record).UpdateColumns(record).Error; err != nil { // UpdateColumns skip hooks and auto-updateTime
				mlog.Error("update db data failed", mlog.Field("error", err))
				return
			}
		}
	}
}

func uuidToIndex(id uuid.UUID, max int) int {
	var v int
	for i := range id {
		v += int(id[i])
	}

	return v % max
}
