package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/sko00o/go-lab/basic/context/go-blog-demo/google"
	"github.com/sko00o/go-lab/basic/context/go-blog-demo/userip"
)

func main() {
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	timeout, err := time.ParseDuration(r.FormValue("timeout"))
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	userIP, err := userip.FromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx = userip.NewContext(ctx, userIP)

	start := time.Now()
	res, err := google.Search(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)

	if err := resultsTemplate.Execute(w, struct {
		Results          google.Results
		Timeout, Elapsed time.Duration
	}{
		Results: res,
		Timeout: timeout,
		Elapsed: elapsed,
	}); err != nil {
		log.Print(err)
		return
	}
}

var resultsTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results}}
    <li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
  {{end}}
  </ol>
  <p>{{len .Results}} results in {{.Elapsed}}; timeout {{.Timeout}}</p>
</body>
</html>
`))
