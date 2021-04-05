package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name                  string
	Age                   uint16
	Money                 int16
	Avg_Grades, Happiness float64
	Hobbies               []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("Name: %s, Age: %d, Money: %d, Average Grades: %.2f, "+
		"Happiness: %.2f", u.Name, u.Age, u.Money, u.Avg_Grades, u.Happiness)
}

func (u *User) setName(newName string) {
	u.Name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Reading", "Gaming", "Hiking"}}
	// bob.setName("Alex")
	// fmt.Fprint(w, "User name is: ", bob.Name, "\n")
	// fmt.Fprint(w, bob.getAllInfo())
	tmpl, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, bob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Contact information goes here.")
}

func handleRequest() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil)
}

func main() {
	// var bob User = ....
	// bob := User{name: "Bob", age: 25, money: -50, avg_grades: 4.2, happiness: 0.8}

	handleRequest()
}
