package backup

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"gorm.io/gorm/clause"
)

func Recover[T any](t doBackupRecover[T]) error {
	dir := "./recover/" + t.Dir()
	entry, err := os.ReadDir(dir)
	if err != nil {
		mlog.Error(fmt.Sprintf("read dir '%s' failed, error: %v", dir, err))
		return err
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		var fileInfo fs.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			mlog.Error("get file info failed", mlog.Field("error", err))
			continue
		}

		if !strings.HasSuffix(fileInfo.Name(), ".json") {
			continue // ignore not go files
		}

		if err := recoverFile(dir+fileInfo.Name(), t); err != nil {
			return err
		}
	}

	return nil
}

func recoverFile[T any](path string, t doBackupRecover[T]) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		mlog.Error("read file failed", mlog.Field("error", err))
		return err
	}

	fileData := t.EmptySlice()
	if err := json.Unmarshal(fileBytes, &fileData); err != nil {
		mlog.Error("unmarshal file failed", mlog.Field("error", err))
		return err
	}

	// 仅测试使用，为了方便看出一条数据库记录是预设的，还是恢复的
	//{
	//	for i := range fileData {
	//		t.DoSomeChangeForTest(fileData[i])
	//	}
	//}

	clauseSkipAutoTime := clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(t.ColumnNames()), // all fields
	}

	if err := dal.DB().Model(t.Model()).Clauses(clauseSkipAutoTime).Create(fileData).Error; err != nil {
		mlog.Error("save file failed", mlog.Field("error", err))
		return err
	}

	return nil
}
