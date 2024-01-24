// A simple radix like tree specialized for URL path routing
// with supports for wildcards.
package gradix

import (
  "maps"
  "strings"
)

type Radix[P any] struct {
	payload *P
	children map[string]*Radix[P]
	wildcards map[string]*Radix[P]
}

// Create an empty tree
func New[P any]() *Radix[P] {
	return &Radix[P]{}
}

// Add an uri path template and the corresponding payload.
// Path component must be separated by '/'.
// Wildcards path component use a ':' prefix.
// Examples: 
//  radix := New[string]()
//  radix.Add("/users", "List users")
//  radix.Add("/users/:id", "Get a user by id")
//  radix.Add("/users/self", "Get the current user")
//  radix.Add("/users/:user_id/pets/", "...")
//  radix.Add("/users/:user_id/pets/:pet_id", "...")
//  radix.Add("/users/:user_id/friends/:friend_id", "...")
//  radix.Add("/users/:user_id/friends/", "...")
func (self *Radix[P]) Add(url string, payload P) {
  self.add(strings.FieldsFunc(url, func(c rune) bool { return c == '/' }), payload)
}

func (self *Radix[P]) add(path []string, payload P) {
  if len(path) == 0 {
    self.payload = &payload
    return
  }

  root := path[0]
  var target map[string]*Radix[P]

  if root[:1] == ":" {
    root = root[1:]
    if self.wildcards == nil {
      self.wildcards = map[string]*Radix[P]{}
    }
    target = self.wildcards
  } else {
    if self.children == nil {
      self.children = map[string]*Radix[P]{}
    }
    target = self.children 
  }

  child := target[root]
  if child == nil {
    child = &Radix[P]{}
    target[root] = child
  }
  child.add(path[1:], payload)
}

// Match result
type Result[P any] struct {
  // The payload of the added path that matched
  Payload P
  // The values of the wildcard components of the match
  Parameters map[string]string
}

// Search the radix tree with a given path.
// It returns all the match found.
// If there are no match, return an empty slice.
// When a path component match multiple fixeds or wildcards registered paths,
// all paths are returned bu the fixed paths appear first, then the wildcard path.
// Example:
//  radix.Search("/users/toto") // => [{"Get a user by id", [id:toto]}}]
//  radix.Search("/users/self") // => [{"Get the current user", []}, {"Get a user by id", [id:self]}}]
//  radix.Search("/test") // => empty
//  radix.Search("nowhere") // => empty
//  radix.Search("") // => empty
func (self *Radix[P]) Search(url string) []Result[P] {
  results := []Result[P]{}
  parameters := make(map[string]string)
  return self.search(strings.FieldsFunc(url, func(c rune) bool { return c == '/' }), results, parameters)
}

func (self *Radix[P]) search(path []string, results []Result[P], parameters map[string]string) []Result[P] {
  if len(path) == 0 {
    if (self.payload != nil) {
      results = append(results, Result[P]{*self.payload, parameters})
    }
    return results
  }

  if child := self.children[path[0]]; child != nil {
    results = child.search(path[1:], results, parameters)  
  }
  
  if wildcards := self.wildcards; wildcards != nil {
    for name, child := range wildcards {
      sub_parameters := maps.Clone(parameters)
      sub_parameters[name] = path[0]
      results = child.search(path[1:], results, sub_parameters)  
    }
  }

  return results
}
