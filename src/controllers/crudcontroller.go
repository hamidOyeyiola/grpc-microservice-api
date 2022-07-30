package controller

import (
	"net/http"

	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
)

type CRUDAPIInitializer interface {
	CRUDAPIInitialize(objectTag string,
		validateTag string,
		dataObject model.Model,
		validateObject model.Model)
}

type Creater interface {
	Create(rw http.ResponseWriter, req *http.Request)
}

type ValidateCreater interface {
	ValidateCreate(rw http.ResponseWriter, req *http.Request)
}

type Retriever interface {
	Retrieve(rw http.ResponseWriter, req *http.Request)
}

type ValidateRetriever interface {
	ValidateRetrieve(rw http.ResponseWriter, req *http.Request)
}

type Updater interface {
	Update(rw http.ResponseWriter, req *http.Request)
}

type ValidateUpdater interface {
	ValidateUpdate(rw http.ResponseWriter, req *http.Request)
}

type Deleter interface {
	Delete(rw http.ResponseWriter, req *http.Request)
}

type ValidateDeleter interface {
	ValidateDelete(rw http.ResponseWriter, req *http.Request)
}

type CreateRetrieveUpdateDeleter interface {
	Creater
	Retriever
	Updater
	Deleter
	ValidateCreater
	ValidateRetriever
	ValidateUpdater
	ValidateDeleter
	CRUDAPIInitializer
}
