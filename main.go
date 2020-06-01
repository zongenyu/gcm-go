// main.go
package main

import (
	"bitbucket.org/raylios/cloudpost-go/slog"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

type info struct {
	Port string
}

var Dat *info
var GCM_DATA []byte

const GCM_KEY_PATH = "./gcmKeys.json"

func init() {

	Dat = new(info)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	if !readGcmKeyFile() {
		slog.Err("Failed to parse gcm key file!")
		return
	}

	parseName()
	parseFlags()

	http.HandleFunc("/probe", probeHandler)
	http.HandleFunc("/gcm", gcmHandler)

	err := http.ListenAndServe(Dat.Port, nil)
	if err != nil {
		slog.Crit("ListenAndServe: ", err)
	}
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warning(">>> %v %v, %v", r.Method, r.URL.Path, r.RemoteAddr)
	io.WriteString(w, "")
}

func gcmHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warning(">>> %v %v, %v", r.Method, r.URL.Path, r.RemoteAddr)

	r.ParseForm()

	token := r.FormValue("token")
	if token == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		slog.Err("Error token: %v %v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	data := r.FormValue("data")
	if data == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		slog.Err("data id: %v %v", http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err := send(token, data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Err("gcm send: %v %v", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return

	}

	io.WriteString(w, "")
}

func readGcmKeyFile() bool {

	bytes, err := ioutil.ReadFile(GCM_KEY_PATH)
	if len(bytes) <= 0 || err != nil {
		slog.Err("Error while reading %v with error: %v", GCM_KEY_PATH, err)
		return false
	}

	GCM_DATA = make([]byte, len(bytes))
	copy(GCM_DATA, bytes)
	//slog.Debug("Read GCM Keys: %v", string(GCM_DATA))

	return true
}

func parseName() {

	s := fmt.Sprintln(os.Args)
	slog.Info("command-line: %v", s)
}

func parseFlags() {

	port := flag.String("p", ":8098", "port")

	flag.Parse()

	Dat.Port = *port
}
