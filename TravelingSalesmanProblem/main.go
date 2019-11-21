package main

import (
	ga "TravelingSalesmanProblem/population"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"src/github.com/paulmach/go.geo"
	"time"
)

var (
	numbereOfTowns = flag.Int("n", 4, "number of towns to go through")
	cities geo.PointSet
	allCityDistances map[[2]int]float64
	bestEverPopulation ga.Population
	sameBestCount int
	populationsCount int
)

func GeneratePoints() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *numbereOfTowns; i++ {
		//cities.Push(geo.NewPoint(rand.Float64(), rand.Float64()))
		cities.Push(geo.NewPoint(float64(rand.Intn(100)), float64(rand.Intn(100))))
	}
}

func CalculateAllDistances() {
	cityDistances := map[[2]int]float64{}

	for i := 0; i < len(cities) - 1; i++ {
		pointA := *cities.GetAt(i)
		for j := i + 1; j < len(cities); j++ {
			pointB := *cities.GetAt(j)
			curr := [2]int{i, j}
			cityDistances[curr] = ga.PointDistance(pointA, pointB)
		}
	}

	allCityDistances = cityDistances
}

func GeneratePopulations() ([]ga.Population, float64, float64) {
	populations := make([]ga.Population, populationsCount)
	maxDistance, minDistance := 0.0, math.Inf(1)

	for i := 0; i < populationsCount; i++ {
		perm := rand.Perm(cities.Length())
		newPopulation := ga.NewPopulation(perm, allCityDistances)
		currDistance := newPopulation.GetDistance()
		populations[i] = *newPopulation

		if bestEverPopulation.GetDistance() > currDistance {
			bestEverPopulation = *newPopulation
		}

		if maxDistance < currDistance {
			maxDistance = currDistance
		}

		if minDistance > currDistance {
			minDistance = currDistance
		}
	}

	return populations, minDistance, maxDistance
}

func Rank(population []ga.Population) []ga.Population {
	sort.Slice(population, func(i, j int) bool {
		return population[i].GetFitness() > population[j].GetFitness()
	})
	return population
}

func Mutate(individual []int) []int {
	geneA, geneB := 0, 0
	for geneA == geneB {
		geneA = rand.Intn(len(individual))
		geneB = rand.Intn(len(individual))
	}
	individual[geneA], individual[geneB] = individual[geneB], individual[geneA]

	return individual
}

func Cross(parentA ga.Population, parentB ga.Population) (*ga.Population, *ga.Population ) {
	cutIndex := 0
	parentARoute := parentA.GetRoute()
	parentBRoute := parentB.GetRoute()

	for cutIndex == 0 {
		cutIndex = rand.Intn(len(parentARoute) - 1)
	}

	child1Has := map[int]bool{}
	child1Route := make([]int, len(parentARoute))
	child2Has := map[int]bool{}
	child2Route := make([]int, len(parentARoute))

	for i := 0; i < cutIndex; i++ {
		child1Route[i] = parentARoute[i]
		child1Has[parentARoute[i]] = true

		child2Route[i] = parentBRoute[i]
		child2Has[parentBRoute[i]] = true
	}
	child1Index := cutIndex
	child2Index := cutIndex

	for i := cutIndex; i < len(parentBRoute); i++ {
		if !child1Has[parentBRoute[i]] {
			child1Route[child1Index] = parentBRoute[i]
			child1Has[parentBRoute[i]] = true
			child1Index += 1
		}

		if !child2Has[parentARoute[i]] {
			child2Route[child2Index] = parentARoute[i]
			child2Has[parentARoute[i]] = true
			child2Index += 1
		}
	}

	for i := 0; i < cutIndex; i++ {
		if !child1Has[parentBRoute[i]] {
			child1Route[child1Index] = parentBRoute[i]
			child1Has[parentBRoute[i]] = true
			child1Index += 1
		}

		if !child2Has[parentARoute[i]] {
			child2Route[child2Index] = parentARoute[i]
			child2Has[parentARoute[i]] = true
			child2Index += 1
		}
	}

	child1 := ga.NewPopulation(child1Route, allCityDistances)
	child2 := ga.NewPopulation(child2Route, allCityDistances)

	return child1, child2
}

