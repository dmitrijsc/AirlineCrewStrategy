package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// City represents a Latvian city where flights operate.
type City string

const (
	Riga       City = "Riga"
	Daugavpils      = "Daugavpils"
	Liepaja         = "Liepāja"
	Jelgava         = "Jelgava"
	Ventspils       = "Ventspils"
)

// Airplane represents an airplane.
type Airplane string

const (
	Plane1 Airplane = "Airplane 1"
	Plane2          = "Airplane 2"
	Plane3          = "Airplane 3"
	Plane4          = "Airplane 4"
)

// Pilot represents a pilot.
type Pilot struct {
	Name            string
	CurrentCity     City     // Track the current location of the pilot
	CurrentAirplane Airplane // Track the current airplane of the pilot
}

// Flight represents a flight between two cities, operated by an airplane and requiring a pilot.
type Flight struct {
	FromCity   City
	ToCity     City
	Airplane   Airplane
	FirstPilot Pilot
	Time       time.Time
}

// GenerateCities returns the 5 largest cities in Latvia.
func GenerateCities() []City {
	return []City{Riga, Daugavpils, Liepaja, Jelgava, Ventspils}
}

func GenerateAirplanes() []Airplane {
	return []Airplane{Plane1, Plane2, Plane3, Plane4}
}

// GeneratePilots returns a list of 10 pilots starting at Riga.
func GeneratePilots(count int) []Pilot {
	pilots := make([]Pilot, count)
	for i := 0; i < count; i++ {
		pilots[i] = Pilot{Name: fmt.Sprintf("Pilot %d", i+1), CurrentCity: Riga}
	}
	return pilots
}

// GenerateFlights creates synthetic flight schedules between cities while respecting continuity and time order.
func GenerateFlights(cities []City, airplanes []Airplane, pilots []Pilot, numFlights int) []Flight {

	flights := make([]Flight, 0, numFlights)

	// Start airplane locations in Riga or Liepaja
	planeLocations := map[Airplane]City{}

	for i, airplane := range airplanes {
		if i%2 == 0 {
			planeLocations[airplane] = Riga
		} else {
			planeLocations[airplane] = Liepaja
		}
	}

	// Helper function to pick a city, with 70% chances favoring Riga and Liepaja
	chooseCity := func(exclude City) City {

		candidate := cities[rand.Intn(len(cities))]

		if rand.Float64() < 0.7 {
			if rand.Float64() < 0.5 {
				candidate = Riga
			}
			candidate = Liepaja
		}

		for candidate == exclude {
			candidate = cities[rand.Intn(len(cities))]
		}
		return candidate
	}

	// Helper function to generate time
	baseTime := time.Now()

	// Assign flights while maintaining airplane continuity and respecting time order
	for i := 0; i < numFlights; i++ {
		airplane := airplanes[rand.Intn(len(airplanes))]
		fromCity := planeLocations[airplane]
		toCity := chooseCity(fromCity)

		firstPilot := pilots[rand.Intn(len(pilots))]

		flightTime := baseTime.Add(time.Duration(i*2) * time.Hour) // Each flight happens 2 hours after the previous one
		flight := Flight{
			FromCity:   fromCity,
			ToCity:     toCity,
			Airplane:   airplane,
			FirstPilot: firstPilot,
			Time:       flightTime,
		}
		flights = append(flights, flight)

		// Update the airplane's new location to the destination city
		planeLocations[airplane] = toCity
	}

	return flights
}

// Check if the pilot starts or ends in Riga or Liepaja, or it's the first or last flight in the day.
func isValidFlight(flight Flight, isFirstFlight, isLastFlight bool) bool {

	allowedCities := []City{Riga, Liepaja}

	if isFirstFlight || isLastFlight {
		return true
	}
	// Check if the flight starts or ends in Riga or Liepaja
	for _, city := range allowedCities {
		if flight.FromCity == city || flight.ToCity == city {
			return true
		}
	}
	return false
}

// Evaluate the schedule based on flight rules and pilot continuity.
// Adds a penalty for using more pilots, rewards for fewer.
func evaluate(flights []Flight, pilots []Pilot) int {
	conflicts := 0

	// We aim to minimize conflicts. Calculate rule violations.
	for i, flight := range flights {
		isFirstFlight := (i == 0)
		isLastFlight := (i == len(flights)-1)

		// Check if the pilot respects the start/end rules
		if !isValidFlight(flight, isFirstFlight, isLastFlight) {
			conflicts++
		}
	}

	// Ensure pilot continuity
	pilotLocations := make(map[string]City)
	pilotAirplanes := make(map[string]Airplane)

	for _, pilot := range pilots {
		pilotLocations[pilot.Name] = pilot.CurrentCity // Initialize the pilots' starting positions
		pilotAirplanes[pilot.Name] = pilot.CurrentAirplane
	}

	for i, flight := range flights {

		pilotName := flight.FirstPilot.Name
		if currentCity, exists := pilotLocations[pilotName]; exists {
			// Check if pilot can start from the current city's location
			if currentCity != flight.FromCity {
				conflicts += 100 // Penalize if the pilot cannot start from their last city
			}

			// Update the pilot's location to the destination of the flight
			pilotLocations[pilotName] = flight.ToCity
		} else {
			pilotLocations[pilotName] = flight.ToCity
		}

		// Check if the pilot is using the same airplane
		if currentAirplane, exists := pilotAirplanes[pilotName]; exists {

			if currentAirplane != flight.Airplane {
				conflicts += 50 // Penalize more heavily if the pilot is switching airplanes
			} else {
				conflicts -= 10 // Reward for staying on the same airplane
			}
			// Update the pilot's airplane to the one used in the flight
			pilotAirplanes[pilotName] = flight.Airplane
		} else {
			pilotAirplanes[pilotName] = flight.Airplane
		}

		// Penalize pilot for jumping between airports unnecessarily (i.e., without flying)
		if i > 0 && flights[i-1].FirstPilot.Name == flight.FirstPilot.Name && flights[i-1].ToCity != flight.FromCity {
			conflicts += 250 // High penalty for "jumping" between airports
		}
	}

	// Reward fewer pilots, penalize more
	pilotFlights := make(map[string]int)
	for _, flight := range flights {
		pilotFlights[flight.FirstPilot.Name]++
	}

	numPilotsUsed := len(pilotFlights)

	// Penalize large number of pilots used
	conflicts += numPilotsUsed * 100

	return conflicts
}

