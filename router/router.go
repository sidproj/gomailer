package router

import (
	"fmt"
	"net/http"
)

type handlers struct{
	get func(w http.ResponseWriter,r *http.Request)
	post func(w http.ResponseWriter,r *http.Request)
	put func(w http.ResponseWriter,r *http.Request)
	delete func(w http.ResponseWriter,r *http.Request)
}

var routes = map[string]handlers{}

func Get(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.get = handler
	routes[path] = h
}

func Post(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.post = handler
	routes[path] = h
}

func Put(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.put = handler
	routes[path] = h
}

func Delete(path string,handler func(w http.ResponseWriter, r *http.Request)){
	h:=routes[path]
	h.delete = handler
	routes[path] = h
}

// provides mapping for all the requests
func wrapper(path string,methods handlers) func (w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		fmt.Printf("Request: %s , method: %s\n",path,r.Method)
		if _,ok :=routes[r.URL.Path];!ok {
			http.ServeFile(w,r,"views\\404.html")
			return
		}

		switch r.Method{
			case http.MethodOptions:
				requestedMethod := r.Header.Get("Access-Control-Request-Method")
				fmt.Println("requested method for pre flight:",requestedMethod)
				switch requestedMethod{
					case http.MethodPost:methods.post(w,r)
					case http.MethodPut:methods.put(w,r)
					case http.MethodDelete:methods.delete(w,r)
					case http.MethodGet:methods.get(w,r)
				}
				return
			case http.MethodGet:
				if methods.get != nil{
					methods.get(w,r)
					return
				}
			case http.MethodPost:
				if methods.post != nil{
					methods.post(w,r)
					return
				}
			case http.MethodPut:
				if methods.put != nil{
					methods.put(w,r)
					return
				}
			case http.MethodDelete:
				if methods.delete != nil{
					methods.delete(w,r)
					return
				}
			default:
				fmt.Fprintf(w,"Invalid route")
				w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func LoadRoutes(){
	for k,v := range routes{
		http.HandleFunc(k,wrapper(k,v))
	}
}