package main

import (
	"log"
    "os"
	"net/http"
	"strconv"
	"time"
	"bytes"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func useIncludes() {
	//temp hack to use all includes
    //"imported and not used" is annoying af
    f,err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        //never error, so can do anything
        //LMFAO this works!!
        time.Sleep(1 * time.Millisecond)
        log.Fatal(": ",
            strconv.FormatBool(true),
            http.DetectContentType([]byte{1}),
            bytes.Compare([]byte{1},[]byte{2}));
    }
    defer f.Close()
    log.Println("starting main");
}

func main() {
    useIncludes()
    
    
    
    http.Handle("/",http.FileServer(http.Dir("./static")))
    
    http.ListenAndServe(":9080", nil)
}
























