package backup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func Backup[T any](t doBackupRecover[T]) error {
	dir := "./backup/" + t.Dir()
	err := emptyDir(dir)
	if err != nil {
		mlog.Error("mkdir failed", mlog.Field("error", err))
		return err
	}

	dbData := t.EmptySlice()
	err = dal.DB().Model(t.Model()).Where(t.Condition()).Find(&dbData).Error // todo: query in page
	if err != nil {
		mlog.Error("get data need to backup failed", mlog.Field("error", err))
		return err
	}

	for _, record := range dbData {
		// gen file path
		index := uuidToIndex(t.ID(record), 16) // hard code: max 16 files for each table
		filePath := fmt.Sprintf("%s%d.json", dir, index)

		// read data from file
		fileData := t.EmptySlice()
		fileBytes, err := os.ReadFile(filePath) // error if file not exist
		if err == nil {
			err = json.Unmarshal(fileBytes, &fileData)
			if err != nil {
				mlog.Error("unmarshal file failed", mlog.Field("error", err))
				return err
			}
		}

		// set 'backupAt', update file data
		t.Update(record)

		isExist := false
		for i := range fileData {
			if t.ID(fileData[i]) == t.ID(record) {
				fileData[i] = record
				isExist = true
				break
			}
		}

		if !isExist { // 如果是新的待备份数据
			fileData = append(fileData, record)
		}

		// write file, update db record
		fileBytes, err = json.Marshal(fileData)
		if err != nil {
			mlog.Error("marshal file failed", mlog.Field("error", err))
			return err
		}

		err = os.WriteFile(filePath, fileBytes, 0644)
		if err != nil {
			mlog.Error("write file failed", mlog.Field("error", err))
			return err
		}

		// think: if necessary get 'id' from method? (like: where(sprintf("%s = ?", t.idName()), t.id(record)) )
		err = dal.DB().Model(t.Model()).Where("id = ?", t.ID(record)).Save(record).Error
		if err != nil {
			mlog.Error("update db data failed", mlog.Field("error", err))
			return err
		}
	}

	return nil
}

func uuidToIndex(id uuid.UUID, max int) int {
	var v int
	for i := range id {
		v += int(id[i])
	}

	return v % max
}

func emptyDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		mlog.Error("remove dir failed", mlog.Field("error", err))
		return err
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		mlog.Error("mkdir failed", mlog.Field("error", err))
		return err
	}

	return nil
}
