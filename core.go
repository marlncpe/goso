package goso

import (
  "net/http"
  "net/url"
  "encoding/json"
  "fmt"
  "os"
)

const (
  GET = "GET"
  POST = "POST"
  PUT = "PUT"
  DELETE = "DELETE"
  )

type Configuration struct {
    Name    string
    Version string
    Port    string
    Status  string
}

type Resource interface {
  Get(values url.Values) (int, interface{})
  Post(values url.Values) (int, interface{})
  Put(values url.Values) (int, interface{})
  Delete(values url.Values) (int, interface{})
}

type (
  GetNotSupported struct {}
  PostNotSupported struct {}
  PutNotSupported struct {}
  DeleteNotSupported struct {}
)


func (GetNotSupported) Get(values url.Values) (int, interface{}) {
  return 405, ""
}

func (PostNotSupported) Post(values url.Values) (int, interface{}) {
  return 405, ""
}

func (PutNotSupported) Put(values url.Values) (int, interface{}) {
  return 405, ""
}

func (DeleteNotSupported) Delete(values url.Values) (int, interface{}) {
  return 405, ""
}

type API struct{}

func (api *API) Abort(rw http.ResponseWriter, statusCode int){
  rw.WriteHeader(statusCode)
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc {
  return func (rw http.ResponseWriter, request *http.Request){
    var data interface{}
    var code int

    request.ParseForm()
    method := request.Method
    values := request.Form

    switch method {
    case GET:
      code, data = resource.Get(values)
    case POST:
      code, data = resource.Post(values)
    case PUT:
      code, data = resource.Put(values)
    case DELETE:
      code, data = resource.Delete(values)
    default:
      api.Abort(rw, 405)
      return
    }

    content, err := json.Marshal(data)
    if err != nil {
      api.Abort(rw, 500)
    }
    rw.WriteHeader(code)
    rw.Write(content)
  }
}

func (api *API) AddResource(resource Resource, path string) {
  http.HandleFunc(path, api.requestHandler(resource))
}

func Run(route string){
  file, _ := os.Open(route)
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  if err != nil {
    fmt.Println("error:", err)
  }
  port := fmt.Sprint(":",configuration.Port)
  fmt.Println(configuration.Name,"funciona en el puerto: ",configuration.Port)
  http.ListenAndServe(port, nil)
}