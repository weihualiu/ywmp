
package channelCore

import "fmt"
import "io"
import "net/http"

type UserHelloController struct {
    
}

func (uc *UserHelloController) Process(w http.ResponseWriter, req *http.Request) {
    fmt.Println("userHelloController channel")
    io.WriteString(w, req.RequestURI)
}