// TabuSearch represents a basic implementation of the Tabu Search algorithm while respecting time order.
func TabuSearch(flights []Flight, pilots []Pilot, maxIterations int, tabuTenure int) []Flight {
	bestSolution := flights
	currentSolution := flights
	tabuList := make(map[string]int) // Maps pilot name to tenure in tabu list
	iteration := 0

	// Helper function to generate a key for the flight for tabu check
	generateTabuKey := func(flight Flight) string {
		return flight.FirstPilot.Name + string(flight.Airplane) + flight.Time.Format(time.RFC3339)
	}

	// Make a move by swapping pilots between flights with adjacent or same times
	makeMove := func(flights []Flight) []Flight {
		newSolution := make([]Flight, len(flights))
		copy(newSolution, flights)

		// Randomly select two adjacent flights to swap pilots
		i := rand.Intn(len(flights) - 1) // Ensuring that i and i+1 are valid
		j := i + 1

		// Swap the pilots in the two selected flights
		newSolution[i].FirstPilot, newSolution[j].FirstPilot = newSolution[j].FirstPilot, newSolution[i].FirstPilot

		return newSolution
	}

	// Tabu Search main loop
	for iteration < maxIterations {
		iteration++

		// Generate a neighboring solution by making a random move
		newSolution := makeMove(currentSolution)

		// Evaluate the new solution
		newConflicts := evaluate(newSolution, pilots)

		// Check if the new solution is better or if it's tabu
		for _, flight := range newSolution {
			tabuKey := generateTabuKey(flight)
			if tenure, found := tabuList[tabuKey]; found && tenure > 0 {
				newConflicts += 1000 // Add a large penalty if the move is tabu
			}
		}

		// If the new solution is better, update the current solution
		if newConflicts < evaluate(currentSolution, pilots) {
			currentSolution = newSolution

			// Update the tabu list
			for _, flight := range currentSolution {
				tabuKey := generateTabuKey(flight)
				tabuList[tabuKey] = tabuTenure
			}
		}

		// Update the best solution found so far
		if evaluate(currentSolution, pilots) < evaluate(bestSolution, pilots) {
			bestSolution = currentSolution
		}

		// Decrease the tabu tenure for all entries in the tabu list
		for key := range tabuList {
			if tabuList[key] > 0 {
				tabuList[key]--
			}
		}

		// Early exit if there are no conflicts
		if evaluate(bestSolution, pilots) == 0 {
			break
		}
	}

	return bestSolution
}

func main() {

	seed := rand.Intn(40) + 10
	rand.Seed(int64(seed)) // Set the seed for reproducible results

	cities := GenerateCities()
	airplanes := GenerateAirplanes()
	pilots := GeneratePilots(8)

	flights := GenerateFlights(cities, airplanes, pilots, seed)

	// Create a file to save the results
	filename := fmt.Sprintf("results_seed_%d_%d.txt", seed, len(pilots))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Function to write to both console and file
	writeOutput := func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
		fmt.Fprintf(file, format, a...)
	}

	writeOutput("Generated Flights:\n")

	for _, airplane := range airplanes {

		writeOutput("Airplane: %s\n", airplane)

		for _, flight := range flights {
			if airplane == flight.Airplane {
				writeOutput("\t%s -> %s \t| %s \t| %s \t| %v\n", flight.FromCity, flight.ToCity, flight.Airplane, flight.FirstPilot.Name, flight.Time.Format(time.RFC822))
			}
		}
	}

	initialConflicts := evaluate(flights, pilots)
	writeOutput("Initial Conflicts: %d\n", initialConflicts)

	// Perform Tabu Search to optimize pilot assignment
	maxIterations := 50000
	tabuTenure := 500

	start := time.Now()
	optimizedFlights := TabuSearch(flights, pilots, maxIterations, tabuTenure)
	duration := time.Since(start)

	writeOutput("\nOptimized Flights:\n")

	for _, pilot := range pilots {

		writeOutput("Pilot: %s\n", pilot.Name)

		for _, flight := range optimizedFlights {
			if pilot == flight.FirstPilot {
				writeOutput("\t %s -> %s \t| %s \t| %v\n", flight.FromCity, flight.ToCity, flight.Airplane, flight.Time.Format(time.RFC822))
			}
		}
	}

	optimizedConflicts := evaluate(optimizedFlights, pilots)
	writeOutput("TabuSearch execution time: %v\n", duration)
	writeOutput("Optimized Conflicts: %d\n", optimizedConflicts)
}
