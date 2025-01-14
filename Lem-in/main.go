package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	//"strconv"
)

var startRoom, endRoom *string
var antsNumber string
var visitedRooms []string
var roomsAndTunnels = make(map[string][]string)

func main() {
	readInput()
	fmt.Println("Ants number: ", antsNumber)
	fmt.Println("Start room: ", *startRoom)
	fmt.Println("End room: ", *endRoom)
	for room, tunnelTo := range roomsAndTunnels {
		fmt.Printf("room: %s, connected rooms: %s\n", room, tunnelTo)
	}
	var currentRoom = *startRoom
	var possiblePath []string
	foundPaths := searchForPath(currentRoom, visitedRooms, possiblePath)
	uniquepaths := isUntersect(foundPaths)

	pathSets := [][][]string{}
	for i := range foundPaths {
		setsOfPaths(foundPaths, [][]string{}, i, &pathSets)
	}

	for _, set := range pathSets {
		fmt.Println(set)
	}

	fmt.Println("Lenth of found paths: ", len(foundPaths))
	fmt.Println("Found paths: ", foundPaths)
	for room, tunnelTo := range uniquepaths {
		fmt.Printf("combination number: %v, paths: %s\n", room, tunnelTo)
	}
	//moveAnts(uniquepaths)
}

func readInput() {
	input, err := os.ReadFile("input.txt")
	errorCheck("Error reading the input file", err)
	sanitizedInput := strings.ReplaceAll(string(input), "\r\n", "\n")
	data := strings.Split(string(sanitizedInput), "\n")
	antsNumber = data[0]
	for i := 0; i < len(data); i++ {
		splittedData := strings.Split(data[i], " ")
		if len(splittedData) == 3 {
			roomsAndTunnels[splittedData[0]] = []string{}
		}
		if data[i] == "##start" {
			fmt.Println(data[i])
			splittedData = strings.Split(data[i+1], " ")
			startRoom = &splittedData[0]
		}
		if data[i] == "##end" && startRoom != nil {
			splittedData = strings.Split(data[i+1], " ")
			endRoom = &splittedData[0]
		}
		if endRoom != nil {
			splittedData = strings.Split(data[i], "-")
			if len(splittedData) == 2 {
				for i := range roomsAndTunnels {
					if i == splittedData[0] {
						roomsAndTunnels[i] = append(roomsAndTunnels[i], splittedData[1])
					}
					if i == splittedData[1] {
						roomsAndTunnels[i] = append(roomsAndTunnels[i], splittedData[0])
					}
				}
			}
		}
	}
	checkInput()
}

func checkInput() {
	if startRoom == nil || endRoom == nil || len(roomsAndTunnels) == 0 {
		fmt.Println("Input data is corrupted")
		os.Exit(1)
	}
	for rooms := range roomsAndTunnels {
		if rooms[0] == '#' || rooms[0] == 'L' {
			fmt.Println("Wrong naming of the rooms")
			os.Exit(1)
		}
	}
}

// var iteration int
func searchForPath(currentRoom string, visitedRooms, possiblePath []string) [][]string {
	var succesfulPaths [][]string
	for i := 0; i < len(roomsAndTunnels[currentRoom]); i++ {
		visitedRooms = append(visitedRooms, currentRoom)
		possiblePath = append([]string(nil), possiblePath...)

		if currentRoom == *endRoom {
			possiblePath = append(possiblePath, currentRoom)
			succesfulPaths = append(succesfulPaths, possiblePath)
			//break
			continue
		}

		nextRoom := roomsAndTunnels[currentRoom][i]
		if !isVisited(nextRoom, visitedRooms) {
			possiblePath = append(possiblePath, currentRoom)
			visitedRooms = append(visitedRooms, nextRoom)
			results := searchForPath(nextRoom, visitedRooms, possiblePath)
			//fmt.Printf("iteration number: %v, Current Room: %s, Path: %v, Visited: %v\n", iteration, currentRoom, possiblePath, visitedRooms)
			//iteration++
			succesfulPaths = append(succesfulPaths, results...)
			if len(visitedRooms) > 0 && len(possiblePath) > 0 {
				visitedRooms = visitedRooms[:len(visitedRooms)-1]
				possiblePath = possiblePath[:len(possiblePath)-1]
			}
		}
	}
	return isDuplicate(succesfulPaths)
}

