package main

import (
  "encoding/json"
  "fmt"
  zmq "github.com/alecthomas/gozmq"
  "github.com/altonymous/api.6fusion.com/models"
  riak "github.com/tpjg/goriakpbc"
)

func main() {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind("tcp://127.0.0.1:5000")
  socket.Bind("tcp://127.0.0.1:6000")

  for {
    msg, _ := socket.Recv(0)

    var machine models.Machine
    err := json.Unmarshal(msg, &machine)

    machineJson, err := json.Marshal(machine)
    if err != nil {
      fmt.Println("error:", err)
    }

    // println("Got", string(msg))
    storeData(machine.UUID, machineJson)
    socket.Send(msg, 0)
  }
}

func storeData(machineUUID string, machineJson []byte) {
  err := riak.ConnectClient("127.0.0.1:8087")
  if err != nil {
    fmt.Println("Cannot connect, is Riak running?")
    return
  }

  bucket, _ := riak.NewBucket("6fusion")
  obj := bucket.NewObject(machineUUID)
  obj.ContentType = "application/json"

  obj.Data = []byte(machineJson)
  obj.Store()

  fmt.Printf("Stored an object in Riak, vclock = %v\n", obj.Vclock)

  riak.Close()

}
