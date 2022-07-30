package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/hamidOyeyiola/grpc-mail-service-api/api"
	controller "github.com/hamidOyeyiola/grpc-mail-service-api/controllers"
	model "github.com/hamidOyeyiola/grpc-mail-service-api/models"
)

func main() {
	r := mux.NewRouter()
	signUp := controller.NewMySQLCRUDController("hamid:kolajoke2055@tcp(localhost:3306)/registrationandlogin")
	api.MakeCreaterAPI(r, signUp, "/api/signup", "id", "",
		model.User{}, nil)
	signIn := controller.NewMySQLCRUDController("hamid:@tcp(localhost:3306)/registrationandlogin")
	api.MakeCreaterAPI(r, signIn, "/api/signin", "session", "user",
		model.Session{}, model.User{})
	updateUser := controller.NewMySQLCRUDController("hamid:@tcp(localhost:3306)/registrationandlogin")
	api.MakeUpdaterAPI(r, updateUser, "/api/updateuser", "user", "session", model.User{}, model.Session{})
	findUser := controller.NewMySQLCRUDController("hamid:@tcp(localhost:3306)/registrationandlogin")
	api.MakeRetrieverAPI(r, findUser, "/api/finduser", "user", "session", model.User{}, model.Session{})
	mailboxCtrlr := controller.NewGrpcServiceController("hamid:@tcp(localhost:3306)/registrationandlogin")
	api.MakeServiceAPI(r, mailboxCtrlr, "/api/mailbox", "mailbox", "session", model.Session{})
	go func() {
		err := http.ListenAndServe("localhost:8000", r)
		if err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("Server is now ready to take your orders!")
	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt)

	<-e

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting Down!")
	os.Exit(1)
}
