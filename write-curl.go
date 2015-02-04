package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"net/http"
	
	"github.com/writeas/writeas-telnet/store"
)

var (
	outDir string
	indexPage []byte
)

func poster(w http.ResponseWriter, r *http.Request) {
	post := r.FormValue("w")

	if post == "" {
		fmt.Fprintf(w, "%s", indexPage)
		return
	}

	filename, err := store.SavePost(outDir, []byte(post))
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "Couldn't save :(\n")
		return
	}
	fmt.Fprintf(w, "https://write.as/%s\n", filename)
}

func main() {
	outDirPtr := flag.String("o", "/home/matt", "Directory where text files will be stored.")
	staticDirPtr := flag.String("s", "./static", "Directory where required static files exist.")
	portPtr := flag.Int("p", 8080, "Port to listen on.")
	flag.Parse()

	outDir = *outDirPtr

	fmt.Print("Initializing...")
	var err error
	indexPage, err = ioutil.ReadFile(*staticDirPtr + "/cmd.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DONE")

	fmt.Printf("Serving on http://localhost:%d\n", *portPtr)

	http.HandleFunc("/", poster)
	http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), nil)
}
