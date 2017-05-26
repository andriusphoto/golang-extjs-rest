package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	r "gopkg.in/gorethink/gorethink.v3"
)

type jsonReturnArray struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
	Length  int           `json:"total"`
}
type jsonReturn struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type filter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
type sorter struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}

func main() {

	router := routing.New()

	//options := cors.Options{AllowOrigins: "http://dev.mro.flts.local", AllowCredentials: true, AllowHeaders: "*"}
	router.Use(
		// all these handlers are shared by every route
		// access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		Connect,
		UseTable,
	)
	api := router.Group("/api")

	// api.Options("/<table>", conect, useTable, Get)

	api.Get("/<table>", AddFilter, Total, AddSorter, AddPagination, Get)
	api.Get("/<table>/<id>", GetOne)
	api.Delete("/<table>/<id>", Delete)
	api.Post("/<table>", Create)
	api.Put("/<table>/<id>", Update)
	print("Server started... port:7776")
	http.Handle("/", router)
	http.ListenAndServe(":7776", nil)

}
func Connect(c *routing.Context) error {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost",
	})
	c.Set("session", session)
	if err != nil {

		log.Fatalln(err)
	}
	return nil
}
func UseTable(c *routing.Context) error {
	c.Set("q", r.DB("test").Table(c.Param("table")))

	return nil
}
func Total(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	var data interface{}
	_, err := q.Count(&data).Run(c.Get("session").(r.QueryExecutor))
	if err != nil {
		log.Fatalln(err)
	}
	c.Set("total", data)
	return nil
}
func AddSorter(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	sortstr := c.Request.FormValue("sort")
	sorter := []sorter{}
	json.Unmarshal([]byte(sortstr), &sorter)

	for _, item := range sorter {
		if item.Direction == "ASC" {
			q = q.OrderBy(item.Property)

		}
	}
	c.Set("q", q)
	return nil
}
func AddPagination(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {

		log.Fatalln(err)
	}
	start, err := strconv.Atoi(c.Request.FormValue("start"))
	if err != nil {

		log.Fatalln(err)
	}
	limit, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err != nil {

		log.Fatalln(err)
	}
	fmt.Println(limit)
	sortstr := c.Request.FormValue("sort")
	sorter := []sorter{}
	json.Unmarshal([]byte(sortstr), &sorter)

	q = q.Slice(start, limit*page)

	c.Set("q", q)
	return nil
}
func AddFilter(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	filterstr := c.Request.FormValue("filter")
	filter := []filter{}
	json.Unmarshal([]byte(filterstr), &filter)
	for _, item := range filter {
		if item.Operator == "like" {
			q = q.Filter(func(el r.Term) r.Term {
				return el.Field(item.Property).Match(item.Value)
			})
		}

	}
	c.Set("q", q)
	return nil
}
