package controller

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
	"google.golang.org/grpc"

	context "context"
	"log"

	mailbox "github.com/hamidOyeyiola/grpc-mail-service-api/mail_service"
)

type GrpcServiceController struct {
	dataSource   string
	db           *sql.DB
	conns        int
	mu           sync.Mutex
	serviceTag   string
	validaterTag string
	validater    model.Model
}

func (sc *GrpcServiceController) ServiceAPIInitialize(serviceTag string, validaterTag string, validater model.Model) {

	sc.serviceTag = serviceTag
	sc.validaterTag = validaterTag
	sc.validater = validater
}

func (sc *GrpcServiceController) Open() bool {
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

func (sc *GrpcServiceController) Close() bool {
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

func NewGrpcServiceController(datasrc string) *GrpcServiceController {
	return &GrpcServiceController{dataSource: datasrc}
}

func (sc *GrpcServiceController) validate(rw http.ResponseWriter, req *http.Request, response *Response) (string, bool) {
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

func (sc *GrpcServiceController) Serve(rw http.ResponseWriter, req *http.Request) {
	response := new(Response)
	ok := sc.Open()
	defer sc.Close()

	var value string

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	value, ok = sc.validate(rw, req, response)

	if !ok {
		h, b := GetStatusFailedDependencyRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	routine := req.URL.Query().Get("routine")
	switch routine {
	case "register":
		register(value, response)
	default:
		post(value, routine, req.Body, response)
	}
	response.Write(rw)
}

func (sc *GrpcServiceController) deleteHelper(q model.SQLQueryToDelete, rw http.ResponseWriter, response *Response) bool {

	_, err := sc.db.Exec(string(q))
	if err != nil {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return false
	}
	return true
}

func register(userid string, response *Response) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	mc := mailbox.NewMailServiceClient(conn)

	u := mailbox.UserID{
		Id: userid,
	}

	r, err := mc.Register(context.Background(), &u)

	if err != nil {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b)
		return
	}
	o, _ := json.MarshalIndent(r.Mails, "", "  ")
	b := new(Body).
		AddContentType("application/json").
		AddContent(string(o))
	h, _ := GetStatusOKRes()
	response.AddHeader(h).
		AddBody(b)
}

func post(userid string, routine string, r io.Reader, response *Response) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	mc := mailbox.NewMailServiceClient(conn)
	m := mailbox.Mail{}
	err = json.NewDecoder(r).Decode(&m)
	if err != nil {
		log.Fatalf("Could not decode message body: %v", err)
	}
	u := mailbox.UserID{
		Id: userid,
	}
	req := mailbox.Request{
		User:    &u,
		Email:   &m,
		Routine: routine,
	}

	n, err := mc.PostRequest(context.Background(), &req)

	if err != nil {
		h, b := GetStatusBadRequestRes()
		response.AddHeader(h).
			AddBody(b)
		return
	}
	o, _ := json.MarshalIndent(n.Mails, "", "  ")
	b := new(Body).
		AddContentType("application/json").
		AddContent(string(o))
	h, _ := GetStatusOKRes()
	response.AddHeader(h).
		AddBody(b)
}
