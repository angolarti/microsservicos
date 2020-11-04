package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"github.com/hashicorp/go-retryablehttp"
)

type Result struct {
	Status string
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Ola</h1>")
	tmpt := template.Must(template.ParseFiles("templates/home.html"))
	tmpt.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {

	result := makeHttpCall("http://127.0.0.1:9091", r.FormValue("coupon"), r.FormValue("cc-number"))
	tmpt := template.Must(template.ParseFiles("templates/home.html"))
	tmpt.Execute(w, result)
}

func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		// log.Fatal("Microsservice pagamento out")
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close() // so é executado quando todo programa terminar de executar
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}
	json.Unmarshal(data, &result)

	return result

}