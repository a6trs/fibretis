package grass

import (
	"../soil"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var sstore = sessions.NewCookieStore([]byte("these-are-very-important-yeah"))

var templates, _ = template.New("IDONTKNOW").
	Funcs(template.FuncMap{"validuser": validUser, "account": account, "project": project, "post": post, "raw": rawhtml, "timestr": timestr, "nutshell": nutshell}).
	ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/_icons.svg", "flowers/index.html", "flowers/login.html", "flowers/signup.html", "flowers/profedit.html", "flowers/projects.html", "flowers/project_create.html", "flowers/project_page.html", "flowers/post_create.html", "flowers/post_page.html")

func validUser(aid int) bool {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	return (err == nil)
}

func account(aid int) *soil.Account {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	if err == nil {
		return acc
	} else {
		return nil
	}
}

func project(prjid int) *soil.Project {
	prj := &soil.Project{ID: prjid}
	err := prj.Load(soil.KEY_Project_ID)
	if err == nil {
		return prj
	} else {
		return nil
	}
}

func post(pstid int) *soil.Project {
	pst := &soil.Project{ID: pstid}
	err := pst.Load(soil.KEY_Post_ID)
	if err == nil {
		return pst
	} else {
		return nil
	}
}

func rawhtml(s string) template.HTML {
	return template.HTML(s)
}

func timestr(t time.Time) string {
	return t.Format(time.RFC822)
}

func nutshell(body string) string {
	// Simply remove all HTML tags.
	r, _ := regexp.Compile(`<\/?\w+(?:\s+\w+=['"].*['"])*>`)
	body = r.ReplaceAllString(body, "")
	// Help on rune arrays:
	// http://www.cnblogs.com/howDo/archive/2013/04/20/GoLang-String.html
	br := []rune(body)
	if len(body) <= 80 {
		return string(br)
	} else {
		return string(br[:80])+"..."
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, title string, arg map[string]interface{}) {
	arg["aid"] = accountInSession(w, r)
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func accountInSession(w http.ResponseWriter, r *http.Request) int {
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	s := sess.Values["id"]
	if s == nil {
		s = -1
	}
	return s.(int)
}
