# redis-protocol
Redis protocol parser in golang

Document
--------
- [GoDoc](http://godoc.org/github.com/quorzz/redis-protocol)

Getting Started
------
```
package main

import (
    "bufio"
    "fmt"
    "github.com/quorzz/redis-protocol"
    "net"
    "time"
)

// Just for example, ignore all the possible errors
func main() {

    conn, _ := net.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)

    raw, _ := protocol.PackCommand("hmset", "users:1", "name", "lily")
    conn.Write(raw)

    reader := bufio.NewReader(conn)
    message, _ := protocol.UnpackFromReader(reader)

    fmt.Println(message.String()) // OK <nil>
    fmt.Println(message.Bool())   // true <nil>



    // There is also anthor usage more flexibily :

    protocol.NewWriter(conn).WriteCommand("hgetall", "users:1")

    msg, _ := protocol.NewReader(conn).ReadMessage()

    fmt.Println(msg.StringMap()) // map[name:lily] <nil>

    fmt.Println(msg.Strings()) // ["name", "lily"] <nil>



    // If you want to use slice or map as parameters :

    args := protocol.NormalizeArgs(map[string]string{"name": "bob"})

    protocol.PackCommand("hmset", "users:1", args)


    args2 := protocol.NormalizeArgs("lpush", "list:users", []string{"no.1", "no.2", "no,3"})

    protocol.NewWriter(conn).WriteCommand(args2)

}
```

Some useful methods of `protocol.Message`:

- `Message.Bool()`
- `Message.Int()`
- `Messsage.Int64()`
- `Message.String()`
- `Message.Strings()`
- `Message.Bytes()`
- `Message.StringMap()`
- `Some more soon ......`

License
-------

[The MIT License (MIT) Copyright (c) 2016 quorzz](http://opensource.org/licenses/MIT)
