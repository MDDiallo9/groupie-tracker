package main

import (
	"groupie-tracker/functions"
)

func main() {
	// #MARK: Lancement du serveur

	functions.DataAPI()
	functions.DataRecolt(functions.Links{})
	// functions.ServerGroupieTracker()

}
