package main

import (
	"api/config"
	"api/controllers"
	"api/lib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/auth"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	r "gopkg.in/gorethink/gorethink.v3"
)

func main() {

	router := routing.New()

	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		lib.Connect,
		lib.UseDB,
		// func(c *routing.Context) error {
		// 	// id, err := authenticate(c)
		// 	// if err != nil {
		// 	// 	return err
		// 	// }
		// 	username, password := parseBasicAuth(c.Request.Header.Get("Authorization"))
		// 	if username == "demo" && password == "foo" {
		// 		// auth.JWT(signingKey)
		// 	}
		// 	token, err := auth.NewJWT(jwt.MapClaims{
		// 		"id": "10000",
		// 	}, signingKey)
		// 	if err != nil {

		// 		return err

		// 	}

		// 	return c.Write(token)

		// },
	)
	api := router.Group("/api")
	api.Post("/login", func(c *routing.Context) error {
		// id, err := authenticate(c)
		// if err != nil {
		// 	return err
		// }
		var ret interface{}
		username, password := lib.ParseBasicAuth(c.Request.Header.Get("Authorization"))
		q := c.Get("q").(r.Term)

		res, err := q.Table("users").Filter(map[string]interface{}{
			"username": username,
			"password": password,
		}).Run(c.Get("session").(r.QueryExecutor))
		if err != nil {
			fmt.Print(err)

			return nil
		}
		var users []map[string]interface{}
		err = res.All(&users)
		if len(users) > 0 {
			token, err := auth.NewJWT(jwt.MapClaims{
				"username": users[0]["username"],
			}, config.SigningKey)
			if err != nil {
				fmt.Print(err)

				return nil
			}
			ret = lib.JsonReturn{token, true}
		} else {
			ret = lib.JsonReturn{nil, false}
		}

		json, _ := json.Marshal(ret)

		return c.Write(string(json))
	})

	// api.Options("/<table>", conect, useTable, Get)
	controllers.ImportRoutes(*api)

	print("Server started... port:7776")
	http.Handle("/", router)
	http.ListenAndServe(":7776", nil)

}
