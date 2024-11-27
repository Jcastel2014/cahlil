package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"sysadmin.com/final/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
	}
	ts, err := template.ParseFiles("./ui/tmpl/index.tmpl")

	if err != nil {
		panic(err.Error())
	}

	err = ts.Execute(w, nil)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	username := strings.TrimSpace(r.PostForm.Get("username"))
	password := strings.TrimSpace(r.PostForm.Get("password"))

	errors_user := make(map[string]string)

	if username == "" {
		errors_user["username"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(username) > 30 {
		errors_user["username"] = "This field is too long"
	} else if utf8.RuneCountInString(username) < 3 {
		errors_user["username"] = "This field is too short"
	}

	if password == "" {
		errors_user["password"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(password) > 60 {
		errors_user["password"] = "This field is too long"
	} else if utf8.RuneCountInString(password) < 5 {
		errors_user["password"] = "This field is too short"
	}

	if len(errors_user) > 0 {
		ts, err := template.ParseFiles("./ui/tmpl/signup.page.tmpl")

		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, &TemplateData{
			ErrorsFromForm:  errors_user,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})
		if err != nil {
			panic(err.Error())
		}
		return
	}

	err = app.user.Insert(username, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			errors_user["username"] = "username already exists in the system"
			ts, err := template.ParseFiles("./ui/tmpl/signup.page.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &TemplateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
			})
			if err != nil {
				app.serverError(w, err)
			}
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	app.session.Put(r, "flash", "new user added")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/tmpl/signup.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
	})

	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/tmpl/login.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
	})

	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	errors_user := make(map[string]string)
	id, err := app.user.Authenticate(username, password)
	if err != nil {
		log.Print("didnt authenticate")
		if errors.Is(err, models.ErrInvalidCredentials) {
			errors_user["default"] = "username or password is incorrect"
			ts, err := template.ParseFiles("./ui/tmpl/login.page.tmpl") //load the template file

			if err != nil {
				log.Println(err.Error())
				app.serverError(w, err)
				return
			}

			err = ts.Execute(w, &TemplateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
			})

			if err != nil {
				log.Panicln(err.Error())
				app.serverError(w, err)
				return
			}
			return
		}
		return
	}
	app.session.Put(r, "authenticatedID", id)
	app.session.Put(r, "flash", "You have logged in.")
	log.Print("authenticated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedID")
	app.session.Put(r, "flash", "You have logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
