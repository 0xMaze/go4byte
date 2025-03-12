package fnprocessor

import (
	"encoding/json"
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
	s := strings.TrimSpace(fs.String())
	var outputsPart string
	if idx := strings.Index(s, "returns"); idx != -1 {
		outputsPart = strings.TrimSpace(s[idx:])
	}

	fnName, inputParamTypes, err := fs.parse()
	if err != nil {
		return "", err
	}
	var inputs []ABIParameter

	for _, p := range inputParamTypes {
		inputs = append(inputs, ABIParameter{
			Name: "",
			Type: string(p),
		})
	}

	var outputs []ABIParameter
	if outputsPart != "" {
		outputsPart = strings.TrimPrefix(outputsPart, "returns")
		outputsPart = strings.TrimSpace(outputsPart)
		if strings.HasPrefix(outputsPart, "(") && strings.HasSuffix(outputsPart, ")") {
			outputsPart = outputsPart[1 : len(outputsPart)-1]
		}
		for {
			var token string
			token, outputsPart, found := strings.Cut(outputsPart, ",")
			token = strings.TrimSpace(token)
			if token != "" {
				parts := strings.Fields(token)
				var outType, outName string
				outType = parts[0]
				if len(parts) > 1 {
					outName = parts[1]
				}
				outputs = append(outputs, ABIParameter{
					Name: outName,
					Type: outType,
				})
			}
			if !found {
				if len(outputsPart) > 0 {
					token = strings.TrimSpace(outputsPart)
					if token != "" {
						parts := strings.Fields(token)
						var outType, outName string
						outType = parts[0]
						if len(parts) > 1 {
							outName = parts[1]
						}
						outputs = append(outputs, ABIParameter{
							Name: outName,
							Type: outType,
						})
					}
				}
				break
			}
		}
	}

	closeParen := strings.Index(fs.String(), ")")
	modifiers := ""
	if closeParen != -1 && closeParen+1 < len(fs.String()) {
		modifiers = strings.TrimSpace(fs.String()[closeParen+1:])
	}
	stateMutability := NonPayable
	if strings.Contains(modifiers, string(View)) {
		stateMutability = View
	} else if strings.Contains(modifiers, string(Pure)) {
		stateMutability = Pure
	} else if strings.Contains(modifiers, string(Payable)) {
		stateMutability = Payable
	}

	abiEntry := ABIEntry{
		Name:            string(fnName),
		Type:            "function",
		Inputs:          inputs,
		Outputs:         outputs,
		StateMutability: stateMutability,
	}
	jsonBytes, err := json.MarshalIndent(abiEntry, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
