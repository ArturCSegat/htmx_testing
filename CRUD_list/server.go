package main

import (
	"net/http"
	"html/template"
    "github.com/gorilla/mux"
)

type User struct {
    Name string
    Role string
}

type Content struct {
    Users []User
}

var Users []User = []User{
    User{Name: "Artur", Role: "Funny"},
    User{Name: "Jullia", Role: "Linda"},
    User{Name: "Peter", Role: "Griffin"},
}

func home(w http.ResponseWriter, r * http.Request) {
    if r.Method == "GET" {
        tmpl, err := template.ParseGlob("templates/*.html")
        if err != nil {
            panic(err)
        }
        cont := Content{Users: Users}

        tmpl.Execute(w, cont)
    }
}

// I know i can write this better with empty return but i want it like this
func delete_user(w http.ResponseWriter, r * http.Request) {
    if r.Method == "DELETE" {
        tmpl, err := template.ParseGlob("templates/*.html")
        if err != nil {
            panic(err)
        }

        vars := mux.Vars(r)
        name, ok := vars["name"]
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            tmpl.ExecuteTemplate(w, "users", Content{Users: Users})
            return
        }

        for i, user := range Users {
            if user.Name == name {
                Users = append(Users[:i], Users[i+1:]...)
                cont := Content{Users: Users}
                w.WriteHeader(http.StatusOK)
                tmpl.ExecuteTemplate(w, "users", cont)
                return
            }
        }
        w.WriteHeader(http.StatusNotFound)
        tmpl.ExecuteTemplate(w, "users", Content{Users: Users})
    }
}


func swap_lol (w http.ResponseWriter, r * http.Request) {
    tmpl, err := template.ParseFiles("templates/user.html")
    if err != nil {
        panic(err)
    }

    tmpl.ExecuteTemplate(w, "user", User{Name: "Lol", Role: "Funnny"})
}

func new_user (w http.ResponseWriter, r * http.Request) {
    if r.Method == "POST" { 
        tmpl, err := template.ParseFiles("templates/users.html", "templates/user.html")
        if err != nil {
            panic(err)
        }
        Users = append(Users, User{Name: r.FormValue("name"), Role: r.FormValue("role")})

        tmpl.ExecuteTemplate(w, "users", Content{Users: Users})
    }
}

func get_all_users_with (w http.ResponseWriter, r * http.Request) {
    if r.Method == "POST" {
        crr_str := r.FormValue("search")
        valid_users := make([]User, 0)
        
        if crr_str != "" {
            outer: for _, usr := range Users {
                for i, _:= range crr_str {
                    if crr_str[i] != usr.Name[i] {
                        continue outer           
                    }
                }
                valid_users = append(valid_users, usr)
            }
        }
        tmpl, err := template.ParseFiles("templates/search_results.html", "templates/search_item.html")
        if err != nil {
            panic(err)
        }
            
        type SearchResults struct {
            Results []User
        }

        tmpl.ExecuteTemplate(w, "search-results", SearchResults{Results: valid_users})
    }
}

func main(){
    router := mux.NewRouter()
    router.HandleFunc("/", home)
    router.HandleFunc("/delete-user/{name}", delete_user)
    router.HandleFunc("/swap", swap_lol)
    router.HandleFunc("/new-user", new_user)
    router.HandleFunc("/search-user", get_all_users_with)

    http.ListenAndServe("localhost:6969", router)
}
