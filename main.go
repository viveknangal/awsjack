// MIT License

// Copyright (c) 2020 Vivek Aggarwal

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This is main package which start the server & invokes handlers
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	awsjack "github.com/viveknangal/awsjack/utils"
)

// This Slice holds list of regions provided via CommandLine arguments
var RegionList []string

// "contentHandler" is used for rendering EC2 UI screen
func contentHandler(resp http.ResponseWriter, req *http.Request) {

	region := req.FormValue("region")

	// fmt.Println(">>>>print thessss Region =>", req.Form)

	d := awsjack.Ec2Details(region)
	d.Regionlist = RegionList
	templates := populateTemplates()

	tmpl := templates.Lookup("content.html")
	tmpl.Execute(resp, d)
}

// "iamHandler" is used for rendering IAM user UI screen
func iamHandler(resp http.ResponseWriter, req *http.Request) {

	templates := populateTemplates()
	tmpl := templates.Lookup("iam.html")
	d := awsjack.IamDetails()
	d.Regionlist = RegionList
	fmt.Println("inside iam=", d)
	tmpl.Execute(resp, d)

}

// "viewHandler" is used for rendering Home UI screen
func viewHandler(resp http.ResponseWriter, req *http.Request) {
	templates := populateTemplates()
	tmpl := templates.Lookup("home.html")
	tmpl.Execute(resp, RegionList)

}

// "populateTemplates" is Parsing all HTML templates
func populateTemplates() *template.Template {
	result := template.New("templates")
	//fmt.Println("template=", result.Tree)
	const basePath = "templates"
	template.Must(result.ParseGlob(basePath + "/*.html"))
	return result
}

// This is main function which go willinvoke for starting this app

func main() {
	fmt.Println("AWSJack App Starting....")
	reginArg := os.Args[1]

	RegionList = strings.Split(reginArg, ",")
	//fmt.Println("############Command argument in main =", RegionList, "====", len(RegionList))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/get", contentHandler)
	http.HandleFunc("/iam", iamHandler)
	http.HandleFunc("/", viewHandler)

	http.ListenAndServe(":8080", nil)

	log.Fatal("error")

}
