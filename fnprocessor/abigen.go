package fnprocessor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ABIParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ABIEntry struct {
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	Inputs          []ABIParameter  `json:"inputs"`
	Outputs         []ABIParameter  `json:"outputs"`
	StateMutability StateMutability `json:"stateMutability"`
}

type StateMutability string

const (
	NonPayable StateMutability = "nonpayable"
	View       StateMutability = "view"
	Pure       StateMutability = "pure"
	Payable    StateMutability = "payable"
)

func (fs *FnSig) GenerateABI() (string, error) {
	fnName, inputParams, err := fs.parse()
	if err != nil {
		return "", err
	}

	abiEntry := ABIEntry{
		Name:            string(fnName),
		Type:            "function",
		Inputs:          buildInputParameters(inputParams),
		Outputs:         parseOutputs(fs.getOutputsPart()),
		StateMutability: fs.determineStateMutability(),
	}

	return marshalABI(abiEntry)
}

func buildInputParameters(params []Param) []ABIParameter {
	inputs := make([]ABIParameter, 0, len(params))
	for _, p := range params {
		inputs = append(inputs, ABIParameter{
			Name: p.Name,
			Type: string(p.Type),
		})
	}
	return inputs
}

func (fs *FnSig) getOutputsPart() string {
	s := strings.TrimSpace(fs.String())
	if idx := strings.Index(s, "returns"); idx != -1 {
		return strings.TrimSpace(s[idx:])
	}
	return ""
}

func parseOutputs(outputsPart string) []ABIParameter {
	if outputsPart == "" {
		return nil
	}

	outputsPart = strings.TrimPrefix(outputsPart, "returns")
	outputsPart = strings.TrimSpace(outputsPart)

	if strings.HasPrefix(outputsPart, "(") && strings.HasSuffix(outputsPart, ")") {
		outputsPart = outputsPart[1 : len(outputsPart)-1]
	}

	return parseParameterList(outputsPart)
}

func parseParameterList(paramStr string) []ABIParameter {
	var params []ABIParameter
	var current strings.Builder
	depth := 0

	for _, r := range paramStr {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				params = append(params, parseParamToken(current.String()))
				current.Reset()
				continue
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		params = append(params, parseParamToken(current.String()))
	}

	return params
}

func parseParamToken(token string) ABIParameter {
	token = strings.TrimSpace(token)
	parts := strings.Fields(token)
	if len(parts) == 0 {
		return ABIParameter{}
	}

	outType := parts[0]
	var outName string
	if len(parts) > 1 {
		outName = parts[1]
	}

	return ABIParameter{
		Type: outType,
		Name: outName,
	}
}

func (fs *FnSig) determineStateMutability() StateMutability {
	modifiers := fs.getModifiers()
	for _, m := range strings.Fields(modifiers) {
		switch strings.ToLower(m) {
		case string(View):
			return View
		case string(Pure):
			return Pure
		case string(Payable):
			return Payable
		}
	}
	return NonPayable
}

func (fs *FnSig) getModifiers() string {
	s := fs.String()
	closeParen := strings.Index(s, ")")
	if closeParen == -1 || closeParen+1 >= len(s) {
		return ""
	}
	return strings.TrimSpace(s[closeParen+1:])
}

func marshalABI(entry ABIEntry) (string, error) {
	jsonBytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal ABI: %w", err)
	}
	return string(jsonBytes), nil
}
