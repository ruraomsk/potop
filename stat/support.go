package stat

func GetDescription() (int, int) {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	return Statistics.diaps, Statistics.counts
}
func GetCountValues() []int {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	ret := make([]int, 0)
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		for _, v := range r.CountValues.Value {
			ret = append(ret, v)
		}
	}
	return ret
}
func GetSpeedValues() []int {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	ret := make([]int, 0)
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		for _, v := range r.SpeedValues.Value {
			ret = append(ret, v)
		}
	}
	return ret
}
func ClearCountValues() {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		r.CountValues.clear(Statistics.diaps)
	}
}
func ClearSpeedValues() {
	Statistics.mutex.Lock()
	defer Statistics.mutex.Unlock()
	for i := 0; i < Statistics.counts; i++ {
		r, is := Statistics.chanels[i]
		if !is {
			continue
		}
		r.SpeedValues.clear(Statistics.diaps)
	}
}
