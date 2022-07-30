package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hamidOyeyiola/grpc-mail-service-api/interfaces"
	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLCRUDController struct {
	dataSource   string
	db           *sql.DB
	conns        int
	mu           sync.Mutex
	objectTag    string
	validaterTag string
	dataObject   model.Model
	validater    model.Model
}

func (sc *MySQLCRUDController) CRUDAPIInitialize(objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	sc.objectTag = objectTag
	sc.validaterTag = validaterTag
	sc.dataObject = dataObject
	sc.validater = validater
}

func (sc *MySQLCRUDController) Open() bool {
	sc.mu.Lock()
	if sc.db == nil {
		db, err := sql.Open("mysql", sc.dataSource)
		if err != nil {
			return false
		}
		db.SetConnMaxLifetime(3 * time.Minute)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		sc.db = db
	}
	sc.conns++
	sc.mu.Unlock()
	return true
}

func (sc *MySQLCRUDController) Close() bool {
	sc.mu.Lock()
	sc.conns--
	if sc.conns == 0 {
		err := sc.db.Close()
		if err != nil {
			return false
		}
		sc.db = nil
	}
	sc.mu.Unlock()
	return true
}

func NewMySQLCRUDController(datasrc string) *MySQLCRUDController {
	return &MySQLCRUDController{dataSource: datasrc}
}

func (sc *MySQLCRUDController) Validate(rw http.ResponseWriter, req *http.Request, response *Response) (string, bool) {

	value, ok := model.GetParamFromRequest(req, sc.validaterTag)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b)
		return "", false
	}
	q, _ := sc.validater.SelectFromWhere(value)
	res, err := sc.db.Query(string(q))
	if err != nil {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b)
		return "", false
	}
	value2, ok, del := sc.validater.Validate(res, req.Body)
	if !ok {
		if del != "" {
			sc.deleteHelper(del, rw, response)
		}
		return "", false
	}
	return value2, ok
}

