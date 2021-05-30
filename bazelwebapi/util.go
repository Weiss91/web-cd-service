package main

import "io/ioutil"

func copy(src string, dst string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
