package fnprocessor

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type FnSig string
type ParamType string
type FnName string

func NewFnSig(s string) FnSig {
	return FnSig(strings.TrimSpace(s))
}

func (fs *FnSig) IsEmpty() bool {
	return !(len(fs.String()) > 0)
}

func (fs *FnSig) String() string {
	return string(*fs)
}

func (fs *FnSig) Set(s string) error {
	*fs = NewFnSig(s)
	return nil
}

func (fs *FnSig) Type() string {
	return "string"
}

func (fs *FnSig) parse() (FnName, []ParamType, error) {
	sig := strings.TrimSpace(fs.String())

	sig = strings.TrimPrefix(sig, "function ")

	if idx := strings.Index(sig, "returns"); idx != -1 {
		sig = sig[:idx]
	}

	openParen := strings.Index(sig, "(")
	closeParen := strings.Index(sig, ")")
	if openParen == -1 || closeParen == -1 || closeParen < openParen {
		return FnName(""), nil, fmt.Errorf("invalid function signature format: %s", fs.String())
	}

	fnName := strings.TrimSpace(sig[:openParen])
	paramsStr := sig[openParen+1 : closeParen]

	var paramTypes []ParamType
	for param := range strings.SplitSeq(paramsStr, ",") {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		parts := strings.Fields(param)
		if len(parts) > 0 {
			paramTypes = append(paramTypes, ParamType(parts[0]))
		}
	}

	return FnName(fnName), paramTypes, nil
}

func (fs *FnSig) FourBytes() (string, error) {
	cSig, err := fs.Canonical()
	if err != nil {
		return "", err
	}

	sigBytes := []byte(string(cSig))
	sigHash := crypto.Keccak256(sigBytes)
	selector := sigHash[:4]
	hexSelector := hexutil.Encode(selector)

	return hexSelector, nil
}
