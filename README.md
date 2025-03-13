# fbyte - Solidity Function Signature Processo
A CLI tool to work with Ethereum/Solidity function signatures. Generate canonical ABIs and calculate 4-byte function selectors from function signatures

## Features
- 🛠 Generate canonical ABI JSON from function signatures
- 🔍 Calculate 4-byte function selectors (EIP-712)
- 📦 Lightweight single-binary tool
- 🚰 Supports standard Solidity function signature syntax

## Installation

```bash
go install github.com/0xMaze/go4byte@latest
```

## Usage

### Basic Syntax
```bash
fbyte --sig "<function signature>" [flags]
```

### Flags
| Flag       | Shorthand | Description                                 |
|------------|-----------|---------------------------------------------|
| `--abi`    | `-a`      | Generate function ABI JSON                  |
| `--four`   | `-f`      | Calculate 4-byte function selector          |
| `--sig`    | `-s`      | Function signature (required)               |
| `--exp`    | `-e`      | Export generated ABI                        |
| `--out`    | `-o`      | Export file path                            |

## Examples

1. **Get both ABI and 4-byte selector**
```bash
fbyte -s "function transfer(address to, uint256 amount)" -a -f
```

2. **Generate ABI only**
```bash
fbyte -s "function balanceOf(address) external view returns (uint256)" -a
```

3. **Calculate 4-byte selector only**
```bash
fbyte -s "function approve(address spender, uint256 value)" -f
```
### Output Examples

### ABI Generation Output
```json
{
  "name": "transfer",
  "type": "function",
  "inputs": [
    {
      "name": "to",
      "type": "address"
    },
    {
      "name": "amount",
      "type": "uint256"
    }
  ],
  "outputs": null,
  "stateMutability": "nonpayable"
}
```

### 4-byte Selector Output
```
0xa9059cbb
```

## How It Works

The tool processes function signatures in standard Solidity format:
1. Parses function name and parameters
2. Normalizes to canonical form (function name + parameter types)
3. Calculates Keccak256 hash of canonical signature
4. Takes first 4 bytes of hash as selector

⚠️ Note: Currently supports basic types (address, uint256, etc). Complex types (tuples, arrays) not yet supported.
