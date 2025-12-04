# Crypto Module

`github.com/janmbaco/go-infrastructure/v2/crypto`

Secure encryption utilities for Go applications with AES-256 support.

## Overview

The crypto module provides simple, secure encryption and decryption using industry-standard AES-256-GCM. It's designed for:

- **Configuration encryption** - Secure sensitive config values
- **Data at rest** - Encrypt files and database fields
- **Token generation** - Create encrypted tokens and secrets
- **Password protection** - Secure password storage mechanisms

## Installation

```bash
go get github.com/janmbaco/go-infrastructure/v2/crypto
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-infrastructure/v2/crypto"
)

func main() {
    // Create cipher with 32-byte key (AES-256)
    key := []byte("your-32-byte-secret-key-here!!!!") // Must be exactly 32 bytes
    cipher, err := crypto.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // Encrypt data
    plaintext := "sensitive information"
    encrypted, err := cipher.Encrypt(plaintext)
    if err != nil {
        panic(err)
    }
    fmt.Println("Encrypted:", encrypted)

    // Decrypt data
    decrypted, err := cipher.Decrypt(encrypted)
    if err != nil {
        panic(err)
    }
    fmt.Println("Decrypted:", decrypted)
}
```

## API Reference

### Cipher

The main interface for encryption operations:

```go
type Cipher interface {
    // Encrypt encrypts a string value and returns base64-encoded ciphertext
    Encrypt(value string) (string, error)
    
    // Decrypt decrypts a base64-encoded ciphertext and returns the original string
    Decrypt(encryptedValue string) (string, error)
}
```

### Creating a Cipher

```go
// NewCipher creates a new cipher with the provided key
// Key must be exactly 32 bytes for AES-256
func NewCipher(key []byte) (Cipher, error)
```

## Usage Examples

### Encrypting Configuration Values

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/janmbaco/go-infrastructure/v2/crypto"
)

type Config struct {
    DatabaseURL string `json:"database_url"`
    APIKey      string `json:"api_key"`
}

func main() {
    // Load encryption key from environment
    key := []byte(os.Getenv("ENCRYPTION_KEY")) // Must be 32 bytes
    cipher, err := crypto.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // Encrypt sensitive values
    dbURL := "postgres://user:password@localhost/mydb"
    encryptedDB, err := cipher.Encrypt(dbURL)
    if err != nil {
        panic(err)
    }

    apiKey := "sk-1234567890abcdef"
    encryptedKey, err := cipher.Encrypt(apiKey)
    if err != nil {
        panic(err)
    }

    // Store encrypted config
    config := Config{
        DatabaseURL: encryptedDB,
        APIKey:      encryptedKey,
    }

    data, _ := json.MarshalIndent(config, "", "  ")
    os.WriteFile("config.encrypted.json", data, 0600)

    // Later: decrypt when needed
    decryptedDB, _ := cipher.Decrypt(config.DatabaseURL)
    fmt.Println("Database URL:", decryptedDB)
}
```

### Encrypting Files

```go
func encryptFile(cipher crypto.Cipher, inputPath, outputPath string) error {
    // Read file
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    // Encrypt content
    encrypted, err := cipher.Encrypt(string(data))
    if err != nil {
        return err
    }

    // Write encrypted file
    return os.WriteFile(outputPath, []byte(encrypted), 0600)
}

func decryptFile(cipher crypto.Cipher, inputPath, outputPath string) error {
    // Read encrypted file
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    // Decrypt content
    decrypted, err := cipher.Decrypt(string(data))
    if err != nil {
        return err
    }

    // Write decrypted file
    return os.WriteFile(outputPath, []byte(decrypted), 0600)
}
```

### Integration with DI Container

```go
package main

import (
    "os"
    
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/crypto"
    "github.com/janmbaco/go-infrastructure/v2/crypto/ioc"
)

func main() {
    // Set encryption key in environment
    os.Setenv("ENCRYPTION_KEY", "your-32-byte-secret-key-here!!!!")

    container := di.NewBuilder().
        AddModule(ioc.NewCryptoModule()).
        MustBuild()

    cipher := di.Resolve[crypto.Cipher](container.Resolver())
    
    encrypted, _ := cipher.Encrypt("secret data")
    decrypted, _ := cipher.Decrypt(encrypted)
}
```

## Security Best Practices

### Key Management

1. **Never hardcode keys in source code**
   ```go
   // BAD - Don't do this
   key := []byte("hardcoded-key-in-source-code!!")
   
   // GOOD - Load from environment
   key := []byte(os.Getenv("ENCRYPTION_KEY"))
   ```

2. **Use strong random keys**
   ```go
   // Generate a secure random key
   key := make([]byte, 32)
   if _, err := rand.Read(key); err != nil {
       panic(err)
   }
   ```

3. **Store keys securely**
   - Use environment variables
   - Use secret management systems (HashiCorp Vault, AWS Secrets Manager)
   - Use key management services (AWS KMS, Google Cloud KMS)

4. **Rotate keys periodically**
   - Implement key rotation strategy
   - Re-encrypt data with new keys
   - Maintain key versioning

### Key Size Requirements

```go
// AES-256 requires exactly 32 bytes
key := []byte("must-be-exactly-32-bytes-long!!")

