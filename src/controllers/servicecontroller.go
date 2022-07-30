package controller

import (
	"net/http"

	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
)

type ServiceAPIInitializer interface {
	ServiceAPIInitialize(serviceTag string,
		validaterTag string,
		validater model.Model)
}

type Server interface {
	Serve(rw http.ResponseWriter, req *http.Request)
}

type Service interface {
	Server
	ServiceAPIInitializer
}
