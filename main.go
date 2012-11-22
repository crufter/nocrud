package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	gt "github.com/opesun/gotrigga"
	"github.com/opesun/nocrud/frame/bang"
	"github.com/opesun/nocrud/frame/config"
	"github.com/opesun/nocrud/frame/top"
	"io"
	"labix.org/v2/mgo"
	"net/http"
)

func main() {
	defer mainErr()
	config := config.New()
	config.LoadFromFile()
	dial := config.DBAddr
	if len(config.DBUser) != 0 || len(config.DBPass) != 0 {
		if len(config.DBUser) == 0 {
			panic("Database password provided but username is missing.")
		}
		if len(config.DBPass) == 0 {
			panic("Database username is provided but password is missing.")
		}
		dial = config.DBUser + ":" + config.DBPass + "@" + config.DBAddr
		if !config.DBAdmMode {
			dial = dial + "/" + config.DBName
		}
	}
	fmt.Println("Connecting to db server.")
	session, err := mgo.Dial(dial)
	if err != nil {
		panic(err)
	}
	var conn *gt.Connection
	if len(config.MsgAddr) > 0 {
		fmt.Println("Connecting to msg server.")
		conn, err = gt.Connect(config.MsgAddr)
		if err != nil {
			panic(err)
		}
	}
	db := session.DB(config.DBName)
	defer session.Close()
	fmt.Println("Waiting for websockets connections.")
	http.HandleFunc("/ws/", func(w http.ResponseWriter, req *http.Request) {
		defer printErr(w)
		wsHandler := func(ws *websocket.Conn) {
			defer ws.Close()
			req.URL.Path = "/" + req.URL.Path[3:]
			ctx, err := bang.NewWS(conn, session, db, w, req, config, ws)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			top := top.New(ctx)
			err = top.RouteWS()
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
		}
		websocket.Handler(wsHandler).ServeHTTP(w, req)
	})
	fmt.Println("Waiting http connections.")
	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			defer printErr(w)
			ctx, err := bang.New(conn, session, db, w, req, config)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			if ctx == nil {
				return
			}
			top := top.New(ctx)
			err = top.Route()
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
		})
	err = http.ListenAndServe(config.Addr+":"+config.PortNum, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func mainErr() {
	if r := recover(); r != nil {
		fmt.Println(fmt.Sprint(r))
	}
}

func printErr(w http.ResponseWriter) {
	if r := recover(); r != nil {
		io.WriteString(w, fmt.Sprint(r))
		panic(fmt.Sprint(r))
	}
}
