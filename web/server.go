package web

// base on simple web server - ref https://reinbach.com/golang-webapps-1.html

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/baabeetaa/glogchain/config"
	"encoding/hex"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"encoding/gob"
)

type Context struct {
	Title  string
	Static string
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "index", context)
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func PostCreateHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "PostCreate"}
	render(w, "PostCreate", context)
}

func PostCreateSave(w http.ResponseWriter, req *http.Request) {
	//var postOperation protocol.PostOperation
	Title := req.FormValue("Title")
	Author := req.FormValue("Author")
	Body := req.FormValue("Body")

	var newPostString string = "{\"Type\": \"PostOperation\" , \"Operation\" : {\"Title\": \"" +  Title + "\", \"Body\": \"" + Body + "\", \"Author\": \"" + Author + "\"} }"
	fmt.Printf("---------------------------------\n")
	fmt.Printf("newPostString: %s\n", newPostString)

	newPostStringHex := make([]byte, len([]byte(newPostString)) * 2)
	hex.Encode(newPostStringHex, []byte(newPostString))
	fmt.Println("newPostStringHex: %s\n", newPostStringHex)

	/// example
	// {"Type": "PostOperation" , "Operation" : {"Title": "aaa", "Body": "aaa", "Author": "aaa"} }
	// http://10.0.0.11:46657/broadcast_tx_commit?tx=%227b2254797065223a2022506f73744f7065726174696f6e22202c20224f7065726174696f6e22203a207b225469746c65223a2022616161222c2022426f6479223a2022616161222c2022417574686f72223a2022616161227d207d%22

	var url_request string = config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit?tx=%22" + string(newPostStringHex[:]) + "%22"
	log.Print("url_request: %#v\n", url_request)
	resp, err := http.Get(url_request)
	//req, err := http.NewRequest("GET", url_request, nil)
	//resp, err := http.PostForm(config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit", url.Values{"tx": {"%22" + string(newPostStringHex[:]) + "%22"}})
	if err != nil {
		log.Print("http.Get error: %#v\n", err)
		return;
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ioutil.ReadAll error: %#v\n", err)
		return;
	}
	json_response_string := string(body[:])
	fmt.Println("json_response_string: %#v\n", json_response_string)

	// delay sometime to make sure hugo loading new page content
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, config.GlogchainConfigGlobal.HugoBaseUrl + "/post/" + Title  + "/", http.StatusFound)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = "/static/"
	tmpl_list := []string{"web/templates/base.html", fmt.Sprintf("web/templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StartWebServer() error  {
	gob.Register(&User{})

	// use https://github.com/gorilla/mux
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	r.HandleFunc("/", AuthWrapper(HomeHandler))

	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)

	r.HandleFunc("/about/", AuthWrapper(AboutHandler))
	r.HandleFunc("/post/create", AuthWrapper(PostCreateHandler))

	// Subrouter
	//s := r.PathPrefix("/secure").Subrouter()
	//// /secure/test
	//s.HandleFunc("/test", AuthWrapper(HomeHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         config.GlogchainConfigGlobal.GlogchainWebAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Bind to a port and pass our router in
	log.Fatal(srv.ListenAndServe())

	return nil
}
