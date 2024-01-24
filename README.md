# Gradix

A simple radix like tree specialized for URL path routing with supports for wildcards.  

## Usage

`GOPROXY=direct go get github.com/globoplox/gradix`

```go
package main

import (
  "github.com/globoplox/gradix"
  "fmt"
)

func main() {
  radix := gradix.New[string]()
  radix.Add("/users", "List users")
  radix.Add("/users/:id", "Get a user by id")
  radix.Add("/users/self", "Get the current user")
  radix.Add("/users/:user_id/pets/", "...")
  radix.Add("/users/:user_id/pets/:pet_id", "...")
  radix.Add("/users/:user_id/friends/:friend_id", "...")
  radix.Add("/users/:user_id/friends/", "...")
  fmt.Printf("%v\n", radix.Search("/users/toto")) // => [{Get a user by id map[id:toto]}]
  fmt.Printf("%v\n", radix.Search("/users/toto/")) // => [{Get a user by id map[id:toto]}]
  fmt.Printf("%v\n", radix.Search("/users/self")) // => [{Get the current user map[]} {Get a user by id map[id:self]}]
  fmt.Printf("%v\n", radix.Search("/test")) // => []
  fmt.Printf("%v\n", radix.Search("nowhere")) // => []
  fmt.Printf("%v\n", radix.Search("")) // => []
  fmt.Printf("%v\n", radix.Search("/////")) // => []
}
```
