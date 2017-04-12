
package channel

import "fmt"
import "io"
import "net/http"

type UserController struct {
    
}

func (uc *UserController) Process(w http.ResponseWriter, req *http.Request) {
    fmt.Println("userController channel")
    io.WriteString(w, req.RequestURI)
}
