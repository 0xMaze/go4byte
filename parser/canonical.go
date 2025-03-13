package fnprocessor

import (
	"fmt"
	"strings"
)

type CanonicalFnSig string

func (fs *FnSig) Canonical() (CanonicalFnSig, error) {
	fnName, paramTypes, err := fs.parse()
	if err != nil {
		return "", err
	}
	var types []string
	for _, p := range paramTypes {
		types = append(types, string(p.Type))
	}
	canonical := fmt.Sprintf("%s(%s)", fnName, strings.Join(types, ","))
	return CanonicalFnSig(canonical), nil
}
