package wizlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const defaultNamingListURL = "https://gist.githubusercontent.com/Astridalia/72fa9fb9699b4a9485cd5a17798cd161/raw/62a67360f88cecfc372c678dcf59c1a511cf0159/w101_names.json"

type Name struct {
	First  string
	Middle string
	Last   string
}

type AcceptedNames struct {
	Names []string `json:"names"`
}

// NameGenerator provides methods for generating valid names based on the input and the accepted names list.
type NameGenerator struct {
	acceptedNames AcceptedNames
}

// NewNameGenerator creates a new instance of NameGenerator and retrieves the default accepted names from the provided URL.
func NewNameGenerator() (*NameGenerator, error) {
	names, err := getDefaultNames()
	if err != nil {
		return nil, err
	}
	return &NameGenerator{
		acceptedNames: names,
	}, nil
}

// getDefaultNames retrieves the default accepted names from the provided URL.
func getDefaultNames() (AcceptedNames, error) {
	resp, err := http.Get(defaultNamingListURL)
	if err != nil {
		return AcceptedNames{}, err
	}
	defer resp.Body.Close()
	var names AcceptedNames
	err = json.NewDecoder(resp.Body).Decode(&names)
	if err != nil {
		return AcceptedNames{}, err
	}
	return names, nil
}

// GenerateName generates a valid name based on the input and the accepted names list.
func (g *NameGenerator) GenerateName(input string) (string, error) {
	pattern := fmt.Sprintf(`(?i)^(%s)( (%s))?((%s))?$`, strings.Join(g.acceptedNames.Names, "|"), strings.Join(g.acceptedNames.Names, "|"), strings.Join(g.acceptedNames.Names, "|"))
	nameRegex := regexp.MustCompile(pattern)

	nameParts := strings.Split(input, " ")
	var name Name
	switch len(nameParts) {
	case 1:
		name.First = nameParts[0]
	case 2:
		name.First = nameParts[0]
		name.Last = nameParts[1]
	default:
		name.First = nameParts[0]
		name.Middle = strings.Join(nameParts[1:len(nameParts)-1], " ")
		name.Last = nameParts[len(nameParts)-1]
	}

	if nameRegex.MatchString(input) {
		return fmt.Sprintf("%s %s%s", name.First, name.Middle, name.Last), nil
	}

	return "", errors.New("sorry, the name is not accepted")
}
