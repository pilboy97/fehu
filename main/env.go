package main

import (
	"core"
)

type Env struct {
}

func NewEnv() *Env {
	return &Env{}
}
func (e *Env) Path() string {
	return core.DBPath
}
func (e *Env) Open(path string) error {
	core.DBPath = path
	return core.Open(path)
}
func (e *Env) Close() {
	core.DBPath = ""
	core.Close()
}
func (e *Env) Code() string {
	return core.Code
}
func (e *Env) SetCode(str string) {
	core.Code = str
}
