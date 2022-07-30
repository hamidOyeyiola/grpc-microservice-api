package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	controller "github.com/hamidOyeyiola/grpc-mail-service-api/controllers"
	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
)

func MakeCRUDAPI(rt *mux.Router, crud controller.CreateRetrieveUpdateDeleter,
	path string, objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	if dataObject == nil || objectTag == validaterTag {
		return
	}

	if validater == nil || validaterTag == "" {
		s := fmt.Sprintf("%s/{%s}", path, objectTag)
		addCreater(rt, path, crud.Create)
		addRetriever(rt, path, crud.Retrieve)
		addRetriever(rt, s, crud.Retrieve)
		addUpdater(rt, path, crud.Update)
		addDeleter(rt, s, crud.Delete)
	} else {
		s := fmt.Sprintf("%s/{%s}/{%s}", path, objectTag, validaterTag)
		s2 := fmt.Sprintf("%s/{%s}", path, validaterTag)
		addCreater(rt, s2, crud.ValidateCreate)
		addRetriever(rt, s2, crud.ValidateRetrieve)
		addRetriever(rt, s, crud.ValidateRetrieve)
		addUpdater(rt, s2, crud.ValidateUpdate)
		addDeleter(rt, s, crud.ValidateDelete)
		addDeleter(rt, s2, crud.ValidateDelete)
	}

	crud.CRUDAPIInitialize(objectTag, validaterTag, dataObject, validater)
	return
}

func MakeCreaterAPI(rt *mux.Router, crud controller.CreateRetrieveUpdateDeleter,
	path string, objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	if dataObject == nil || objectTag == validaterTag {
		return
	}

	if validater == nil || validaterTag == "" {
		addCreater(rt, path, crud.Create)
	} else {
		s := fmt.Sprintf("%s/{%s}", path, validaterTag)
		addCreater(rt, s, crud.ValidateCreate)
	}

	crud.CRUDAPIInitialize(objectTag, validaterTag, dataObject, validater)
	return
}

func MakeRetrieverAPI(rt *mux.Router, crud controller.CreateRetrieveUpdateDeleter,
	path string, objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	if dataObject == nil || objectTag == validaterTag {
		return
	}

	if validater == nil || validaterTag == "" {
		s := fmt.Sprintf("%s/{%s}", path, objectTag)
		addRetriever(rt, path, crud.Retrieve)
		addRetriever(rt, s, crud.Retrieve)
	} else {
		s := fmt.Sprintf("%s/{%s}/{%s}", path, objectTag, validaterTag)
		s2 := fmt.Sprintf("%s/{%s}", path, validaterTag)
		addRetriever(rt, s2, crud.ValidateRetrieve)
		addRetriever(rt, s, crud.ValidateRetrieve)
	}

	crud.CRUDAPIInitialize(objectTag, validaterTag, dataObject, validater)
	return
}

func MakeUpdaterAPI(rt *mux.Router, crud controller.CreateRetrieveUpdateDeleter,
	path string, objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	if dataObject == nil || objectTag == validaterTag {
		return
	}

	if validater == nil || validaterTag == "" {
		addUpdater(rt, path, crud.Update)
	} else {
		s := fmt.Sprintf("%s/{%s}", path, validaterTag)
		addUpdater(rt, s, crud.ValidateUpdate)
	}

	crud.CRUDAPIInitialize(objectTag, validaterTag, dataObject, validater)
	return
}

func MakeDeleterAPI(rt *mux.Router, crud controller.CreateRetrieveUpdateDeleter,
	path string, objectTag string, validaterTag string, dataObject model.Model, validater model.Model) {

	if dataObject == nil || objectTag == validaterTag {
		return
	}

	if validater == nil || validaterTag == "" {
		s := fmt.Sprintf("%s/{%s}", path, objectTag)
		addDeleter(rt, s, crud.Delete)
	} else {
		s := fmt.Sprintf("%s/{%s}/{%s}", path, objectTag, validaterTag)
		s2 := fmt.Sprintf("%s/{%s}", path, validaterTag)
		addDeleter(rt, s, crud.ValidateDelete)
		addDeleter(rt, s2, crud.ValidateDelete)
	}

	crud.CRUDAPIInitialize(objectTag, validaterTag, dataObject, validater)
	return
}

func addCreater(rt *mux.Router, path string, c http.HandlerFunc) {
	rt.HandleFunc(path, c).Methods("POST")
}

func addRetriever(rt *mux.Router, path string, r http.HandlerFunc) {
	rt.HandleFunc(path, r).Methods("GET")
}

func addUpdater(rt *mux.Router, path string, u http.HandlerFunc) {
	rt.HandleFunc(path, u).Methods("PUT")
}

func addDeleter(rt *mux.Router, path string, d http.HandlerFunc) {
	rt.HandleFunc(path, d).Methods("DELETE")
}