// Validate key length
if len(key) != 32 {
    return errors.New("key must be exactly 32 bytes for AES-256")
}
```

### Secure Key Generation

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func generateKey() ([]byte, error) {
    key := make([]byte, 32) // 32 bytes = 256 bits
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    key, err := generateKey()
    if err != nil {
        panic(err)
    }
    
    // Print as base64 for storage
    encoded := base64.StdEncoding.EncodeToString(key)
    fmt.Println("Generated key (base64):", encoded)
    
    // Store in environment or secrets manager
    // export ENCRYPTION_KEY=<encoded>
}
```

## Error Handling

The module defines specific error types:

```go
// Invalid key length
if len(key) != 32 {
    return nil, fmt.Errorf("invalid key length: expected 32, got %d", len(key))
}

// Encryption failure
encrypted, err := cipher.Encrypt(data)
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
}

// Decryption failure
decrypted, err := cipher.Decrypt(encrypted)
if err != nil {
    return fmt.Errorf("decryption failed: %w", err)
}
```

## Technical Details

### Algorithm

- **Algorithm:** AES-256-GCM (Galois/Counter Mode)
- **Key Size:** 256 bits (32 bytes)
- **Block Size:** 128 bits (16 bytes)
- **Authentication:** Built-in with GCM mode

### Why AES-256-GCM?

- **Security:** Industry standard, FIPS 140-2 approved
- **Performance:** Hardware acceleration on modern CPUs
- **Authenticated:** Provides both confidentiality and integrity
- **AEAD:** Authenticated Encryption with Associated Data

### Output Format

Encrypted values are returned as base64-encoded strings:

```
plaintext → AES-256-GCM → binary ciphertext → base64 encoding → string
```

This makes them safe for:
- JSON storage
- Database text fields
- Configuration files
- HTTP headers

## Testing

Example test:

```go
func TestCipher_EncryptDecrypt(t *testing.T) {
    key := []byte("test-key-must-be-32-bytes-long!")
    cipher, err := crypto.NewCipher(key)
    assert.NoError(t, err)

    plaintext := "sensitive data"
    
    // Encrypt
    encrypted, err := cipher.Encrypt(plaintext)
    assert.NoError(t, err)
    assert.NotEqual(t, plaintext, encrypted)
    
    // Decrypt
    decrypted, err := cipher.Decrypt(encrypted)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}

func TestCipher_InvalidKey(t *testing.T) {
    key := []byte("too-short")
    _, err := crypto.NewCipher(key)
    assert.Error(t, err)
}
```

## Common Use Cases

### Database Field Encryption

```go
type User struct {
    ID       int
    Email    string
    SSN      string // Encrypted field
}

func (u *User) SetSSN(cipher crypto.Cipher, ssn string) error {
    encrypted, err := cipher.Encrypt(ssn)
    if err != nil {
        return err
    }
    u.SSN = encrypted
    return nil
}

func (u *User) GetSSN(cipher crypto.Cipher) (string, error) {
    return cipher.Decrypt(u.SSN)
}
```

### API Token Encryption

```go
type APIToken struct {
    Token     string
    CreatedAt time.Time
}

func generateToken(cipher crypto.Cipher, userID string) (string, error) {
    // Create token data
    data := fmt.Sprintf("%s:%d", userID, time.Now().Unix())
    
    // Encrypt
    return cipher.Encrypt(data)
}

func validateToken(cipher crypto.Cipher, token string) (string, error) {
    // Decrypt
    data, err := cipher.Decrypt(token)
    if err != nil {
        return "", err
    }
    
    // Parse and validate
    parts := strings.Split(data, ":")
    if len(parts) != 2 {
        return "", errors.New("invalid token format")
    }
    
    return parts[0], nil // Return userID
}
```

## Performance Considerations

- **Encryption overhead:** ~10-50μs per operation (depends on data size)
- **Base64 overhead:** ~33% size increase
- **Memory:** Minimal allocation, suitable for high-throughput scenarios

## Troubleshooting

### "invalid key length" error

**Solution:** Ensure your key is exactly 32 bytes:
```go
key := []byte(os.Getenv("ENCRYPTION_KEY"))
if len(key) != 32 {
    panic("Key must be 32 bytes")
}
```

### "cipher: message authentication failed" error

**Cause:** Data was corrupted or tampered with, or wrong key used.

**Solution:** 
- Verify you're using the correct key
- Check data hasn't been modified
- Ensure base64 encoding is intact

## Migration from Other Libraries

### From standard crypto/aes

```go
// Before (manual AES-GCM setup)
block, _ := aes.NewCipher(key)
gcm, _ := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize())
io.ReadFull(rand.Reader, nonce)
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

// After (simplified)
cipher, _ := crypto.NewCipher(key)
encrypted, _ := cipher.Encrypt(string(plaintext))
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md)

## License

Apache License 2.0 - see [LICENSE](../LICENSE)
