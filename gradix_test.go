package gradix

import (
	"testing"
	"reflect"
)

func TestGradix(t *testing.T) {
	radix := New[string]()
  radix.Add("/users", "List users")
  radix.Add("/users/:id", "Get a user by id")
  radix.Add("/users/self", "Get the current user")
  radix.Add("/users/:user_id/pets/", "...")
  radix.Add("/users/:user_id/pets/:pet_id", "...")
  radix.Add("/users/:user_id/friends/:friend_id", "...")
  radix.Add("/users/:user_id/friends/", "...")
	radix.Add("/", "root")

	// Fuck type inference i guess
	tests := map[string][]Result[string]{
		// Simple test
		"/users/toto":{Result[string]{"Get a user by id", map[string]string{"id": "toto"}}},

		// Trailing an duplicate slashes
		"/users//toto/":{Result[string]{"Get a user by id", map[string]string{"id": "toto"}}},

		// Wildcard and fixed path collision: fixed first
		"/users/self":{
			Result[string]{"Get the current user", map[string]string{}},
			Result[string]{"Get a user by id", map[string]string{"id": "self"}},
		},

		// No match
		"nowhere":{},
		
		// Various root equivalent
		"":{Result[string]{"root", map[string]string{}}},
		"/":{Result[string]{"root", map[string]string{}}},
		"///":{Result[string]{"root", map[string]string{}}},
	}

	for test, expect := range tests {
		got := radix.Search(test)
		if !reflect.DeepEqual(got, expect) {
			t.Fatalf("Search '%v' got '%v' expected '%v'\n", test, got, expect)
		}
	}
}
