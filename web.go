package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./web/upload.gtpl")
		fc := g.getFileContent()
		t.Execute(w, fc)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./upload"); os.IsNotExist(err) {
		err = os.MkdirAll("upload", 0755)
		if err != nil {
			panic(err)
		}
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	if !checkFile("./upload/" + handler.Filename) {
		http.Redirect(w, r, "/badfile", http.StatusSeeOther)
	} else {
		g.importFile("./upload/" + handler.Filename)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func badfile(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./web/badfile.gtpl")
	t.Execute(w, nil)
}

func path(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	from, ok := r.URL.Query()["from"]
	if !ok || len(from[0]) < 1 {
		log.Println("Url параметр 'from' отсутствует")
		json.NewEncoder(w).Encode([]string{"Url параметр 'from' отсутствует"})
		return
	}
	to, ok := r.URL.Query()["to"]
	if !ok || len(to[0]) < 1 {
		log.Println("Url параметр 'to' отсутствует")
		json.NewEncoder(w).Encode([]string{"Url параметр 'to' отсутствует"})
		return
	}

	result, err := g.extractPath(from[0], to[0])
	if err != nil {
		log.Printf("from: %s, to:%s, result: %s\n", from[0], to[0], err)
		json.NewEncoder(w).Encode([]string{err.Error()})
	} else {
		log.Printf("from: %s, to:%s, result: %s\n", from[0], to[0], result)
		json.NewEncoder(w).Encode(result)
	}
}
