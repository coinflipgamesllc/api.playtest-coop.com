package app

const fromAddress = "no-reply@playtest-coop.com"

// func (s *Server) welcomeEmail(email, name, verificationID string) error {
// 	templateData := struct {
// 		Name string
// 		URL  string
// 	}{
// 		Name: name,
// 		URL:  s.hostname + "/v1/auth/verify-email/" + verificationID,
// 	}

// 	tpl := s.templates["email/welcome"]
// 	buf := new(bytes.Buffer)
// 	if err := tpl.Execute(buf, templateData); err != nil {
// 		return err
// 	}

// 	return s.send(email, "Welcome to Playtest Co-op!", buf.String())
// }

// func (s *Server) verifyEmail(email, name, verificationID string) error {
// 	templateData := struct {
// 		Name string
// 		URL  string
// 	}{
// 		Name: name,
// 		URL:  s.hostname + "/v1/auth/verify-email/" + verificationID,
// 	}

// 	tpl := s.templates["email/verify-email"]
// 	buf := new(bytes.Buffer)
// 	if err := tpl.Execute(buf, templateData); err != nil {
// 		return err
// 	}

// 	return s.send(email, "Verify your email", buf.String())
// }

// func (s *Server) send(toAddress, subject, body string) error {
// 	message := s.mail.NewMessage(
// 		fromAddress,
// 		subject,
// 		"",
// 		toAddress,
// 	)
// 	message.SetHtml(body)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
// 	defer cancel()

// 	_, _, err := s.mail.Send(ctx, message)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
