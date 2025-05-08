package router

import (
	"fmt"
	"net/http"
)

type handlers struct{
	get func(w http.ResponseWriter,r *http.Request)
	post func(w http.ResponseWriter,r *http.Request)
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

func wrapper(path string,methods handlers) func (w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		fmt.Printf("Request: %s , method: %s\n",path,r.Method)
		if _,ok :=routes[r.URL.Path];!ok {
			http.ServeFile(w,r,"views\\404.html")
			return
		} 
		if r.Method == "GET"{
			if methods.get != nil{
				methods.get(w,r)
				return
			}else{
				fmt.Fprintf(w,"Invalid route")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}else if r.Method == "POST"{
			if methods.post != nil{
				methods.post(w,r)
				return
			}else{
				fmt.Fprintf(w,"Invalid route")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}else if r.Method == http.MethodOptions {
			methods.post(w,r)
			return
		}else{
			fmt.Println("In invalid route",r.Method)
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