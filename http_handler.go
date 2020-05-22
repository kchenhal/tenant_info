package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func handleRoot(w http.ResponseWriter, req *http.Request) {
	const tpl = `
<title>Status</title>
<h1>Current Configurations</h1>

<h2>Log Level:   </h2>

<form action="/changeLoglevel">
    <input type="submit" value="Change" />
</form>  

<table class="table  table-striped table-bordered">
    <thead>
        <tr>
            <th>Value</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <td>{{.LogLevel}}</td>
        <td>between 1 to 5, 1 is most verbose, while 5 is least</td>
    </tbody>
</table>
	`

	data := struct {
		LogLevel int64
	}{}

	data.LogLevel = *logLevel

	t, _ := template.New("root").Parse(tpl)
	t.Execute(w, data)
}

func handleChangeLogLevel(w http.ResponseWriter, req *http.Request) {
	const tpl = `
		<h1> Change Log Level</h1>
		<form action="/logLevel" method="POST">
		<label for="logLevel">New log level:</label>
		<input type="number" id="logLevel" name="logLevel" value="{{.}}" min="1" max="5">
		<input type="submit" value="Submit">
		</form>	
	`
	t, _ := template.New("changeLogLevel").Parse(tpl)
	t.Execute(w, *logLevel)
}
func handleUpdateTenantInterval(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Fprintf(w, "%d minutes\n", tenantInterval)
	} else {
		oldLevel := tenantInterval
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "failed to parse interval\n")
			return
		}

		newLevel, err := strconv.ParseInt(req.Form.Get("interval"), 10, 32)
		if err == nil {
			tenantInterval = time.Duration(newLevel) * time.Minute
		} else {
			fmt.Fprintf(w, "failed to convert log level\n")
			return
		}

		fmt.Fprintf(w, "<p>interval changed from %d to %d minutes</p>", oldLevel, tenantInterval)
		fmt.Fprintf(w, "<a href=\"/\">Home</a>")
	}

}

func handleLogLevel(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Fprintf(w, "%d\n", *logLevel)
	} else {
		oldLevel := *logLevel
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "failed to parse logLevel\n")
			return
		}

		newLevel, err := strconv.ParseInt(req.Form.Get("logLevel"), 10, 32)
		if err == nil {
			*logLevel = newLevel
		} else {
			fmt.Fprintf(w, "failed to convert log level\n")
			return
		}

		fmt.Fprintf(w, "<p>log level changed from %d to %d</p>", oldLevel, *logLevel)
		fmt.Fprintf(w, "<a href=\"/\">Home</a>")
	}

}
