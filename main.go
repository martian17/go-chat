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


//message "class" (class isn't really a thing in go but yk)
type Message struct {
    str string
    sender *Client
}

func newMessage(str string, client *Client) *Message {
    return &Message{
       str:    str,
       sender: client,
   }
}




// client "class" (class isn't really a thing in go but yk)
type Client struct {
    //hub, doubly linked
    hub *Hub
    
    //connection pointer
    conn *websocket.Conn
    
    //internal: reader -> client
    handlePost chan *Message
    //send(msg) pseudo method
    send chan *Message
    
}

// Client constructor
func newClient(hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        hub:        hub,
        conn:       conn,
        handlePost: make(chan *Message),
        send:       make(chan *Message),
    }
}

//methods

//will be ran under goroutine
func (this *Client) reader(){
    for{
        msgType, bytes, err := this.conn.ReadMessage()
        //msgType either websocket.TextMessage or websocket.BinaryMessage
        //bytes are of type []byte
        
        if err != nil {
            log.Println("Read Failed: ",err);
            break;
        }
        if msgType != websocket.TextMessage{
            continue;
        }
        
        this.handlePost <- newMessage(string(bytes),this)
    }
    //cleanup
    close(this.handlePost);
    //could be a double close, so may cause error
    //docs doesn't say anything about it so idk what would happen
    this.conn.Close();
}


func (this *Client) run(){
    this.hub.Register <- this;
    go this.reader()
    for{
        select{
        case msg, ok := <- this.send:
            if !ok {
                break
            }
            if msg.sender == this {
                continue
            }
            err := this.conn.WriteMessage(websocket.TextMessage, []byte(msg.str))
            if err != nil {
                log.Println("write failed:", err)
                break
            }
        case msg, ok := <- this.handlePost:
            //received message from client
            if !ok {
                break
            }
            this.hub.broadcast <- msg
        }
    }
    //cleanup
    this.hub.Unregister <- this
    this.conn.Close()
}





// hub "class"
type Hub struct {
    // Registered clients.
    clients map[*Client]bool

    // Inbound messages from the clients.
    // use <- to read from it in each goroutine
    broadcast chan *Message
    
    Register chan *Client
    Unregister chan *Client
}

// hub constructor
func newHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan *Message),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}

func (this *Hub) run() {
    for {
        select{
        case client := <- this.Register:
            this.clients[client] = true
        case client := <- this.Unregister:
            delete(this.clients,client)
        case msg := <- this.broadcast:
            for client := range this.clients {
                client.send <- msg
            }
        }
    }
}










func main() {
    useIncludes()
    
    hub := newHub()
    go hub.run()
    
    //runs for each client
    http.HandleFunc("/socket", func (w http.ResponseWriter, r *http.Request) {
        // Upgrade upgrades the HTTP server connection to the WebSocket protocol.
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal("upgrade failed: ", err)
        }
        client := newClient(hub,conn)
        go client.run()
    })
    
    
    
    http.Handle("/",http.FileServer(http.Dir("./static")))
    
    http.ListenAndServe(":9080", nil)
}
























