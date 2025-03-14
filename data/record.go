// data/record.go
package data

import "fmt"

type Record struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r Record) String() string {
	return fmt.Sprintf("id: %s, name: %s, email: %s", r.ID, r.Name, r.Email)
}
