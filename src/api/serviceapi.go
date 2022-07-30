package api

import (
	"fmt"

	"github.com/gorilla/mux"
	controller "github.com/hamidOyeyiola/grpc-mail-service-api/controllers"
	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
)

func MakeServiceAPI(rt *mux.Router, service controller.Service,
	path string, serviceTag string, validaterTag string, validater model.Model) {

	if validater == nil {
		return
	}

	s := fmt.Sprintf("%s/{%s}", path, validaterTag)
	rt.HandleFunc(s, service.Serve).Methods("POST")
	service.ServiceAPIInitialize(serviceTag, validaterTag, validater)
	return
}
