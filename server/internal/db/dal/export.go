package dal

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"os"
//	"time"
//
//	"github.com/mats0319/unnamed_plan/server/cmd/model"
//	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
//)
//
//// todo: 定时任务+分页查询
//func export(pageNum int) {
//	pageSize := 100 // hard code pageSize=100
//
//	// get data need export
//	qu := Q.User
//	usersFromDB, err := qu.WithContext(context.TODO()).Where(qu.ExportedAt.LtCol(qu.UpdatedAt)).
//		Order(qu.ExportID.Asc()).
//		Offset(pageSize * (pageNum - 1)).Limit(pageSize).Find()
//	if err != nil {
//		mlog.Log("export user failed", mlog.Field("error", err))
//		return
//	}
//
//	if len(usersFromDB) < 1 {
//		return // this page no need to export
//	}
//
//	// read file
//	startID := pageSize*(pageNum-1) + 1
//	endID := pageSize * pageNum
//	filePath := fmt.Sprintf("./export/%s_%dto%d.json", qu.TableName(), startID, endID)
//
//	fileIns, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
//	if err != nil {
//		mlog.Log("open file failed", mlog.Field("error", err))
//		return
//	}
//	defer fileIns.Close()
//
//	dataBytes, err := os.ReadFile(filePath)
//	if err != nil {
//		mlog.Log("read file failed", mlog.Field("error", err))
//		return
//	}
//
//	usersFromFile := make([]*model.User, 0)
//	if len(dataBytes) > 0 { // have file
//		err = json.Unmarshal(dataBytes, &usersFromFile)
//		if err != nil {
//			mlog.Log("json unmarshal failed", mlog.Field("error", err))
//			return
//		}
//
//		timestamp := time.Now().UnixMilli()
//		for i := range usersFromDB {
//			usersFromDB[i].ExportedAt = timestamp
//
//			for j := range usersFromFile {
//				if usersFromDB[i].ExportID == usersFromFile[j].ExportID {
//					usersFromFile[j] = usersFromDB[i]
//					break
//				}
//			}
//		}
//	} else { // create new file
//		timestamp := time.Now().UnixMilli()
//		for i := range usersFromDB {
//			usersFromDB[i].ExportedAt = timestamp
//		}
//
//		usersFromFile = usersFromDB
//	}
//
//	// save file
//	dataBytes, err = json.Marshal(usersFromFile)
//	if err != nil {
//		mlog.Log("json marshal failed", mlog.Field("error", err))
//		return
//	}
//	_, err = fileIns.Write(dataBytes)
//	if err != nil {
//		mlog.Log("write file failed", mlog.Field("error", err))
//		return
//	}
//
//	// update db
//	err = qu.WithContext(context.TODO()).Save(usersFromDB...)
//	if err != nil {
//		mlog.Log("save user failed", mlog.Field("error", err))
//		return
//	}
//}
