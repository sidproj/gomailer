package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type handlers struct{
	get func(w http.ResponseWriter,r *http.Request)
	post func(w http.ResponseWriter,r *http.Request)
	put func(w http.ResponseWriter,r *http.Request)
	delete func(w http.ResponseWriter,r *http.Request)
}

var routes = map[string]handlers{}

type ctxKey string

const pathParamsKey ctxKey = "pathParams"

func GetPathParams(r *http.Request)map[string]string{
	return r.Context().Value(pathParamsKey).(map[string]string)
}

type RouterNode struct{
	handlers handlers 
	path string
	childRouters map[string]*RouterNode 
	// the about member's key for the map can be param. to know if it is dynamic route, the node's path would be empty
}

func Get(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.get = handler
	routes[path] = h
	rootNode.AddChildRouters(path,"GET",handler)
}

func Post(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.post = handler
	routes[path] = h
	rootNode.AddChildRouters(path,"POST",handler)
}

func Put(path string,handler func(w http.ResponseWriter,r *http.Request)){
	h:=routes[path]
	h.put = handler
	routes[path] = h
	rootNode.AddChildRouters(path,"PUT",handler)
}

func Delete(path string,handler func(w http.ResponseWriter, r *http.Request)){
	h:=routes[path]
	h.delete = handler
	routes[path] = h
	rootNode.AddChildRouters(path,"DELETE",handler)
}

func (n* RouterNode)findHandler(path string)(*RouterNode, map[string]string){
	pathSlice := strings.Split(path, "/")[1:]
	var finalNode * RouterNode
	params := map[string]string{}
	for _,pathSeg := range pathSlice{
		var curNode * RouterNode = nil
		if(len(n.childRouters) == 0){
			return nil,nil
		}
		for k,v:= range n.childRouters{
			if k == pathSeg{
				curNode = v
				n = v
				break
			}
		}
		if curNode == nil {
			// check if pathSeg is here
			for k,v:= range n.childRouters{
				if strings.HasPrefix(k,":") {
					curNode = v
					params[k] = pathSeg
					n = v
					break
				}
			}
			if curNode == nil{
				return curNode,nil
			}
		}
		finalNode = curNode
	}
	return finalNode,params;
}

func (r * RouterNode)AddChildRouters(path string,method string,handler func(w http.ResponseWriter,r *http.Request)){
	if path == "/"{
		switch method{
			case "GET": r.handlers.get = handler
			case "POST": r.handlers.post = handler
			case "PUT": r.handlers.put = handler
			case "DELETE": r.handlers.delete = handler
		}
		return
	}
	pathSlices := strings.Split(path,"/")[1:]
	var travelNode = r
	for _,pathSeg := range pathSlices{
		if _,ok := travelNode.childRouters[pathSeg];!ok{
			travelNode.childRouters[pathSeg] = &RouterNode{
				handlers:handlers{},
				path:pathSeg,
				childRouters: make(map[string]*RouterNode),
			}
		}
		travelNode = travelNode.childRouters[pathSeg]
	}

	switch method{
		case "GET": travelNode.handlers.get = handler
		case "POST": travelNode.handlers.post = handler
		case "PUT": travelNode.handlers.put = handler
		case "DELETE": travelNode.handlers.delete = handler
	}
}

func (r*RouterNode)Display(){
	if r == nil {
		return
	}
	
	fmt.Printf("%+v \n\n",*r)
	for _,v:= range r.childRouters{
		v.Display()
	}
}

var rootNode = RouterNode{
	path: "/",
	handlers: handlers{},
    childRouters: make(map[string]*RouterNode),
}

// provides mapping for all the requests
func wrapper() func (w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		if r.URL.Path == ""{
			r.URL.Path = "/"
		}
		fmt.Printf("Request: %s , method: %s\n",r.URL.Path,r.Method)
		
		node := &rootNode
		if r.URL.Path!="/"{
			fmt.Println("Started finding node...")
			findNode,params := rootNode.findHandler(r.URL.Path)
		
			if(findNode==nil){
				http.ServeFile(w,r,"views\\404.html")
				return	
			}
			node = findNode
			if(params != nil){
				ctx := context.WithValue(r.Context(),pathParamsKey,params)
				r = r.WithContext(ctx)
			}
		}

		switch r.Method{
			case http.MethodOptions:
				requestedMethod := r.Header.Get("Access-Control-Request-Method")
				switch requestedMethod{
					case http.MethodPost: node.handlers.post(w,r)
					case http.MethodPut: node.handlers.put(w,r)
					case http.MethodDelete: node.handlers.delete(w,r)
					case http.MethodGet: node.handlers.get(w,r)
				}
				return
			case http.MethodGet:
				if node.handlers.get != nil{
					node.handlers.get(w,r)
					return
				}
			case http.MethodPost:
				if node.handlers.post != nil{
					node.handlers.post(w,r)
					return
				}
			case http.MethodPut:
				if node.handlers.put != nil{
					node.handlers.put(w,r)
					return
				}
			case http.MethodDelete:
				if node.handlers.delete != nil{
					node.handlers.delete(w,r)
					return
				}
			default:
				fmt.Fprintf(w,"Invalid route")
				w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func LoadRoutes(){
	for k := range routes{
		http.HandleFunc(k,wrapper())
	}
}