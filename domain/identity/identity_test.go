package identity

import "github.com/DATA-DOG/godog"

var (
	testedEmail    string
	testedPassword string
)

func iHaveTheEmail(email string) error {
	testedEmail = email
	return nil
}

func thePassword(password string) error {
	testedPassword = password
	return nil
}

func theEmailIsNotAlreadyRegistered() error {
	return godog.ErrPending
}

func iRegisterANewAccount() error {
	return godog.ErrPending
}

func myAccountShouldHaveBeenCreated() error {
	return godog.ErrPending
}

func theEmailIsAlreadyRegistered() error {
	return godog.ErrPending
}

func itShouldComplainsAboutTheEmailBeingAlreadyTaken() error {
	return godog.ErrPending
}

func iAuthenticate() error {
	return godog.ErrPending
}

func itShouldComplainsAboutTheEmailNotBeingRegistered() error {
	return godog.ErrPending
}

func theValidPasswordIs(arg1 string) error {
	return godog.ErrPending
}

func itShouldComplainsAboutThePasswordBeingIncorrect() error {
	return godog.ErrPending
}

func iShouldReceiveAnAccessToken() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have the email "([^"]*)"$`, iHaveTheEmail)
	s.Step(`^the password "([^"]*)"$`, thePassword)
	s.Step(`^the email is not already registered$`, theEmailIsNotAlreadyRegistered)
	s.Step(`^I register a new account$`, iRegisterANewAccount)
	s.Step(`^my account should have been created$`, myAccountShouldHaveBeenCreated)
	s.Step(`^the email is already registered$`, theEmailIsAlreadyRegistered)
	s.Step(`^it should complains about the email being already taken$`, itShouldComplainsAboutTheEmailBeingAlreadyTaken)
	s.Step(`^I authenticate$`, iAuthenticate)
	s.Step(`^it should complains about the email not being registered$`, itShouldComplainsAboutTheEmailNotBeingRegistered)
	s.Step(`^the valid password is "([^"]*)"$`, theValidPasswordIs)
	s.Step(`^it should complains about the password being incorrect$`, itShouldComplainsAboutThePasswordBeingIncorrect)
	s.Step(`^I should receive an access token$`, iShouldReceiveAnAccessToken)
}
