package blockchain

import "errors"

func NextRound(n int, randoms []int) ([]int, error) {
	if len(randoms) > n {
		return nil, errors.New("error array length")
	}
	sum := 0
	for _, v := range randoms {
		sum += v
	}
	value := [n]int{}
	for i := 0; i < len(randoms); i++ {
		//v := sum % randoms[i]
	}
	return value[:], nil
}
