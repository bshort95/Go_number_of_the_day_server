package main

// file imports
import (
	"fmt"
	"hash/fnv"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

//**************************************
// this function is used to calculate a number
// the rules for this number are as follows
// - the number has to be the same number regardless
// of who accesses the server
// - the number must change daily
// the number must be between 0-100
// the number can not be 42
// that is different everyday in the utc time zone
//**************************************
func calculate_number() int {
	var array1 [3]int
	var array2 [5]int
	constarr1 := [3]int{3, 7, 11}
	constarr2 := [5]int{1, 2, 3, 5, 8}
	var totarr1 [3]int
	var totarr2 [5]int
	var totalstr string = ""
	var cons1 string = "Bro"
	var cons2 string = "Jones"
	currentDay := time.Now().Weekday().String()
	currentMonth := time.Now().Month().String()
	currentNumDay := time.Now().Day()
	cons1chars := []rune(cons1)
	cons2chars := []rune(cons2)
	daychars := []rune(currentDay)
	monthchars := []rune(currentMonth)

	for i := 0; i < 5; i++ {
		array2[i] = int(daychars[i]) + int(cons2chars[i])
	}

	for i := 0; i < 3; i++ {
		array1[i] = int(monthchars[i]) + int(cons1chars[i])
	}

	for i := 0; i < 3; i++ {
		totarr1[i] = constarr1[i] * array1[i]
	}

	for i := 0; i < 5; i++ {
		totarr2[i] = constarr2[i] * array2[i]
	}

	for i := 0; i < 3; i++ {
		totalstr = totalstr + string((totarr1[i]%92)+32)
	}

	for i := 0; i < 5; i++ {
		totalstr = totalstr + string((totarr2[i]%92)+32)
	}

	h := fnv.New32a()
	h.Write([]byte(totalstr))

	var totalnum int = int(h.Sum32())

	totalnum = totalnum * currentNumDay
	var finalnum int = totalnum % 100

	if finalnum == 42 {
		finalnum = 41

	}
	return finalnum

}

// *************************************
// the controllers for the server
// *************************************
func routes_and_controllers(w http.ResponseWriter, r *http.Request) {
	// if the url path is not reconized as one of our own, return a 404 page
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "views\\404.html")
		return
	}

	switch r.Method {

	// requests that were made
	//get request
	case "GET":

		t, err := template.ParseFiles("views\\wotd.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		// is able to send variable to the html
		myvar := map[string]interface{}{"MyVar": "take a guess"}
		if err := t.Execute(w, myvar); err != nil {
			http.Error(w, err.Error(), 500)
		}

	// post request
	case "POST":

		// post sends the number we are calculating
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		// number we recieved
		number := r.FormValue("number")

		//secret cool calculation
		var notd int = calculate_number()

		t, err := template.ParseFiles("views\\wotd.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// intValue := 0
		num, err := strconv.Atoi(number)
		if err != nil {
			fmt.Printf("Supplied value %s is not a number\n", number)
		} else {

			//check to see if inputed number is the same as the number of the day
			if notd == num {
				myvar := map[string]interface{}{"MyVar": "number is correct"}
				if err := t.Execute(w, myvar); err != nil {
					http.Error(w, err.Error(), 500)
				}

			} else {
				myvar := map[string]interface{}{"MyVar": "that is incorrect, the correct number is " + strconv.FormatInt(int64(notd), 10)}
				if err := t.Execute(w, myvar); err != nil {
					http.Error(w, err.Error(), 500)
				}

			}
		}
		// handles only expected results
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// main function. sort of like the routes.
func main() {
	http.HandleFunc("/", routes_and_controllers)

	fmt.Printf("Starting server... on 8081\n")
	http.ListenAndServe(":8081", nil)
}
