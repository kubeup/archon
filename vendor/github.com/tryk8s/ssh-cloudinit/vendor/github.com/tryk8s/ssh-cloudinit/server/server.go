package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tryk8s/ssh-cloudinit/client"
	"io/ioutil"
	"net/http"
	"os"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "OK")
		return
	}

	taskid := r.URL.Query().Get("id")
	if taskid == "" {
		fmt.Fprintf(w, "OK")
		return
	}
	conf := &client.Config{
		Port:   22,
		Os:     "ubuntu",
		Stdout: os.Stdout,
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		err := client.Run(conf)
		if conf.Callback != "" {
			var result string
			if err == nil {
				result = "success"
			} else {
				result = "error"
			}
			response := struct {
				id     string
				result string
			}{taskid, result}
			body, err := json.Marshal(response)
			if err != nil {
				return
			}
			http.Post(conf.Callback, "application/json", bytes.NewBuffer(body))
		}
	}()

	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/run", runHandler)
	http.ListenAndServe(":8080", nil)
}
