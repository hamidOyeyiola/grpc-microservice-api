package model

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hamidOyeyiola/grpc-mail-service-api/interfaces"
	"github.com/hamidOyeyiola/grpc-mail-service-api/utils"
)

type User struct {
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     utils.EmailAddress `json:"email"`
	PhoneNo   string             `json:"phoneno"`
	Password  string             `json:"password"`
	MessageID string             `json:"-"`
	CreatedOn string             `json:"-"`
	UpdatedOn string             `json:"-"`
	ID        int                `json:"-"`
}

func (u User) String() string {
	return fmt.Sprintf("User(firstname, lastname, email, phoneno, password, messageID, createdOn, updatedOn) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s')",
		u.FirstName, u.LastName, u.Email, u.PhoneNo, u.Password, u.MessageID, u.CreatedOn, u.UpdatedOn)
}

func (u User) Name() string {
	return "users"
}

func (u User) PrimaryKey() string {
	return "email"
}

func (u User) Validate(s interfaces.Iterator, r io.Reader) (string, bool, SQLQueryToDelete) {
	v := struct {
		Email    utils.EmailAddress
		Password string
	}{}
	err := json.NewDecoder(r).Decode(&v)
	if s.Next() {
		err = s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &u.ID)
		if err != nil {
			return "", false, ""
		}
	}
	return string(v.Email), (u.Password == utils.EncryptPassword(v.Password) && v.Email == u.Email), ""
}

func (u User) InsertIntoIf(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if u.Email.IsValid() && string(u.Email) == value {
		ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
			u.FirstName, u.LastName, u.Email, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0))
	}
	return ins, sel, err == nil
}

func (u User) InsertInto(r io.Reader) (ins []SQLQueryToInsert, sel []SQLQueryToSelect, ok bool) {
	v := []User{}

	err := json.NewDecoder(r).Decode(&v)

	for _, u = range v {
		if u.Email.IsValid() {
			ins = append(ins, SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
				u.FirstName, u.LastName, u.Email, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0)))
		}
	}
	sel = append(sel, "")
	return ins, sel, err == nil
}

func (u User) InsertIntoWhere(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if utils.EmailAddress(value).IsValid() {
		ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
			u.FirstName, u.LastName, value, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0))
	}
	return ins, sel, err == nil
}

func (u User) Update(r io.Reader) (upt []SQLQueryToUpdate, sel []SQLQueryToSelect, ok bool) {
	v := []User{}
	err := json.NewDecoder(r).Decode(&v)
	upt = make([]SQLQueryToUpdate, len(v))
	sel = make([]SQLQueryToSelect, len(v))
	i := 0
	for _, u = range v {
		upt[i].ID = string(u.Email)
		if u.FirstName != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.FirstName)
		}
		if u.LastName != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.LastName)
		}
		if u.PhoneNo != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.PhoneNo)
		}
		if u.Password != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET password=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, utils.EncryptPassword(u.Password))
		}
		if u.MessageID != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, string(u.MessageID))
		}
		upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
		date := fmt.Sprintf("%s", utils.NewDate())
		upt[i].Values = append(upt[i].Values, date)
		sel[i] = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))
		i++
	}
	return upt, sel, err == nil
}

func (u User) UpdateWhere(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	upt.ID = value
	if u.FirstName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.FirstName)
	}
	if u.LastName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.LastName)
	}
	if u.PhoneNo != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
		upt.Values = append(upt.Values, u.PhoneNo)
	}
	if u.Password != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET password=? WHERE email=?"))
		upt.Values = append(upt.Values, utils.EncryptPassword(u.Password))
	}
	if u.MessageID != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
		upt.Values = append(upt.Values, string(u.MessageID))
	}
	upt.Stmts = append(upt.Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
	date := fmt.Sprintf("%s", utils.NewDate())
	upt.Values = append(upt.Values, date)
	sel = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))

	return upt, sel, err == nil
}

func (u User) UpdateIf(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if string(u.Email) == value {
		upt.ID = string(u.Email)
		if u.FirstName != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
			upt.Values = append(upt.Values, u.FirstName)
		}
		if u.LastName != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
			upt.Values = append(upt.Values, u.LastName)
		}
		if u.PhoneNo != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
			upt.Values = append(upt.Values, u.PhoneNo)
		}
		if u.Password != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET password=? WHERE email=?"))
			upt.Values = append(upt.Values, utils.EncryptPassword(u.Password))
		}
		if u.MessageID != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
			upt.Values = append(upt.Values, string(u.MessageID))
		}
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
		date := fmt.Sprintf("%s", utils.NewDate())
		upt.Values = append(upt.Values, date)
		sel = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))
	}
	return upt, sel, err == nil
}

func (u User) FromQueryResult(s interfaces.Iterator) (JSONObject, int) {
	v := []User{}
	var n int
	var id int
	for n = 0; s.Next(); n++ {
		err := s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &id)
		if err != nil {
			break
		}
		v = append(v, u)
	}
	o, _ := json.MarshalIndent(v, "", "  ")
	return JSONObject(o), n
}

func (u User) FromQueryResultArray(ss []interfaces.Iterator) (JSONObject, int) {
	v := []User{}
	var id int
	n := 0
	for _, s := range ss {
		for ; s.Next(); n++ {
			err := s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &id)
			if err != nil {
				break
			}
			v = append(v, u)
		}
	}
	o, _ := json.MarshalIndent(v, "", "  ")
	return JSONObject(o), n
}

func (u User) SelectFromWhere(value string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect("SELECT * FROM " + u.Name())
	q = q + SQLQueryToSelect(fmt.Sprintf(" WHERE %s = '%s'", u.PrimaryKey(), value))
	return q, true
}

func (u User) DeleteFromWhere(key string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", u.Name(), u.PrimaryKey(), key))
	return q, ok
}
