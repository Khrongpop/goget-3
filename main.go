package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"_topic"`
	Done  bool          `json:"done" bson:"_done"`
}

type handler struct {
	m *mgo.Session
}

func main() {

	e := echo.New()

	// yml
	// mongo:

	// env
	// MONGO_HOST
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	// port := ":" + viper.GetString("port")

	// mongoHost := "13.250.119.252"
	// mongoUser := "root"
	// mongoPass := "example"
	// port := ":1323"

	fmt.Println(mongoHost, mongoUser, mongoPass)
	// session, err := mgo.Dial("root:example@13.250.119.252:27017"
	connString := fmt.Sprintf("%v:%v@%v", mongoUser, mongoPass, mongoHost)
	session, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	e.Use(middleware.Logger())
	e.GET("/todo", h.list)
	e.GET("/todo/:id", h.view)
	e.POST("/todo", h.create)
	e.PUT("/todo/:id", h.done)
	e.DELETE("/todo/:id", h.delete)
	e.Logger.Fatal(e.Start(":1323"))
}

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()

	col := session.DB("workshop").C("todos_su")

	if err := col.Insert(t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	var ts []todo
	col := session.DB("workshop").C("todos_su")

	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	col := session.DB("workshop").C("todos_su")

	var t todo
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) done(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	col := session.DB("workshop").C("todos_su")

	var t todo
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	t.Topic = t.Topic + " update : "
	if err := col.UpdateId(id, t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	col := session.DB("workshop").C("todos_su")

	if err := col.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}
