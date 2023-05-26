package wizlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

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
func NewNameGenerator(repo NameRepository) (*NameGenerator, error) {
	names, err := repo.GetNames()
	if err != nil {
		return nil, err
	}
	return &NameGenerator{
		acceptedNames: names,
	}, nil
}

// NameRepository defines the contract for accessing name data.
type NameRepository interface {
	GetNames() (AcceptedNames, error)
}

// JSONNameRepository is an implementation of the NameRepository using a JSON file.
type JSONNameRepository struct {
	FilePath string
}

// GetNames retrieves the accepted names from a JSON file.
func (r *JSONNameRepository) GetNames() (AcceptedNames, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		return AcceptedNames{}, err
	}
	defer file.Close()

	var names AcceptedNames
	err = json.NewDecoder(file).Decode(&names)
	if err != nil && err != io.EOF {
		return AcceptedNames{}, err
	}

	return names, nil
}

// URLNameRepository is an implementation of the NameRepository using a remote URL.
type URLNameRepository struct {
	URL string
}

// GetNames retrieves the accepted names from a remote URL.
func (r *URLNameRepository) GetNames() (AcceptedNames, error) {
	client := NewAPIClient()
	body, err := client.Get(r.URL)
	if err != nil {
		return AcceptedNames{}, err
	}

	var names AcceptedNames
	err = json.Unmarshal(body, &names)
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
		var sb strings.Builder
		sb.WriteString(name.First)
		if name.Middle != "" {
			sb.WriteString(" ")
			sb.WriteString(name.Middle)
		}
		if name.Last != "" {
			sb.WriteString(name.Last)
		}
		return sb.String(), nil
	}

	return "", errors.New("sorry, the name is not accepted")
}
