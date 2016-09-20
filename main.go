package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	//"github.com/dgrijalva/jwt-go"
	//"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/qbit/homestead/src"
	"github.com/qbit/pgenv"
)

var insecure bool
var cookieSecret string
var crsfSecret string
var jwtSecret string
var templ *template.Template
var store *sessions.CookieStore
var listen string
var tf = "2006-01-02T15:04:05.999Z"

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format(time.RFC1123)
	},
	"shortDate": func(t time.Time) string {
		return t.Format(time.RFC822)
	},
	"printByte": func(b []byte) string {
		return string(b)
	},
	"printHTML": func(b []byte) template.HTML {
		return template.HTML(string(b))
	},
}

func init() {
	var err error
	flag.BoolVar(&insecure, "i", false, "Insecure mode")
	flag.StringVar(&cookieSecret, "cookie", "something-very-secret", "Secret to use for cookie store")
	flag.StringVar(&crsfSecret, "crsf", "32-byte-long-auth-key", "Secret to use for cookie store")
	flag.StringVar(&jwtSecret, "jwt", "super secret neat", "Secret to use for jwt")
	flag.StringVar(&listen, "http", ":8080", "Listen on")

	flag.Parse()

	store = sessions.NewCookieStore([]byte(cookieSecret))
	templ, err = template.New("homestead").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	gob.Register(&homestead.User{})
}

func main() {
	var cstr = pgenv.ConnStr{}
	cstr.SetDefaults()

	db, err := sql.Open("postgres", cstr.ToString())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router := mux.NewRouter()
	router.PathPrefix("/public/").Handler(
		http.StripPrefix("/public/",
			http.FileServer(http.Dir("public"))))

	router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		err = templ.ExecuteTemplate(w, "data.html", nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("/data/store", func(w http.ResponseWriter, r *http.Request) {
		var log homestead.Log
		r.ParseMultipartForm(32 << 20) // 32 mb

		for k, v := range r.Form {
			switch k {
			case "timestamp":
				log.Stamp, err = time.Parse(tf, strings.Join(v, ""))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "sensor":
				log.SensorName = strings.Join(v, "")
			default:
				log.Metrics = append(log.Metrics, fmt.Sprintf(`"%s"=>"%s"`, k, strings.Join(v, "")))
			}

		}

		_, err := log.SetID(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call InsertLog

		logID, err := homestead.InsertLog(db, &log)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fmt.Sprintf(`{"status": "OK", "id": %d}`, logID))
	})

	router.HandleFunc("/data/sensors", func(w http.ResponseWriter, r *http.Request) {
		val, err := homestead.GetSensors(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, *val)
	})

	router.HandleFunc("/data/current/{sensor}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sensor := vars["sensor"]
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val, err := homestead.GetCurrent(db, sensor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, *val)
	})
	/*
				router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
					session, err := store.Get(r, "session-name")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					session.Options = &sessions.Options{
						MaxAge: -1,
					}
					session.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)

				})

	p

				router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
					authHeader := r.Header.Get("Authorization")
					if authHeader == "" {
						http.Error(w, "Not Authorized!", http.StatusUnauthorized)
						return
					}

					token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
						return []byte(jwtSecret), nil
					})

					if err != nil {
						http.Error(w, err.Error(), http.StatusUnauthorized)
						return
					}

					if token.Valid {
						session, err := store.Get(r, "session-name")
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						data := r.FormValue("data")

			}
				})

				router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					err := templ.ExecuteTemplate(w, "index.html", map[string]interface{}{
						csrf.TemplateTag: csrf.TemplateField(r),
					})
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				})
				router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
					session, err := store.Get(r, "session-name")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					uVal := session.Values["user"]
					var u, ok = uVal.(*homestead.User)
					if !ok {
						uVal = &homestead.User{}
						session.Values["user"] = &uVal
						session.Save(r, w)
					}

					if ok && u.Admin {
						err = templ.ExecuteTemplate(w, "admin.html", nil)

						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

					}
				})
				router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
					session, err := store.Get(r, "session-name")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					uVal := session.Values["user"]
					var u, ok = uVal.(*homestead.User)
					if !ok {
						uVal = &homestead.User{}
						session.Values["user"] = &uVal
						session.Save(r, w)
					}

					if ok && u.Admin {
						// grab data here and send it to template
						err = templ.ExecuteTemplate(w, "data.html", nil)

						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

					}
				})

				router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
					session, err := store.Get(r, "session-name")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					user := r.FormValue("user")
					passwd := r.FormValue("passwd")

					if user == "" && passwd == "" {
						http.Redirect(w, r, "/", http.StatusFound)
					} else {
						u, err := homestead.Auth(db, user, passwd)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						if u.Authed {
							session.Values["user"] = u
							session.Save(r, w)
							http.Redirect(w, r, "/data", http.StatusFound)
						}
					}

				})
	*/
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	/*
		if insecure {
			log.Fatal(http.ListenAndServe(listen,
				csrf.Protect([]byte("32-byte-long-auth-key"),
					csrf.Secure(false))(loggedRouter)))
		} else {
			log.Fatal(http.ListenAndServe(listen,
				csrf.Protect([]byte(crsfSecret))(loggedRouter)))
		}
	*/
	log.Fatal(http.ListenAndServe(listen, loggedRouter))
}