func isVisited(currentRoom string, visitedRooms []string) bool {
	for _, room := range visitedRooms {
		if room == currentRoom {
			return true
		}
	}
	return false
}

func isDuplicate(succesfulPaths [][]string) [][]string {
	var uniquePaths [][]string
	for _, path := range succesfulPaths {
		keepTrack := make(map[string]bool)
		isDublicate := false
		for _, room := range path {
			if keepTrack[room] {
				isDublicate = true
				break
			}
			keepTrack[room] = true
		}
		if !isDublicate {
			uniquePaths = append(uniquePaths, path)
		}
	}
	return uniquePaths
}

func intersects(set [][]string, path []string) bool {
	for _, rooms := range set {
		for _, room1 := range rooms[1 : len(rooms)-1] {
			for _, room2 := range path[1 : len(path)-1] {
				if room1 == room2 {
					return true
				}
			}
		}
	}
	return false
}

//*** this function was made with assistance from Marcus
func setsOfPaths(paths, curSet [][]string, index int, sets *[][][]string) {
	curSet = append(curSet, paths[index])
	paths = paths[index+1:]

	nonIntersecting := [][]string{}
	for _, path := range paths {
		if !intersects(curSet, path) {
			nonIntersecting = append(nonIntersecting, path)
		}
	}

	if len(nonIntersecting) == 0 {
		*sets = append(*sets, curSet)
	}

	for i := range nonIntersecting {
		setsOfPaths(nonIntersecting, curSet, i, sets)
	}
}

// ***While doing func setsOfPaths, Marcus based on this function below. Is there still a need to keep this code or remove it?
func isIntersect(paths [][]string) map[int][][]string {
	pathsCombinations := make(map[int][][]string)
	index := 1

	for i := 0; i < len(paths); i++ {
		noJamPaths := [][]string{paths[i]}
		for j := 0; j < len(paths); j++ {
			if i != j {
				intersect := false
				for _, path := range noJamPaths {
					for k := 0; k < len(path); k++ {
						for l := 0; l < len(paths[j]); l++ {
							if path[k] == paths[j][l] && path[k] != *startRoom && path[k] != *endRoom {
								intersect = true
								break
							}
						}
						if intersect {
							break
						}
					}
					if intersect {
						break
					}
				}
				if !intersect {
					noJamPaths = append(noJamPaths, paths[j])
				}
			}
		}
		if len(noJamPaths) > 1 {
			pathsCombinations[index] = noJamPaths
			index++
		}
	}
	return pathsCombinations
}

/*func moveAnts(pathsCombinations map[int][][]string) {
	ants, err := strconv.Atoi(antsNumber)
	errorCheck("Error reading the ant number", err)

	var shortestFlow string
	var shortestCombination [][]string
	shortestLength := 10000

	for index, paths := range pathsCombinations {
		fmt.Printf("\nCombination %d:\n", index)

		antPositions := make([]int, ants)
		antPaths := make([]int, ants)
		for i := 0; i < ants; i++ {
			antPaths[i] = i % len(paths)
		}

		flow := ""
		finishedAnts := 0

		for finishedAnts < ants {
			moves := []string{}
			occupiedRooms := make(map[string]bool)

			for ant := 0; ant < ants; ant++ {
				if antPositions[ant] == -1 {
					continue // Skip finished ants
				}

				path := paths[antPaths[ant]]
				currentPos := antPositions[ant]

				if currentPos == len(path)-1 {
					antPositions[ant] = -1
					finishedAnts++
					continue // Skip finished ants
				}

				nextRoom := path[currentPos+1]

				if !occupiedRooms[nextRoom] || nextRoom == *endRoom {
					moves = append(moves, fmt.Sprintf("L%d-%s", ant+1, nextRoom))
					antPositions[ant]++
					if nextRoom != *endRoom {
						occupiedRooms[nextRoom] = true
					}
				}
			}

			if len(moves) > 0 {
				flow += strings.Join(moves, " ") + "\n"
			}
		}

		if strings.Count(flow, "\n") < shortestLength {
			shortestLength = strings.Count(flow, "\n")
			shortestFlow = flow
			shortestCombination = paths
		}
	}

	fmt.Println("\nShortest flow:")
	fmt.Print(shortestFlow)
	fmt.Println("Combination of paths for the shortest flow:")
	for _, path := range shortestCombination {
		fmt.Println(path)
	}
}*/

