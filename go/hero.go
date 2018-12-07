//https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"log"
	"net/http"
	//"reflect"
	"strconv"
	"sync"
)

// PROBLEM DESCRIPTION:
// Your goal here is to design an API that allows for hero tracking, much like the Vue problem
// You are to implement an API (for which the skeleton already exists) that has the following capabilities
// - Get      : return a JSON representation of the hero with the name supplied
// - Make     : create a superhero according to the JSON body supplied
// - Calamity : a calamity of the supplied level requires heroes with an equivalent combined powerlevel to address it.
//              Takes a calamity with powerlevel and at least 1 hero. On success return a 200 with json response indicating the calamity has been resolved.
//              Otherwise return a response indicating that the heroes were not up to the task. Addressing a calamity adds 1 point of exhaustion.
// - Rest     : recover 1 point of exhaustion
// - Retire   : retire a superhero, someone may take up his name for the future passing on the title
// - Kill     : a superhero has passed away, his name may not be taken up again.

// On success all endpoints should return a status code 200.

// You are free to decide what your API endpoints should be called and what shape they should take. You can modify any code in this file however you'd like.

// NOTE: you may want to install postman or another request generating software of your choosing to make testing easier. (api is running on localhost port 8081)

// NOTE the second: the API is receiving asynchronous requests to manage our super friends. As such, your hero access should be thread-safe for writes.
// Even if the operations are extremely lightweight we want to make our application scalable.

// NOTE the third: There are many ways to make whatever package-level tracking you implement thread-safe, feel free to change anything about this file (without changing the functionality of the program) to do so.
// i.e. add package-level maps, add functions that take the hero struct as reference, modify the hero struct, creating access control paradigms
// I highly recommend looking into channels, mutexes, and other golang memory management patterns and pick whatever you're most comfortable with.
// For mad props: a timeout on the memory management process which returns a resource not available.

// Bonus: If you're having fun (this is by no means necessary) you can make the calamity hold the heroes up for a time and delay their unlocking in a go-routine
// example:
// go func(h *hero) {
//     time.Sleep(calamityTime)
//     // release lock on hero
// }(heroPtr)

var maxExhaustion = 10

type hero struct {
	Name       string `json:"Name"`
	PowerLevel int    `json:"PowerLevel"`
	Exhaustion int    `json:"Exhaustion"`
	Status     string `json:"Status"`
}

type allheros struct {
	Total      int    `json:"Total Heros"`
	Active     int    `json:"Total Active"`
	Retired    int    `json:"Total Retired"`
	KIA        int    `json:"Total KIA"`
	Power      int    `json:"Total Power"`
	Exhaustion int    `json:"Total Exhaustion"`
	Heros      []hero `json:"Hero List"`
}

type restresult struct {
	Status string `json:"Status"`
	Hero   hero   `json:"Hero"`
}

type retireresult struct {
	Status string `json:"Status"`
	Hero   hero   `json:"Hero"`
}

type killedresult struct {
	Status string `json:"Status"`
	Hero   hero   `json:"Hero"`
}

type makeresult struct {
	Status string `json:"Status"`
	Hero   hero   `json:"Hero"`
}

type calamityresult struct {
	Status string `json:"Status"`
	Heros  []hero `json:"Heros Involved"`
}

var heros []hero

var mu sync.Mutex

// TODO: add storage and memory management

func herosGet(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var result allheros
	result.Heros = heros
	var tot int
	var pow int
	var exh int
	var kia int
	var ret int
	var act int
	tot = 0
	pow = 0
	exh = 0
	kia = 0
	ret = 0
	act = 0
	for _, item := range heros {
		pow = pow + item.PowerLevel
		exh = exh + item.Exhaustion
		if item.Status == "KIA" {
			kia = kia + 1
		} else if item.Status == "Retired" {
			ret = ret + 1
		} else if item.Status == "Active" {
			act = act + 1
		}
		tot = tot + 1
	}
	if tot > 0 {
		result.Total = tot
		result.Active = act
		result.Retired = ret
		result.KIA = kia
		result.Power = pow
		result.Exhaustion = exh
		result.Heros = heros
		json.NewEncoder(w).Encode(result)
	} else {
		json.NewEncoder(w).Encode("No heros found.")
	}
	mu.Unlock()
	return
}

func heroGet(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var name string
	var found bool
	var ok bool
	var last int
	if name, ok = mux.Vars(r)["name"]; !ok {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		found = false
		for index, item := range heros {
			if item.Name == name {
				found = true
				last = index
			}
		}
		if found == false {
			json.NewEncoder(w).Encode("Unable to find hero.")
		} else {
			json.NewEncoder(w).Encode(heros[last])
		}
	}
	mu.Unlock()
	return
}

func heroMake(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var h hero
	var result makeresult
	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		var valid = true
		for _, item := range heros {
			if item.Name == h.Name {
				if item.Status != "Retired" {
					valid = false
					h.Status = item.Status
					if item.Status == "Active" {
						json.NewEncoder(w).Encode("Hero creation failed. This hero name is already in active use.")
					} else if item.Status == "KIA" {
						json.NewEncoder(w).Encode("Hero creation failed. This hero name was retired when the previous holder was killed in action.")
					}
				}
			}
		}
		if valid == true {
			h.Status = "Active"
			h.Exhaustion = 0
			heros = append(heros, h)
			result.Hero = h
			result.Status = "Hero sucessfully created."
			json.NewEncoder(w).Encode(result)
		}
	}
	mu.Unlock()
	return
}

