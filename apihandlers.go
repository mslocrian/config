package main

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"strings"
	//"net/url"
	"strconv"
	"github.com/nu7hatch/gouuid"
)

func GetConfigObj(r *http.Request, obj models.ConfigObj) (models.ConfigObj, error) {
	var retObj models.ConfigObj
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return retObj, err
	}
	if err := r.Body.Close(); err != nil {
		return retObj, err
	}
	return obj.UnmarshalObject(body)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/jsoni;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		panic(err)
	}
}

func ShowConfigObject(w http.ResponseWriter, r *http.Request) {
	logger.Println("####  ShowConfigObject called")
}

func ConfigObjectsBulkGet(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), "/")
	resource = strings.Split(resource, "?")[0]
	resource = resource[:len(resource)-1]

	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		currentIndex, objCount := ExtractGetBulkParams(r)
		gMgr.objHdlMap[resource].owner.GetBulkObject(obj, currentIndex, objCount)
	}
}

func ExtractGetBulkParams(r *http.Request) (currentIndex int64, objectCount int64) {
	valueMap := r.URL.Query()
	if currentIndexStr, ok1 := valueMap["CurrentMarker"]; ok1 {
		currentIndex, _ = strconv.ParseInt(currentIndexStr[0], 10, 64)
	} else {
		currentIndex = 100
	}

	if objectCountStr, ok := valueMap["Count"]; ok {
		objectCount, _ = strconv.ParseInt(objectCountStr[0], 10, 64)
	} else {
		objectCount = 100
	}
	return currentIndex, objectCount
}

func ConfigObjectCreate(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimPrefix(r.URL.String(), "/")
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		//_, success := gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
		objectId, success := gMgr.objHdlMap[resource].owner.CreateObject(obj, gMgr.dbHdl)
		if success == true {

			UUId, err := uuid.NewV4()
			if err != nil {
			    logger.Println("### Failed to get UUID ", UUId, err)
			}
			objKey, err := obj.GetKey()
			if err != nil {
				logger.Println("### Failed to get objKey after executing ", objKey, err)
			}
			dbCmd := fmt.Sprintf(`INSERT INTO UuidMap (Uuid, Key) VALUES ('%v', '%v') ;`, UUId, objectId)
			_, err = models.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
			if err != nil {
				logger.Println("### Failed to insert uuid entry in db ", dbCmd, err)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err = json.NewEncoder(w).Encode(UUId); err != nil {
				logger.Println("### Failed to encode the UUId for object ", resource, UUId)
			}
		}
	}
	return
}

func ConfigObjectDelete(w http.ResponseWriter, r *http.Request) {
	var objKey int64
	resource := strings.Split(r.URL.String(), "/")[1]
	vars := mux.Vars(r)

	err := gMgr.dbHdl.QueryRow("select Key from UuidMap where Uuid = ?", vars["UUId"]).Scan(&objKey)
	if err != nil {
		logger.Println("### Failure in getting objKey for Uuid ", resource, vars["UUId"], err)
		return
	}

/*
	objId, err := strconv.ParseInt(vars["objId"], 10, 64)
	if err != nil {
		logger.Println("### Failure in deleting object with Id ", resource, vars["objId"], err)
		return
	}
*/
	if objHdl, ok := models.ConfigObjectMap[resource]; ok {
		obj, _ := GetConfigObj(r, objHdl)
		//success := gMgr.objHdlMap[resource].owner.DeleteObject(obj, objId, gMgr.dbHdl)
		success := gMgr.objHdlMap[resource].owner.DeleteObject(obj, objKey, gMgr.dbHdl)
		if success == true {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			dbCmd := "delete from " + "UuidMap" + " where Uuid = " + vars["UUId"]
			fmt.Println("### DB Deleting Uuid")
			_, err := models.ExecuteSQLStmt(dbCmd, gMgr.dbHdl)
			if err != nil {
				logger.Println("### Failure in deleting Uuid map entry for ", vars["UUId"], err)
			}

		}
	}
	return
}
