package main

import (
	"core"
)

type Env struct {
	code string
}

func NewEnv() *Env {
	return &Env{
		code: "KRW",
	}
}
func (e *Env) Path() string {
	return core.DBPath
}
func (e *Env) Open(path string) {
	core.DBPath = path
	core.Open(path)
}
func (e *Env) Close() {
	core.DBPath = ""
	core.Close()
}