//showing all the diff combinations and the flow of each

func moveAnts(pathsCombinations map[int][][]string) {
	ants, _ := strconv.Atoi(antsNumber)
	var shortestFlow string
	var shortestCombination [][]string
	var shortestLength = 10000 // Initialize with a large number

	for index, paths := range pathsCombinations {
		fmt.Printf("\nCombination %d:\n", index)

		// Print the combination of paths
		fmt.Println("Combination of paths:")
		for _, path := range paths {
			fmt.Println(path)
		}

		antPositions := make([]int, ants) // Position of each ant in their path, starts at -1 to indicate in start room
		for i := range antPositions {
			antPositions[i] = 1 // Start at index 1 since 0 is the start room
		}

		var flow string
		for {
			hasMove := false
			moves := []string{}
			occupiedRooms := make(map[string]bool)

			// Try to move each ant
			for ant := 0; ant < ants; ant++ {
				pathIndex := ant % len(paths)
				path := paths[pathIndex]

				if antPositions[ant] < len(path) {
					nextRoom := path[antPositions[ant]]
					// Skip if room is already occupied (unless it's the end room)
					if nextRoom != *endRoom && occupiedRooms[nextRoom] {
						continue
					}

					hasMove = true
					moves = append(moves, fmt.Sprintf("L%d-%s", ant+1, nextRoom))
					occupiedRooms[nextRoom] = true
					antPositions[ant]++
				}
			}

			if !hasMove {
				break
			}
			flow += strings.Join(moves, " ") + "\n"
		}

		// Print the flow for the current combination
		fmt.Println("Flow:")
		fmt.Print(flow)

		// Check if the current flow is shorter than the shortest flow found so far
		if strings.Count(flow, "\n") < shortestLength {
			shortestLength = strings.Count(flow, "\n")
			shortestFlow = flow
			shortestCombination = paths
		}
	}

	// Print the shortest flow and its combination
	fmt.Println("\nShortest flow:")
	fmt.Print(shortestFlow)
	fmt.Println("Combination of paths for the shortest flow:")
	for _, path := range shortestCombination {
		fmt.Println(path)
	}
}

/*func moveAnts(pathsCombinations map[int][][]string) {
	ants, err := strconv.Atoi(antsNumber)
	errorCheck("Error converting ants number", err)

	for _, combination := range pathsCombinations {
		// Assign each ant to a specific path
		paths := combination
		antPositions := make([]int, ants) // Tracks each ant's current position in its path

		for step := 0; ; step++ { // Loop until all ants have finished
			allFinished := true

			for ant := 0; ant < ants; ant++ {
				pathIndex := ant % len(paths) // Assign ants to paths cyclically
				path := paths[pathIndex]

				if antPositions[ant] < len(path) { // Check if the ant still has rooms to move to
					allFinished = false
					fmt.Printf("L%v-%v ", ant+1, path[antPositions[ant]])
					antPositions[ant]++ // Move the ant to the next room
				}
			}

			fmt.Println() // Newline after each turn

			if allFinished { // Stop if all ants have finished their paths
				break
			}
		}
	}
}*/

func isFound(possiblePath []string) bool {
	for _, room := range possiblePath {
		if room == *endRoom {
			return true
		}
	}
	return false
}

func errorCheck(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}
