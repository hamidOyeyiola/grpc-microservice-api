package controller

import (
	"fmt"
	"net/http"
)

type priority uint8

const (
	statusOK priority = 1 + iota
	statusBadRequest
	statusCreated
	statusNotModified
	statusNotFound
	statusFailedDependency
)

type Response struct {
	header *Header
	body   *Body
}

type Header struct {
	Priority priority
	Code     int
}

func (h *Header) AddPriority(p priority) *Header {
	h.Priority = p
	return h
}

func (h *Header) AddCode(c int) *Header {
	h.Code = c
	return h
}

type Body struct {
	ContentType string
	Content     string
}

func (b *Body) AddContentType(c string) *Body {
	b.ContentType = c
	return b
}

func (b *Body) AddContent(c string) *Body {
	b.Content = c
	return b
}

func (res *Response) Write(rw http.ResponseWriter) *Response {
	rw.Header().Set("Content-Type", res.body.ContentType)
	rw.WriteHeader(res.header.Code)
	fmt.Fprintf(rw, "%s", res.body.Content)
	return res
}

func (res *Response) AddHeader(h *Header) *Response {
	if res.header == nil || h.Priority > res.header.Priority {
		res.header = h
	}
	return res
}

func (res *Response) AddBody(b *Body) *Response {
	res.body = b
	return res
}

var statusOKHeader *Header
var statusOKBody *Body

func GetStatusOKRes() (*Header, *Body) {
	if statusOKHeader == nil {
		statusOKHeader = new(Header).
			AddPriority(statusOK).
			AddCode(http.StatusOK)
		statusOKBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("OK")
	}
	return statusOKHeader, statusOKBody
}

var statusBadRequestHeader *Header
var statusBadRequestBody *Body

func GetStatusBadRequestRes() (*Header, *Body) {
	if statusBadRequestHeader == nil {
		statusBadRequestHeader = new(Header).
			AddPriority(statusBadRequest).
			AddCode(http.StatusBadRequest)
		statusBadRequestBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("Bad Request.")
	}
	return statusBadRequestHeader, statusBadRequestBody
}

var statusNotFoundHeader *Header
var statusNotFoundBody *Body

func GetStatusNotFoundRes() (*Header, *Body) {
	if statusNotFoundHeader == nil {
		statusNotFoundHeader = new(Header).
			AddPriority(statusNotFound).
			AddCode(http.StatusNotFound)
		statusNotFoundBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("Not Found")
	}
	return statusNotFoundHeader, statusNotFoundBody
}

var statusNotModifiedHeader *Header
var statusNotModifiedBody *Body

func GetStatusNotModifiedRes() (*Header, *Body) {
	if statusNotModifiedHeader == nil {
		statusNotModifiedHeader = new(Header).
			AddPriority(statusNotModified).
			AddCode(http.StatusNotModified)
		statusNotModifiedBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("Not Modified.")
	}
	return statusNotModifiedHeader, statusNotModifiedBody
}

var statusFailedDependencyHeader *Header
var statusFailedDependencyBody *Body

func GetStatusFailedDependencyRes() (*Header, *Body) {
	if statusFailedDependencyHeader == nil {
		statusFailedDependencyHeader = new(Header).
			AddPriority(statusFailedDependency).
			AddCode(http.StatusFailedDependency)
		statusFailedDependencyBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("Failed Dependency..")
	}
	return statusFailedDependencyHeader, statusFailedDependencyBody
}

var statusCreatedHeader *Header
var statusCreatedBody *Body

func GetStatusCreatedRes() (*Header, *Body) {
	if statusCreatedHeader == nil {
		statusCreatedHeader = new(Header).
			AddPriority(statusCreated).
			AddCode(http.StatusCreated)
		statusCreatedBody = new(Body).
			AddContentType("text/plaintext").
			AddContent("Created.")
	}
	return statusCreatedHeader, statusCreatedBody
}
