package model

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hamidOyeyiola/grpc-mail-service-api/interfaces"
	"github.com/hamidOyeyiola/grpc-mail-service-api/utils"
)

type Session struct {
	Sessionkey string `json:"sessionkey"`
	UserKey    string `json:"userkey"`
	Deadline   string `json:"-"`
}

func (se Session) String() string {
	return fmt.Sprintf("User(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		se.Sessionkey, se.UserKey, se.Deadline)
}

func (se Session) Name() string {
	return "session"
}

func (se Session) PrimaryKey() string {
	return "sessionkey"
}

func (se Session) Validate(s interfaces.Iterator, r io.Reader) (string, bool, SQLQueryToDelete) {
	var id int
	if s.Next() {
		err := s.Scan(&se.Sessionkey, &se.UserKey, &se.Deadline, &id)
		if err != nil {
			return "", false, ""
		}
	}
	d, _ := se.DeleteFromWhere(se.Sessionkey)
	return se.UserKey, utils.IsNotDeadline(se.Deadline), d
}

func (se Session) InsertInto(r io.Reader) (ins []SQLQueryToInsert, sel []SQLQueryToSelect, ok bool) {
	v := []struct {
		Email string
	}{}
	err := json.NewDecoder(r).Decode(&v)
	for _, s := range v {
		ins = append(ins, SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
			utils.GetSessionToken(), s.Email, utils.MakeDeadline(time.Minute*5))))
	}
	sel = append(sel, "")
	return ins, sel, err == nil
}

func (se Session) InsertIntoWhere(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {

	ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		utils.GetSessionToken(), value, utils.MakeDeadline(time.Minute*5)))
	return ins, sel, true
}

func (se Session) InsertIntoIf(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		utils.GetSessionToken(), value, utils.MakeDeadline(time.Minute*5)))
	return ins, sel, true
}

func (se Session) Update(r io.Reader) (upt []SQLQueryToUpdate, sel []SQLQueryToSelect, ok bool) {
	return upt, sel, true
}

func (se Session) UpdateWhere(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	return upt, sel, true
}

func (se Session) UpdateIf(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	return upt, sel, true
}

func (se Session) FromQueryResult(s interfaces.Iterator) (JSONObject, int) {
	v := []Session{}
	var n int
	var id int
	for n = 0; s.Next(); n++ {
		err := s.Scan(&se.Sessionkey, &se.UserKey, &se.Deadline, &id)
		if err != nil {
			break
		}
		v = append(v, se)
	}
	o, _ := json.MarshalIndent(v, "", "  ")
	return JSONObject(o), n
}

func (se Session) FromQueryResultArray(ss []interfaces.Iterator) (JSONObject, int) {
	v := []Session{}
	var n int
	var id int
	for _, s := range ss {
		for n = 0; s.Next(); n++ {
			err := s.Scan(&se.Sessionkey, &se.Deadline, &id)
			if err != nil {
				break
			}
			v = append(v, se)
		}
	}
	o, _ := json.MarshalIndent(v, "", "  ")
	return JSONObject(o), n
}

func (se Session) SelectFromWhere(value string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect("SELECT * FROM " + se.Name())
	q = q + SQLQueryToSelect(fmt.Sprintf(" WHERE %s = '%s'", se.PrimaryKey(), value))
	return q, true
}

func (se Session) DeleteFromWhere(value string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", se.Name(), se.PrimaryKey(), value))
	return q, true
}
