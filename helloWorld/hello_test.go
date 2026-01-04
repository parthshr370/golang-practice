package main

import "testing"

func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Mellow("Carol", "English")
		want := "HeyCarol"
		asserCorrectMessage(t, got, want)

	})
	t.Run("say hello world when the string is not supplied", func(t *testing.T) {
		got := Mellow("", "English")
		want := "HeyWorld"
		asserCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {

		got := Mellow("Elodie", "Spanish")
		want := "HolaElodie"
		asserCorrectMessage(t, got, want)
	})

}

// wrapping ourcode to make code more readable
func asserCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
