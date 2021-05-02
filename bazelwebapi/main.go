package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	c, err := loadconfig()
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		c: c,
	}

	tmpl, err := template.New("dockerconf").Parse(dockerconftmpl)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(c.DockerConfPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(f, c)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	mux := http.NewServeMux()
	// should publish an image. Target required
	mux.HandleFunc("/execute/task", s.ExecuteTask)

	log.Println("WebAPI running on port 8088")
	log.Fatal(http.ListenAndServe(":8088", mux))
}

const dockerconftmpl = `
{
	"auths": {
		"{{.Registry}}": {
			"auth": "{{.Auth}}"
		}
	}
}
`
