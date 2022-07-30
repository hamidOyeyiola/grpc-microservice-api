package model

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamidOyeyiola/grpc-mail-service-api/interfaces"
)

type SQLQueryToInsert string

type SQLQueryToSelect string

type SQLQueryToDelete string

type SQLQueryToUpdate struct {
	Stmts  []string
	Values []string
	ID     string
}

type JSONObject string

type Model interface {
	Name() string
	PrimaryKey() string
	Validate(interfaces.Iterator, io.Reader) (string, bool, SQLQueryToDelete)
	FromQueryResult(interfaces.Iterator) (JSONObject, int)
	FromQueryResultArray([]interfaces.Iterator) (JSONObject, int)
	InsertInto(jsonObject io.Reader) ([]SQLQueryToInsert, []SQLQueryToSelect, bool)
	InsertIntoWhere(jsonObject io.Reader, value string) (SQLQueryToInsert, SQLQueryToSelect, bool)
	InsertIntoIf(jsonObject io.Reader, value string) (SQLQueryToInsert, SQLQueryToSelect, bool)
	UpdateWhere(jsonObject io.Reader, value string) (SQLQueryToUpdate, SQLQueryToSelect, bool)
	UpdateIf(jsonObject io.Reader, value string) (SQLQueryToUpdate, SQLQueryToSelect, bool)
	Update(jsonObject io.Reader) ([]SQLQueryToUpdate, []SQLQueryToSelect, bool)
	SelectFromWhere(value string) (SQLQueryToSelect, bool)
	DeleteFromWhere(value string) (SQLQueryToDelete, bool)
}

func GetQueryToSelect(m Model, id string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", m.Name(), id))
	return q, true
}

func GetQueryToSelectAll(m Model) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s ", m.Name()))
	return q, true
}

func GetQueryToDelete(m Model, id string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", m.Name(), id))
	return q, ok
}

func GetParamFromRequest(req *http.Request, param string) (string, bool) {
	vars := mux.Vars(req)
	id, ok := vars[param]
	return id, ok
}
