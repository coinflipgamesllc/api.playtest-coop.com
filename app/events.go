package app

import (
	"fmt"
	"log"

	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/pubsub"
)

func (s *Server) listenForEvents() {
	userCreated := make(chan pubsub.Message)
	pubsub.Instance.Subscribe("User/Created", userCreated)

	userEmailUnverified := make(chan pubsub.Message)
	pubsub.Instance.Subscribe("User/EmailUnverified", userEmailUnverified)

	for {
		select {
		case evt := <-userCreated:
			go s.userCreated(evt)
		case evt := <-userEmailUnverified:
			go s.userEmailUnverified(evt)
		}
	}
}

func (s *Server) userCreated(msg pubsub.Message) {
	log.Printf("[User/Created]: %+v\n\n", msg)

	data := msg.Data.(map[string]interface{})

	err := s.welcomeEmail(data["email"].(string), data["name"].(string), data["verificationID"].(string))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func (s *Server) userEmailUnverified(msg pubsub.Message) {
	log.Printf("[User/EmailUnverified]: %+v\n\n", msg)

	data := msg.Data.(map[string]interface{})

	err := s.verifyEmail(data["email"].(string), data["name"].(string), data["verificationID"].(string))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}
