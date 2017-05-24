package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	db        *sql.DB
	db_prefix = "so_"
	//config config.Config
	//v      config.Vars
	t = template.Must(template.ParseGlob("views/*"))
)

type post struct {
	ID           int
	Post_title   string
	Post_content template.HTML
	Post_date    string
}

func init() {
	db, _ = sql.Open("mysql", "root:root@/wordpres1?charset=utf8")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{page:[0-9]+}.html", HomeHandler)
	//http.ListenAndServe(":9090", nil)
	log.Fatal(http.ListenAndServe(":9090", r))
}

/*首页*/
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	page := mux.Vars(r)["page"]
	if page != "" {
		page = "1"
	}

	rows, err := db.Query("SELECT ID,post_title,post_content,post_date FROM " + db_prefix + "posts ORDER BY ID DESC LIMIT 1,10")

	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.ID, &p.Post_title, &p.Post_content, &p.Post_date)
		p.ID = p.ID
		p.Post_title = p.Post_title
		fmt.Println(p.Post_content)
		p.Post_content = template.HTML(p.Post_content)
		p.Post_date = p.Post_date

		posts = append(posts, p)
	}
	//data := map[string][]string{}ced

	//data["list"] = posts
	//fmt.Println(posts)

	renderTemplate(w, "index.html", posts)
	//t.ExecuteTemplate(w, "index.html", posts)
}

/*模板解析*/
//func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	//t = t.Funcs(template.FuncMap{"unescaped": unescaped})
	err := t.ExecuteTemplate(w, tmpl, data)

	// Things will be more elegant than this: just a placeholder for now!
	if err != nil {
		http.Error(w, "error 500:"+" "+err.Error(), http.StatusInternalServerError)
	}
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

//http://studygolang.com/articles/1741
func unescaped(x string) interface{} {
	return template.HTML(x)
}
