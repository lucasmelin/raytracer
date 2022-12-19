package geometry

type Rnd interface {
	Float64() float64
}

func FloatInRange(rnd Rnd, min float64, max float64) float64 {
	return min + rnd.Float64()*(max-min)
}
