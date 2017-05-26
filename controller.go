package main

import (
	"encoding/json"
	"fmt"
	"log"

	routing "github.com/go-ozzo/ozzo-routing"

	r "gopkg.in/gorethink/gorethink.v3"
)

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
	fmt.Println(c.Get("total"))
	total := c.Get("total").(int)
	ret := jsonReturnArray{rows, true, total}

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
	ret := jsonReturn{rows, true}
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
	ret := jsonReturn{rows, true}
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
	ret := jsonReturn{rows, true}
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
	ret := jsonReturn{rows, true}
	json, _ := json.Marshal(ret)

	return c.Write(string(json))

}
