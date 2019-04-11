package userlib

import "io/ioutil"

type fileReader func(string)([]byte, error)

const (
	FILEERRORCODE = 404
	FILEERRORMSG = "File Read Error"
	TIMEOUTERRORCODE = 408
	SUCCESSCODE = 200
)

var F fileReader = func(filename string)(data []byte, err error){
	data, err = ioutil.ReadFile(filename)
	return
}

func ReadFile(filename string)(data []byte, err error){
	data, err = F(filename)
	return
}


func ReplaceReadFile(newfunc func(string)([]byte, error)){
	F = newfunc
}