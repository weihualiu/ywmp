
package channelCore

import "fmt"
import "io"
import "net/http"

type UserExchangeController struct {
}

func (uec *UserExchangeController)Process(w http.ResponseWriter, req *http.Request) {
    fmt.Println("userExchangeController channel")
    io.WriteString(w, req.RequestURI)
}