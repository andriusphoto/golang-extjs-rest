package controllers

import (
	config "api/config"
	lib "api/lib"
	"encoding/json"
	"log"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/auth"
	r "gopkg.in/gorethink/gorethink.v3"
)

func ImportRoutes(api routing.RouteGroup) {

	api.Get("/restricted", auth.JWT(config.SigningKey), lib.GetJWTclaims, restricted)
	api.Get("/<table>", lib.UseTable, lib.AddFilter, lib.Total, lib.AddSorter, lib.AddPagination, Get)
	api.Get("/<table>/<id>", auth.JWT(config.SigningKey), lib.UseTable, GetOne)
	api.Delete("/<table>/<id>", auth.JWT(config.SigningKey), lib.UseTable, Delete)
	api.Post("/<table>", auth.JWT(config.SigningKey), lib.UseTable, Create)
	api.Put("/<table>/<id>", auth.JWT(config.SigningKey), lib.UseTable, Update)
}
func restricted(c *routing.Context) error {

	claims := c.Get("claims")
	ret := lib.JsonReturn{claims, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
func Get(c *routing.Context) error {
	q := c.Get("q").(r.Term)

	res, err := q.Run(c.Get("session").(r.QueryExecutor))
	if err != nil {
		log.Fatalln(err)
	}

	var rows []interface{}
	err = res.All(&rows)
	if err != nil {
		log.Fatalln(err)
	}

	total := c.Get("total").(int)

	ret := lib.JsonReturnArray{rows, true, total}

	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
func GetOne(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	res, err := q.Get(c.Param("id")).Run(c.Get("session").(r.QueryExecutor))
	if err != nil {
		log.Fatalln(err)
	}

	var rows []interface{}
	err = res.All(&rows)
	if err != nil {
		// error
	}
	ret := lib.JsonReturn{rows, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
func Delete(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	_, err := q.Get(c.Param("id")).Delete().Run(c.Get("session").(r.QueryExecutor))
	if err != nil {
		log.Fatalln(err)
	}

	var rows []interface{}
	ret := lib.JsonReturn{rows, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
func Create(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	var data interface{}
	err := c.Read(&data)
	if err != nil {
		log.Fatalln(err)
	}
	ins, err1 := q.Insert(data).RunWrite(c.Get("session").(r.QueryExecutor))
	if err1 != nil {
		log.Fatalln(err1)
	}

	res, err2 := q.Get(ins.GeneratedKeys[0]).Run(c.Get("session").(r.QueryExecutor))
	if err2 != nil {
		log.Fatalln(err2)
	}

	var rows []interface{}
	err3 := res.All(&rows)
	if err3 != nil {
		// error
	}
	ret := lib.JsonReturn{rows, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
func Update(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	var data interface{}
	err := c.Read(&data)
	if err != nil {
		log.Fatalln(err)
	}
	_, err1 := q.Get(c.Param("id")).Update(data).RunWrite(c.Get("session").(r.QueryExecutor))
	if err1 != nil {
		log.Fatalln(err1)
	}

	res, err2 := q.Get(c.Param("id")).Run(c.Get("session").(r.QueryExecutor))
	if err2 != nil {
		log.Fatalln(err2)
	}

	var rows []interface{}
	err3 := res.All(&rows)
	if err3 != nil {
		// error
	}
	ret := lib.JsonReturn{rows, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
