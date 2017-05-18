package main

import (
	"log"
	"net/http"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	r "gopkg.in/gorethink/gorethink.v3"
)

type jsonReturnArray struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
}

func main() {

	router := routing.New()

	//options := cors.Options{AllowOrigins: "http://dev.mro.flts.local", AllowCredentials: true, AllowHeaders: "*"}
	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)
	api := router.Group("/api")
	// api.Options("/<table>", conect, useTable, Get)

	api.Get("/<table>", conect, useTable, Get)
	api.Get("/<table>/<id>", conect, useTable, GetOne)
	api.Delete("/<table>/<id>", conect, useTable, Delete)
	api.Post("/<table>", conect, useTable, Create)
	api.Put("/<table>/<id>", conect, useTable, Update)
	print("Server started... port:7776")
	http.Handle("/", router)
	http.ListenAndServe(":7776", nil)

}
func conect(c *routing.Context) error {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost",
	})
	c.Set("session", session)
	if err != nil {

		log.Fatalln(err)
	}
	return nil
}
func useTable(c *routing.Context) error {
	c.Set("q", r.DB("test").Table(c.Param("table")))

	return nil
}
