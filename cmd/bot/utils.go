package main

import "io/ioutil"

func loadDCA() ([]byte, error) {
	buf, err := ioutil.ReadFile("assets/zaht.dca")
	if err != nil {
		return nil, err
	}

	return buf, nil
}
