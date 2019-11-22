package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var MIRRORS = []string {"https://ftp.fau.de/gentoo", "https://ftp-stud.hs-esslingen.de/pub/Mirrors/gentoo",
					 "http://ftp.fi.muni.cz/pub/linux/gentoo", "http://gentoo.mirror.web4u.cz",
					 "http://gentoo.mirror.web4u.cz/", "http://gentoo.modulix.net/gentoo",
					 "http://ftp-stud.hs-esslingen.de/pub/Mirrors/gentoo",
					 "https://mirror.eu.oneandone.net/linux/distributions/gentoo/gentoo",
					 "https://mirror.netcologne.de/gentoo/",
					 "https://ftp.halifax.rwth-aachen.de/gentoo/",
					 "http://ftp.ntua.gr/pub/linux/gentoo/",
					 "https://mirrors.evowise.com/gentoo/",
					 "https://ftp.snt.utwente.nl/pub/os/linux/gentoo/",
					 "https://mirror.leaseweb.com/gentoo/",
					 "http://ftp.vectranet.pl/gentoo/",
					 "http://ftp.dei.uc.pt/pub/linux/gentoo/",
					 "https://gentoo.wheel.sk/",
					 "http://tux.rainside.sk/gentoo/",
					 "https://mirror.bytemark.co.uk/gentoo/",
					 "http://mirror.isoc.org.il/pub/gentoo/",
}

func findMirror(requestURI string, w http.ResponseWriter, r *http.Request) {
	for _, mirrorURL := range MIRRORS {
		url := strings.TrimRight(mirrorURL, "/") + requestURI
		res, err := http.Head(url)
		//defer res.Body.Close()

		if err != nil {
			log.Println(err)
			continue
		}
		if res.StatusCode == http.StatusOK {
			log.Printf("Requested URL %s found at %s", requestURI, url)
			http.Redirect(w, r, mirrorURL + requestURI, http.StatusTemporaryRedirect)
			return
		}

	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 page not found")
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello index")
}
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		serveRoot(w, r)
		return
	}

	findMirror(r.RequestURI, w, r)
}


func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(catchAllHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
