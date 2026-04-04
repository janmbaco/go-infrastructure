# Crypto Module

`github.com/janmbaco/go-infrastructure/v2/crypto`

The `crypto` package provides symmetric encryption using AES-GCM. Its API is byte-oriented: it encrypts `[]byte` and returns `[]byte`.

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/crypto
```

## API

```go
type Cipher interface {
    Encrypt(value []byte) ([]byte, error)
    Decrypt(value []byte) ([]byte, error)
}

func NewCipher(key []byte) (Cipher, error)
```

`NewCipher` expects a 32-byte key for AES-256.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/janmbaco/go-infrastructure/v2/crypto"
)

func main() {
    key := []byte("0123456789abcdef0123456789abcdef")

    cipher, err := crypto.NewCipher(key)
    if err != nil {
        panic(err)
    }

    encrypted, err := cipher.Encrypt([]byte("sensitive information"))
    if err != nil {
        panic(err)
    }

    decrypted, err := cipher.Decrypt(encrypted)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(decrypted))
}
```

## Working With Strings

Because the package returns raw bytes, string-based storage should happen at the boundary of your application. Base64 is a common option:

```go
package main

import (
    "encoding/base64"
    "fmt"

    "github.com/janmbaco/go-infrastructure/v2/crypto"
)

func main() {
    key := []byte("0123456789abcdef0123456789abcdef")
    cipher, _ := crypto.NewCipher(key)

    encrypted, _ := cipher.Encrypt([]byte("api-token"))
    encoded := base64.StdEncoding.EncodeToString(encrypted)

    rawCiphertext, _ := base64.StdEncoding.DecodeString(encoded)
    decrypted, _ := cipher.Decrypt(rawCiphertext)

    fmt.Println(encoded)
    fmt.Println(string(decrypted))
}
```

## Encrypting Files

```go
func encryptFile(cipher crypto.Cipher, inputPath, outputPath string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    encrypted, err := cipher.Encrypt(data)
    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, encrypted, 0o600)
}

func decryptFile(cipher crypto.Cipher, inputPath, outputPath string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    decrypted, err := cipher.Decrypt(data)
    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, decrypted, 0o600)
}
```

## DI Integration

`crypto/ioc` registers `crypto.Cipher` and expects a `key` parameter when resolving it:

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/v2/crypto"
    cryptoioc "github.com/janmbaco/go-infrastructure/v2/crypto/ioc"
    "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

func main() {
    key := []byte("0123456789abcdef0123456789abcdef")

    container := dependencyinjection.NewBuilder().
        AddModule(cryptoioc.NewCryptoModule()).
        MustBuild()

    cipher := dependencyinjection.ResolveWithParams[crypto.Cipher](
        container.Resolver(),
        map[string]interface{}{"key": key},
    )

    _, _ = cipher.Encrypt([]byte("secret"))
}
```

## Security Notes

- Use a strong random 32-byte key.
- Treat ciphertext as binary data unless you explicitly encode it.
- Rotate keys through your application boundary, not by editing ciphertext in place.
- Do not use this package for password hashing. Use a password hashing algorithm such as bcrypt, scrypt or Argon2 for passwords.
