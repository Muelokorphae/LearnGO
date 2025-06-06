package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	ID   int  `json:"id"`
	Name string  `json:"name"`
	Price      float64 `json:"price"`
	Instructor string  `json:"instructor"`
}

var Courselist []Course

func init() {
	CourseJSON := `[
		{"id": 1, "name": "Go Programming", "price": 199.99, "instructor": "John Doe"},
		{"id": 2, "name": "Python Programming", "price": 149.99, "instructor": "Jane Smith"},
		{"id": 3, "name": "JavaScript Basics", "price": 99.99, "instructor": "Alice Johnson"}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &Courselist)
	if err != nil {
		log.Fatal(err)
	}
}

func getNewID() int {
	highestID := -1
	for _, course := range Courselist {
		if highestID < course.ID {
			highestID = course.ID
		}
	}
	return highestID + 1
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}	
	course, listItemIndex:= findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("No course with ID %d found",ID), http.StatusNotFound)
		return
	}

	switch r.Method {
		case http.MethodGet:
			courseJSON, err := json.Marshal(course)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(courseJSON)
		case http.MethodPut:
			var updatedCourse Course
			Bodybyte, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(Bodybyte, &updatedCourse)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}	
			if updatedCourse.ID != ID {
				w.WriteHeader(http.StatusBadRequest)
				return
				
			}
			course = &updatedCourse
			Courselist[listItemIndex] = *course
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)


	}

}

func findID(ID int) (*Course, int) {
	for i , course := range Courselist {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON, err := json.Marshal(Courselist)
	switch r.Method {
		case http.MethodGet:
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}	
			w.Header().Set("Content-Type", "application/json")
			w.Write(courseJSON)
		case http.MethodPost:
			var newCourse Course
			Bodybyte, err := io.ReadAll(r.Body)		
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(Bodybyte, &newCourse)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			newCourse.ID = getNewID()
			Courselist = append(Courselist, newCourse)
			w.WriteHeader(http.StatusCreated)
			return
	}	
		
}

// //middleware 
// func MiddlewareHandler(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("before handler middleware start")
// 		handler.ServeHTTP(w, r)
// 		fmt.Printf("\nafter handler middleware end")
// 	})


// }

// CorsMiddleware
func enableCorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With, Origin, Cache-Control, Pragma, X-HTTP-Method-Override, X-HTTP-Method, X-HTTP-Method-Override, X-HTTP-Method-Override, X-HTTP-Method, X-HTTP-Method-Override")
		handler.ServeHTTP(w, r)
	})
}






func main() {
	courseItemHandler := http.HandlerFunc(courseHandler)
	coruseListHandler := http.HandlerFunc(coursesHandler)
	http.Handle("/course/", enableCorsMiddleware(courseItemHandler))
	http.Handle("/course", enableCorsMiddleware(coruseListHandler))
	http.ListenAndServe(":5000", nil)
}