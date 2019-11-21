package population

import (
	"math"
	"src/github.com/paulmach/go.geo"
)

type Population struct {
	route []int
	distance float64
	fitness float64
}

func NewPopulation(path []int, allCityDistances map[[2]int]float64) *Population {
	population := &Population{
		route: path,
		distance: math.Inf(1),
		fitness: 0,
	}

	population.CalculateDistance(allCityDistances)

	return population
}

func (population *Population) GetRoute() []int {
	return population.route
}

func (population *Population) SetRoute(route []int) {
	population.route = route
}

func (population *Population) GetDistance() float64 {
	return population.distance
}

func (population *Population) SetDistance(initialDistance float64) {
	population.distance = initialDistance
}

func (population *Population) GetFitness() float64 {
	return population.fitness
}

func (population *Population) SetFitness(fitness float64) {
	population.fitness = fitness
}

func (population *Population) CalculateDistance(allCityDistances map[[2]int]float64) float64 {
	distance := 0.0
	var pointA, pointB int
	for i := 0; i < len(population.route); i++ {
		if i == len(population.route) - 1 {
			pointA = population.route[0]
			pointB = population.route[i]
			if pointA > pointB {
				pointA, pointB = pointB, pointA
			}
		} else {
			pointA = population.route[i]
			pointB = population.route[i + 1]
			if pointA > pointB {
				pointA, pointB = pointB, pointA
			}
		}
		points := [2]int{pointA, pointB}
		distance += allCityDistances[points]
	}
	population.distance = distance

	return population.distance
}

func (population *Population) CalculateFitness(minDistance float64, maxDistance float64) float64 {
	fitness := (maxDistance - population.distance) / (maxDistance - minDistance)
	population.fitness = fitness

	return population.fitness
}

func (population *Population) NormalizeFitness(total float64) {
	population.fitness /= total
}

func PointDistance(pointA geo.Point, pointB geo.Point) float64{
	d1 := pointB.X() - pointA.X()
	d2 := pointB.Y() - pointA.Y()

	return math.Sqrt(d1*d1 + d2*d2)
}