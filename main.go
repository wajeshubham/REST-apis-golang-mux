package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Author model (Struct)
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Course model (Struct)
type Course struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Price  string  `json:"price"`
	Link   string  `json:"link"`
	Author *Author `json:"author"`
}

// Initialize Course variable as slice (which nothing but an Array data structure)
var courses []Course

// function to get all the books
// it's like we do in express as app.get("/",(req,res)=>{}) this is same as following
func getCourses(res http.ResponseWriter, req *http.Request) {
	// we have to set the header "Content-Type: application/json"
	// because we are sending JSON data with a request through postman
	res.Header().Set("Content-Type", "application/json")
	// we are taking variable 'courses' in which we've appended dummy data and just returning that as a response
	json.NewEncoder(res).Encode(courses)
}

// function to get a single course
func getSingleCourse(res http.ResponseWriter, req *http.Request) {
	// set the header "Content-Type: application/json"
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // we are extracting 'id' of the Course which we are passing in the url
	// we are looping throuh courses and
	// if we get matching course id with 'id' that we've passed in url
	// we are encoding that course only
	for _, item := range courses {
		if item.ID == params["id"] {
			json.NewEncoder(res).Encode(item) // we are sending matched course in json format
			return
		}
	}
	// if we don't get any matching course we will simply return a message
	// which means course dowsn't exist with provided id
	json.NewEncoder(res).Encode("No course found")

}

func createCourse(res http.ResponseWriter, req *http.Request) {
	// set the header "Content-Type: application/json"
	res.Header().Set("Content-Type", "application/json")
	// we are assigning course variable of type Course struct
	var course Course
	// we will take everyhing that we are passing through postman
	// which is nothing but req.Body and we will store decoded req.Body in 'course' variable
	_ = json.NewDecoder(req.Body).Decode(&course)
	course.ID = strconv.Itoa(rand.Intn(1000000)) // just creating dummy id as we are not sending an id from postman - not safe for production
	// Now, we will just append that course with decoded req.Body in our
	// courses slice (Array) which we have defined at the start
	courses = append(courses, course)
	// At the end we will simply send the corse that we have created in the response
	json.NewEncoder(res).Encode(course)
}

func updateCourse(res http.ResponseWriter, req *http.Request) {
	// set the header "Content-Type: application/json"
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // we are extracting 'id' of the Course which we are passing in the url

	// this method is bit cheeky because we don't have a database
	for i, item := range courses {
		if item.ID == params["id"] {
			// this is the same slicing technique that we used to delete the course in deleteCourse() func
			courses = append(courses[:i], courses[i+1:]...)
			// we are defining a course variable of type Course struct
			var course Course
			// then we are taking data which we are sending from postman which is nothing but req.Body
			// and creating the new course from it
			_ = json.NewDecoder(req.Body).Decode(&course)
			course.ID = params["id"] // we are keeping the id same because we are updating the existing course not creating new one - not safe for production
			// Now, we will just append that course in courses slice(array)
			courses = append(courses, course)
			json.NewEncoder(res).Encode(course)
			return
		}
	}
}

func deleteCourse(res http.ResponseWriter, req *http.Request) {
	// set the header "Content-Type: application/json"
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // we are extracting 'id' of the Course which we are passing in the url

	// We are looping through all the courses we have
	//then we will compare each course's ID with the id that we've passed in the url (params)
	for i, item := range courses {
		if item.ID == params["id"] {
			// This is slicing of a slice(array)
			// we are appending all the courses in `courses` slice (array) except the one which has ID equal to the id we've passed in url
			// so suppose if ID of course, whose index is "3" is equal to the id we've passed in url. So,
			// it will take items from index "0-2" and "4-`last index`" i.e. except index "3" from `courses` slice(array)
			// and will store that slice(array) in courses itself
			courses = append(courses[:i], courses[i+1:]...)
			break
		}
	}
	json.NewEncoder(res).Encode(courses) // it will return all the other courses except the deleted one.
}

func main() {
	//Initialize router
	router := mux.NewRouter()

	// Dummy data of books as we don't have a database (slice)
	// we are just appending a Course struct
	// which we are creating on the spot in courses variable that we've defined above
	courses = append(courses, Course{ID: "124134", Name: "FullStack Django Developer Freelance ready", Price: "299", Link: "https://courses.learncodeonline.in/learn/FullStack-Django-Developer-Freelance-ready",
		Author: &Author{Firstname: "Hitesh", Lastname: "Choudhary"}})
	courses = append(courses, Course{ID: "154434", Name: "Full stack with Django and React", Price: "299", Link: "https://courses.learncodeonline.in/learn/Full-stack-with-Django-and-React",
		Author: &Author{Firstname: "Hitesh", Lastname: "Choudhary"}})
	courses = append(courses, Course{ID: "198767", Name: "Complete React Native bootcamp", Price: "199", Link: "https://courses.learncodeonline.in/learn/Complete-React-Native-Mobile-App-developer",
		Author: &Author{Firstname: "Hitesh", Lastname: "Choudhary"}})

	// Handel the routes
	router.HandleFunc("/api/courses", getCourses).Methods("GET")
	router.HandleFunc("/api/course/{id}", getSingleCourse).Methods("GET")
	router.HandleFunc("/api/courses/create", createCourse).Methods("POST")
	router.HandleFunc("/api/courses/update/{id}", updateCourse).Methods("PUT")
	router.HandleFunc("/api/courses/delete/{id}", deleteCourse).Methods("DELETE")

	// Initialize a server
	log.Fatal(http.ListenAndServe(":8000", router))

}
