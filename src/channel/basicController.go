
package channel

import "fmt"
import "io"
import "net/http"

/*
 * 
 */
type indexController struct {
}

func (ic *indexController) Process(w http.ResponseWriter, req *http.Request) {
    fmt.Println("indexController channel")
    io.WriteString(w, req.RequestURI)
}

type heartbeatController struct {
    
}

func (ic *heartbeatController) Process(w http.ResponseWriter, req *http.Request) {
    fmt.Println("heartbeatController channel")
    io.WriteString(w, req.RequestURI)
}