func heroKill(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var found bool
	var name string
	var ok bool
	var last int
	found = false
	var result killedresult
	if name, ok = mux.Vars(r)["name"]; !ok {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		found = false
		for index, item := range heros {
			if item.Name == name {
				last = index
				found = true
			}
		}
		if found == false {
			json.NewEncoder(w).Encode("Unable to find any hero by that name.")
		} else {
			result.Hero = heros[last]
			if heros[last].Status == "Active" {
				heros[last].Status = "KIA"
				heros[last].Exhaustion = 0
				heros[last].PowerLevel = 0
				result.Hero = heros[last]
				result.Status = "Hero Sucessfully Killed"
			} else if heros[last].Status == "KIA" {
				result.Status = "This hero was already dead."
			} else if heros[last].Status == "Retired" {
				result.Status = "This hero is already retired and cannot be killed in action."
			}
			json.NewEncoder(w).Encode(result)
		}
	}
	mu.Unlock()
	return
}

func heroRetire(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var found bool
	var name string
	var ok bool
	var last int
	found = false
	var result retireresult

	if name, ok = mux.Vars(r)["name"]; !ok {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		found = false
		for index, item := range heros {
			if item.Name == name {
				last = index
				found = true
			}
		}
		if found == false {
			json.NewEncoder(w).Encode("Unable to find any hero by that name.")
		} else {
			result.Hero = heros[last]
			if heros[last].Status == "Active" {
				heros[last].Status = "Retired"
				heros[last].Exhaustion = 0
				heros[last].PowerLevel = 0
				result.Hero = heros[last]
				result.Status = "Hero Sucessfully Retired"
			} else if heros[last].Status == "KIA" {
				result.Status = "This hero was killed in action before they could retire."
			} else if heros[last].Status == "Retired" {
				result.Status = "This hero is already retired."
			}
			json.NewEncoder(w).Encode(result)
		}
	}
	mu.Unlock()
	return
}

func heroRest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var name string
	var found bool
	found = false
	var ok bool
	var last int
	var result restresult

	if name, ok = mux.Vars(r)["name"]; !ok {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		for index, item := range heros {
			if item.Name == name {
				last = index
				found = true
			}
		}
		if found == false {
			json.NewEncoder(w).Encode("Unable to find hero.")
		} else {
			result.Hero = heros[last]
			if heros[last].Status == "KIA" {
				result.Status = "Hero is dead and cannot rest."
			} else if heros[last].Status == "Retired" {
				result.Status = "Hero is retired and does not need to rest."
			} else if heros[last].Status == "Active" {
				if heros[last].Exhaustion > 0 {
					heros[last].Exhaustion = heros[last].Exhaustion - 1
					result.Status = "Hero successfully rested."
				} else {
					result.Status = "Hero was already fully rested."
				}
			}
			json.NewEncoder(w).Encode(result)
		}
	}
	mu.Unlock()
	return
}

func calamity(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var p string
	var ok bool
	var result calamityresult
	if p, ok = mux.Vars(r)["power"]; !ok {
		json.NewEncoder(w).Encode("Input invalid.")
	} else {
		var actives []int //a list of indexes of active heros
		powerNeeded, err := strconv.Atoi(p)
		if err != nil {
			json.NewEncoder(w).Encode("Input could not be parsed as an integer.")
		} else {
			for i, item := range heros {
				if item.Status == "Active" {
					actives = append(actives, i)
				}
			}
			var hold int
			for i, spot := range actives { //sorts by exhaustion level. Heros with lowest exhaustion are used first.
				var min = i
				for j, item := range actives[i:] {
					if heros[item].Exhaustion < heros[min].Exhaustion {
						min = j
					}
				}
				hold = spot
				spot = actives[min]
				actives[min] = hold
			}
			var totalPower = 0
			for _, item := range actives {
				totalPower = totalPower + heros[item].PowerLevel
			}
			if totalPower < powerNeeded {
				var finalherolist []hero
				for _, item := range actives {
					heros[item].Exhaustion = heros[item].Exhaustion + 1
					finalherolist = append(finalherolist, heros[item])
				}
				result.Status = "The combined power of all of the active heros was not enough to overcome the calamity. All active heros have become more exhausted."
				result.Heros = finalherolist
				json.NewEncoder(w).Encode(result)
			} else {
				var starters []int
				for powerNeeded > 0 {
					powerNeeded = powerNeeded - heros[actives[0]].PowerLevel
					starters = append(starters, actives[0])
					actives = actives[1:]
				}
				var finalherolist []hero
				for _, item := range starters {
					heros[item].Exhaustion = heros[item].Exhaustion + 1
					finalherolist = append(finalherolist, heros[item])
				}
				result.Status = "The heros took care of the calamity. All heros involved became more exhausted."
				result.Heros = finalherolist
				json.NewEncoder(w).Encode(result)
			}
		}
	}
	mu.Unlock()
	return
}

func linkRoutes(r *mux.Router) {
	r.HandleFunc("/hero", heroMake).Methods("POST")
	r.HandleFunc("/hero", herosGet).Methods("GET")
	r.HandleFunc("/hero/{name}", heroGet).Methods("GET")
	r.HandleFunc("/hero/{name}/retire", heroRetire).Methods("PATCH")
	r.HandleFunc("/hero/{name}/rest", heroRest).Methods("PATCH")
	r.HandleFunc("/hero/{name}/kill", heroKill).Methods("PATCH")
	r.HandleFunc("/calamity/{power}", calamity).Methods("PATCH")
}

func main() {
	fmt.Printf("1\n")
	// create a router
	router := mux.NewRouter()
	fmt.Printf("2\n")
	// and a server to listen on port 8081
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: router,
	}
	fmt.Printf("3\n")
	// link the supplied routes
	linkRoutes(router)
	fmt.Printf("4\n")
	// wait for requests
	log.Fatal(server.ListenAndServe())
	fmt.Printf("5\n")
}

//https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo

//https://pragmacoders.com/building-a-json-api-in-golang/
