package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/qiucarson/blog/config"
)

var (
	db        *sql.DB
	db_prefix string
	//config config.Config
	//v      config.Vars
	t = template.Must(template.ParseGlob("views/*.html"))
)

const (
	PAGE_MAX = 10
)

type post struct {
	ID           int
	Post_title   string
	Post_content template.HTML
	Post_date    string
}

func init() {
	c := config.MysqlConfig()
	db, _ = sql.Open("mysql", c.Mysql)
	db_prefix = c.Prefix
	//db, _ = sql.Open("mysql", "test:123456@/bolgsong?charset=utf8")
}

func main() {

	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(HomeHandler)
	r.Methods("GET").Path("/page/{page:[0-9]+}").HandlerFunc(HomeHandler)
	r.Methods("GET").Path("/{id:[0-9]+}.html").HandlerFunc(PostsHandler)
	//http.ListenAndServe(":9090", nil)
	log.Fatal(http.ListenAndServe(":9090", r))
}

/*首页*/
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	page := mux.Vars(r)["page"]
	if page == "" {
		page = "1"
	}
	pages, err := strconv.Atoi(page)
	pages = (pages - 1) * PAGE_MAX
	page_max := strconv.Itoa(PAGE_MAX)

	stmt, _ := db.Prepare("SELECT ID,post_title,post_content,post_date FROM " + db_prefix + "posts 	WHERE post_type='post' AND post_status = 'publish' ORDER BY ID DESC LIMIT ?,?")
	rows, err := stmt.Query(pages, page_max)

	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	posts := []post{}
	for rows.Next() {
		p := post{}
		var Post_content string
		rows.Scan(&p.ID, &p.Post_title, &Post_content, &p.Post_date)
		p.ID = p.ID
		p.Post_title = p.Post_title
		//gifmt.Println(p.Post_content)
		p.Post_content = template.HTML(Tagtotext(Post_content))
		//fmt.Println(Tagtotext(Post_content))
		p.Post_date = p.Post_date

		posts = append(posts, p)
	}

	renderTemplate(w, "index.html", posts)
	//t.ExecuteTemplate(w, "index.html", posts)
}

/*文章页*/
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	p := post{}
	var Post_content string
	stmt, _ := db.Prepare("SELECT ID,post_title,post_content,post_date FROM " + db_prefix + "posts 	WHERE post_type='post' AND post_status = 'publish' AND ID=?")
	rows, err := stmt.Query(id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	rows.Next()
	rows.Scan(&p.ID, &p.Post_title, &Post_content, &p.Post_date)
	p.Post_content = template.HTML(Post_content)
	renderTemplate(w, "posts.html", p)

}

/*模板解析*/
//func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {

	err := t.ExecuteTemplate(w, tmpl, data)

	// Things will be more elegant than this: just a placeholder for now!
	if err != nil {
		http.Error(w, "error 500:"+" "+err.Error(), http.StatusInternalServerError)
	}
}

func Tagtotext(content string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllString(content, "\n")
	return content
}
