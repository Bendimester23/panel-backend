package controllers

import "log"

type ServerController struct {
}

func (s *ServerController) NewServer(name string) {
	log.Println(name)
}
