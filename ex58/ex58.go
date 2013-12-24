package main

import (
	"fmt"
	"net/http"
)

type String string

type Struct struct{
	Greeting string
	Punct string
	Who string
}

func (h String) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprint( w, h );
}

func (s Struct) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"%s %s [%s]\n", s.Greeting, s.Punct, s.Who )	
}

type DefaultHandler struct{}
func (dh DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Nothing Interesting here!!!\n")
}


func main(){
	http.Handle( "/string", String("I am a pizza!") )
	http.Handle( "/struct", &Struct{"Hello",":","Lambros"} )
	var df DefaultHandler
	http.Handle( "/", df )
	http.ListenAndServe("localhost:4000", nil)
}


