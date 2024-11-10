package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Define a Person struct
type Person struct {
	ID   int
	Name string
	Age  int
	Bio  string
}

// Parse templates for HTML rendering
var templates = template.Must(template.ParseGlob("static/*.html"))

// Database connection
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "people.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize the database (create table if it doesn't exist)
	createTable()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/auth/", authHandler)
	http.HandleFunc("/poeple-list/", poepleListHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Function to create the table if it doesn't already exist
func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS people (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "name" TEXT,
        "age" INTEGER,
        "bio" TEXT
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Insert sample data if table is empty
	insertSampleData()
}

// Insert sample data if the table is empty
func insertSampleData() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM people").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	sampleData := []Person{
		{Name: "Alice", Age: 30, Bio: "Loves data engineering and Rust programming."},
		{Name: "Bob", Age: 25, Bio: "An enthusiast of Ubuntu and networking."},
		{Name: "Carol", Age: 35, Bio: "Enjoys Ansible, board games, and Factorio."},
	}

	for _, person := range sampleData {
		_, err := db.Exec("INSERT INTO people (name, age, bio) VALUES (?, ?, ?)", person.Name, person.Age, person.Bio)
		if err != nil {
			log.Fatal("Failed to insert sample data:", err)
		}
	}
}

// Handler for the home page (list of people)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	people, err := getAllPeople()
	if err != nil {
		http.Error(w, "Error fetching people", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", people)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Retrieve all people from the database
func getAllPeople() ([]Person, error) {
	rows, err := db.Query("SELECT id, name, age, bio FROM people")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name, &person.Age, &person.Bio)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, nil
}

// Handler for profile pages
func profileHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/profile/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	person, err := getPersonByID(id)
	if err != nil || person == nil {
		http.NotFound(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "profile.html", person)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Retrieve a person by ID from the database
func getPersonByID(id int) (*Person, error) {
	var person Person
	err := db.QueryRow("SELECT id, name, age, bio FROM people WHERE id = ?", id).Scan(&person.ID, &person.Name, &person.Age, &person.Bio)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &person, nil
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already authenticated
	fmt.Println("authHandler")
	cookie, err := r.Cookie("authenticated")
	if err == nil && cookie.Value == "true" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("r.Method", r.Method)
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "admin" && password == "1234" {
			cookie := http.Cookie{
				Name:  "authenticated",
				Value: "true",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	err = templates.ExecuteTemplate(w, "auth.html", nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func poepleListHandler(w http.ResponseWriter, r *http.Request) {
	people, err := getAllPeople()
	if err != nil {
		http.Error(w, "Error fetching people", http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "poeple-list.html", people)
	if err != nil {
		fmt.Println("Error rendering template", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
