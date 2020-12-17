package events

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/pubsub"
	"go.uber.org/zap"
)

// EventHandler routes domain events to the proper handler
type EventHandler struct {
	MailService *app.MailService
	Logger      *zap.SugaredLogger
}

// ListenForEvents creates the channels for recieving domain events and sets up the handlers
func (h *EventHandler) ListenForEvents() {
	userCreated := make(chan pubsub.Message)
	pubsub.Instance.Subscribe("User/Created", userCreated)

	userEmailUnverified := make(chan pubsub.Message)
	pubsub.Instance.Subscribe("User/EmailUnverified", userEmailUnverified)

	userPasswordResetRequested := make(chan pubsub.Message)
	pubsub.Instance.Subscribe("User/PasswordResetRequested", userPasswordResetRequested)

	for {
		select {
		case evt := <-userCreated:
			go h.userCreated(evt)
		case evt := <-userEmailUnverified:
			go h.userEmailUnverified(evt)
		case evt := <-userPasswordResetRequested:
			go h.userPasswordResetRequested(evt)
		}
	}
}

func (h *EventHandler) userCreated(msg pubsub.Message) {
	h.Logger.Infof("[User/Created]: %+v\n\n", msg)

	data := msg.Data.(map[string]interface{})

	err := h.MailService.SendWelcomeEmail(data["email"].(string), data["name"].(string), data["verificationID"].(string))
	if err != nil {
		h.Logger.Error(err)
	}
}

func (h *EventHandler) userEmailUnverified(msg pubsub.Message) {
	h.Logger.Infof("[User/EmailUnverified]: %+v\n\n", msg)

	data := msg.Data.(map[string]interface{})

	err := h.MailService.SendVerifyEmail(data["email"].(string), data["name"].(string), data["verificationID"].(string))
	if err != nil {
		h.Logger.Error(err)
	}
}

func (h *EventHandler) userPasswordResetRequested(msg pubsub.Message) {
	h.Logger.Infof("[User/PasswordResetRequested]: %+v\n\n", msg)

	data := msg.Data.(map[string]interface{})

	err := h.MailService.SendPasswordResetEmail(data["email"].(string), data["name"].(string), data["otp"].(string))
	if err != nil {
		h.Logger.Error(err)
	}
}