func Breed(parents []ga.Population) []ga.Population {
	children := make([]ga.Population, populationsCount)
	parentsCount := len(parents)
	crossedParents := map[[2]int]bool{}

	for i := 0; i < parentsCount; i++ {
		parentIndexes := [2]int{}
		for parentIndexes[0] == parentIndexes[1]{
			parentIndexes[0] = rand.Intn(parentsCount)
			parentIndexes[1] = rand.Intn(parentsCount)

			if parentIndexes[0] > parentIndexes[1] {
				parentIndexes[0], parentIndexes[1] = parentIndexes[1], parentIndexes[0]
			}

			if crossedParents[parentIndexes] {
				parentIndexes[0] = parentIndexes[1]
			}
		}

		child1, child2 := Cross(parents[parentIndexes[0]], parents[parentIndexes[1]])
		crossedParents[parentIndexes] = true

		children[i] = *child1
		children[i + parentsCount] = *child2
	}

	toMutate := int(populationsCount / 5)
	mutated := map[int]bool{}
	for i := 0; i < toMutate; i++ {
		childToMutate := rand.Intn(populationsCount - 1)
		if !mutated[childToMutate] {
			mutatedChild := Mutate(children[childToMutate].GetRoute())
			children[childToMutate].SetRoute(mutatedChild)
			children[childToMutate].CalculateDistance(allCityDistances)
			mutated[childToMutate] = true
		} else {
			i -= 1
		}
	}

	bestNow := ga.Population{}
	bestNow.SetDistance(math.Inf(1))
	bestIndex := -1

	for i := 0; i < populationsCount; i++ {
		if bestNow.GetDistance() > children[i].GetDistance() {
			bestNow = *ga.NewPopulation(children[i].GetRoute(), allCityDistances)
			bestIndex = i
		}
	}

	if bestNow.GetDistance() > bestEverPopulation.GetDistance() {
		children[bestIndex] = parents[0]
	} else if bestNow.GetDistance() < bestEverPopulation.GetDistance() {
		bestEverPopulation = bestNow
		sameBestCount = 0
	} else {
		sameBestCount += 1
	}

	return children
}


func main() {
	flag.Parse()
	populationsCount = 2 * (*numbereOfTowns)
	bestEverPopulation.SetDistance(math.Inf(1))

	// 1. create random points on the 2D
	GeneratePoints()
	CalculateAllDistances()
	fmt.Println("Cities: ", cities)

	// 2. generate populations
	populations, minDistance, maxDistance := GeneratePopulations()
	totalFitness := 0.0
	// Map all fitness values between 0 and 1
	for i := 0; i < len(populations); i++ {
		totalFitness += populations[i].CalculateFitness(minDistance, maxDistance)
	}

	// normalize the fitness values to a probability between 0 and 1
	for i := 0; i < len(populations); i++ {
		populations[i].NormalizeFitness(totalFitness)
		if populations[i].GetDistance() == bestEverPopulation.GetDistance() {
			bestEverPopulation.SetFitness(populations[i].GetFitness())
		}
	}
	populations = Rank(populations)

	generation := 0
	for sameBestCount < 5 {
		populations = Breed(populations[:populationsCount/2])

		totalFitness := 0.0
		// Map all fitness values between 0 and 1
		for i := 0; i < len(populations); i++ {
			totalFitness += populations[i].CalculateFitness(minDistance, maxDistance)
		}

		// normalize the fitness values to a probability between 0 and 1
		for i := 0; i < len(populations); i++ {
			populations[i].NormalizeFitness(totalFitness)
			if populations[i].GetDistance() == bestEverPopulation.GetDistance() {
				bestEverPopulation.SetFitness(populations[i].GetFitness())
			}
		}
		populations = Rank(populations)
		generation += 1

		if generation == 1 || generation == 5 || generation == 10 || generation == 15 {
			fmt.Printf("Best in generation %v: %v; with distance: %v\n", generation,
				bestEverPopulation.GetRoute(), bestEverPopulation.GetDistance())
		}
	}

	fmt.Printf("Last best in generation %v: %v; with distance: %v\n", generation,
		bestEverPopulation.GetRoute(), bestEverPopulation.GetDistance())
}