func (sc *MySQLCRUDController) Create(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()
	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	q := make([]model.SQLQueryToInsert, 1)
	r := make([]model.SQLQueryToSelect, 1)
	var value string

	value, ok = model.GetParamFromRequest(req, sc.objectTag)
	if ok {
		q[0], r[0], ok = sc.dataObject.InsertIntoWhere(req.Body, value)
	} else {
		q, r, ok = sc.dataObject.InsertInto(req.Body)
	}

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.createHelper(q, r, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) createHelper(q []model.SQLQueryToInsert, r []model.SQLQueryToSelect, rw http.ResponseWriter, response *Response) {
	ids := string("0")
	for _, v := range q {
		res, err := sc.db.Exec(string(v))
		if err != nil {
			if ids == "0" {
				h, b := GetStatusNotModifiedRes()
				response.AddHeader(h).
					AddBody(b)
				return
			}
			break
		}
		id, _ := res.LastInsertId()
		ids = ids + fmt.Sprintf(",%d", id)
	}

	h, _ := GetStatusCreatedRes()
	response.AddHeader(h)

	if r[0] == "" {
		q, _ := model.GetQueryToSelect(sc.dataObject, ids)
		sc.retrieveHelper(q, rw, response)
	} else {
		sc.batchRetrieveHelper(r, rw, response)
	}
}

func (sc *MySQLCRUDController) ValidateCreate(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	q := make([]model.SQLQueryToInsert, 1)
	r := make([]model.SQLQueryToSelect, 1)
	var value string

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	value, ok = sc.Validate(rw, req, response)

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	q[0], r[0], ok = sc.dataObject.InsertIntoIf(req.Body, value)

	if !ok {
		h, b := GetStatusNotModifiedRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.createHelper(q, r, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) Retrieve(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	var q model.SQLQueryToSelect
	var value string

	value, ok = model.GetParamFromRequest(req, sc.objectTag)
	if ok {
		q, ok = sc.dataObject.SelectFromWhere(value)
	} else {
		q, ok = model.GetQueryToSelectAll(sc.dataObject)
	}
	sc.retrieveHelper(q, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) retrieveHelper(q model.SQLQueryToSelect, rw http.ResponseWriter, response *Response) {

	res, err := sc.db.Query(string(q))
	defer res.Close()

	if err != nil {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b)
		return
	}
	d, n := sc.dataObject.FromQueryResult(res)
	if n == 0 {
		h, b := GetStatusNotFoundRes()
		response.AddHeader(h).
			AddBody(b)
		return
	}
	b := new(Body).
		AddContentType("application/json").
		AddContent(string(d))
	h, _ := GetStatusOKRes()
	response.AddHeader(h).
		AddBody(b)
}

func (sc *MySQLCRUDController) batchRetrieveHelper(q []model.SQLQueryToSelect, rw http.ResponseWriter, response *Response) {

	res := make([]interfaces.Iterator, len(q))
	for i, v := range q {
		resp, err := sc.db.Query(string(v))
		defer resp.Close()
		if err != nil {
			h, b := GetStatusFailedDependencyRes()
			response.AddHeader(h).
				AddBody(b)
		}
		res[i] = resp
	}

	d, n := sc.dataObject.FromQueryResultArray(res)
	if n == 0 {
		h, b := GetStatusNotFoundRes()
		response.AddHeader(h).
			AddBody(b)
		return
	}
	//fmt.Println(d)
	b := new(Body).
		AddContentType("application/json").
		AddContent(string(d))
	h, _ := GetStatusOKRes()
	response.AddHeader(h).
		AddBody(b)
}

func (sc *MySQLCRUDController) ValidateRetrieve(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	value, ok := sc.Validate(rw, req, response)
	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	q, ok := sc.dataObject.SelectFromWhere(value)
	if !ok {
		h, b := GetStatusNotFoundRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.retrieveHelper(q, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) Update(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	q := make([]model.SQLQueryToUpdate, 1)
	r := make([]model.SQLQueryToSelect, 1)
	var value string

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	value, ok = model.GetParamFromRequest(req, sc.objectTag)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	q[0], r[0], ok = sc.dataObject.UpdateWhere(req.Body, value)
	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.updateHelper(q, r, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) updateHelper(q []model.SQLQueryToUpdate, r []model.SQLQueryToSelect, rw http.ResponseWriter, response *Response) {
	for i := 0; i < len(q); i++ {
		for j := 0; j < len(q[i].Stmts); j++ {
			stmt, err := sc.db.Prepare(string(q[i].Stmts[j]))
			if err != nil {
				h, b := GetStatusNotModifiedRes()
				response.AddHeader(h).
					AddBody(b)
				return
			}
			_, err = stmt.Exec(q[i].Values[j], q[i].ID)
			if err != nil {
				h, b := GetStatusNotModifiedRes()
				response.AddHeader(h).
					AddBody(b)
				return
			}
		}
	}
	sc.batchRetrieveHelper(r, rw, response)
}

func (sc *MySQLCRUDController) ValidateUpdate(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	q := make([]model.SQLQueryToUpdate, 1)
	r := make([]model.SQLQueryToSelect, 1)
	var value string

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	value, ok = sc.Validate(rw, req, response)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	q[0], r[0], ok = sc.dataObject.UpdateIf(req.Body, value)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.updateHelper(q, r, rw, response)
	response.Write(rw)
}

func (sc *MySQLCRUDController) Delete(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)

	ok := sc.Open()
	defer sc.Close()

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	value, ok := model.GetParamFromRequest(req, sc.objectTag)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	var r model.SQLQueryToSelect
	r, ok = sc.dataObject.SelectFromWhere(value)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.retrieveHelper(r, rw, response)
	q, _ := sc.dataObject.DeleteFromWhere(value)
	ok = sc.deleteHelper(q, rw, response)
	if !ok {
		q, _ := model.GetQueryToDelete(sc.dataObject, sc.objectTag)
		sc.deleteHelper(q, rw, response)
	}
	response.Write(rw)
}

func (sc *MySQLCRUDController) deleteHelper(q model.SQLQueryToDelete, rw http.ResponseWriter, response *Response) bool {
	_, err := sc.db.Exec(string(q))
	if err != nil {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b)
		return false
	}
	return true
}

func (sc *MySQLCRUDController) ValidateDelete(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}

	value, ok := sc.Validate(rw, req, response)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	q, ok := sc.dataObject.SelectFromWhere(value)
	if ok {
		sc.retrieveHelper(q, rw, response)
	}

	d, ok := sc.dataObject.DeleteFromWhere(value)
	if !ok {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	sc.deleteHelper(d, rw, response)
	response.Write(rw)
}
