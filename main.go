package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

var keyValueStore map[string]string
var vStoreMutex sync.RWMutex

func init() {
	keyValueStore = make(map[string]string)
	vStoreMutex = sync.RWMutex{}
}

func main() {

	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	http.HandleFunc("/remove", remove)
	http.HandleFunc("/list", list)
	http.ListenAndServe(":3000", nil)

}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}
		if len(values.Get("key")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", "Wrong input key.")
			return
		}
		vStoreMutex.RLock()
		value := keyValueStore[string(values.Get("key"))]
		fmt.Printf("value:%v", value)
		vStoreMutex.RUnlock()
		fmt.Fprint(w, value)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only GET accepted.")
	}
}

//post
func set(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}

		//k
		if len(values.Get("key")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", "Wrong input key.")
			return
		}

		//v
		if len(values.Get("value")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", "Wrong input value.")
			return
		}

		vStoreMutex.Lock()
		keyValueStore[string(values.Get("key"))] = string(values.Get("value"))
		vStoreMutex.Unlock()
		fmt.Fprint(w, "success")

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only POST accepted.")
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", err)
			return
		}

		//k
		if len(values.Get("key")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error:", "Wrong input key.")
			return
		}
		vStoreMutex.Lock()
		delete(keyValueStore, values.Get("value"))
		vStoreMutex.Unlock()
		fmt.Fprint(w, "success")

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only GET accepted.")
	}

}

func list(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vStoreMutex.RLock()
		for k, v := range keyValueStore {
			fmt.Fprintln(w, k, ":", v)
		}
		vStoreMutex.Unlock()

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Only GET accepted.")
	}
}
