package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCipher_WhenValidKey_ThenReturnsCipher(t *testing.T) {
	// Arrange
	key := []byte("1234567890123456") // 16 bytes for AES-128

	// Act
	c, err := NewCipher(key)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewCipher_WhenInvalidKey_ThenReturnsError(t *testing.T) {
	// Arrange
	key := []byte("short")

	// Act
	c, err := NewCipher(key)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestCipher_Encrypt_WhenValidValue_ThenReturnsEncrypted(t *testing.T) {
	// Arrange
	key := []byte("1234567890123456")
	c, _ := NewCipher(key)
	value := []byte("hello world")

	// Act
	encrypted, err := c.Encrypt(value)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, encrypted)
	assert.NotEqual(t, value, encrypted)
}

func TestCipher_Decrypt_WhenValidEncrypted_ThenReturnsOriginal(t *testing.T) {
	// Arrange
	key := []byte("1234567890123456")
	c, _ := NewCipher(key)
	original := []byte("hello world")
	encrypted, _ := c.Encrypt(original)

	// Act
	decrypted, err := c.Decrypt(encrypted)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, original, decrypted)
}

func TestCipher_Decrypt_WhenInvalidEncrypted_ThenReturnsError(t *testing.T) {
	// Arrange
	key := []byte("1234567890123456")
	c, _ := NewCipher(key)
	invalid := []byte("invalid")

	// Act
	decrypted, err := c.Decrypt(invalid)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decrypted)
}

func TestCipher_RoundTrip_WhenEncryptThenDecrypt_ThenMatchesOriginal(t *testing.T) {
	// Arrange
	key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	c, _ := NewCipher(key)
	original := []byte("test message")

	// Act
	encrypted, _ := c.Encrypt(original)
	decrypted, _ := c.Decrypt(encrypted)

	// Assert
	assert.Equal(t, original, decrypted)
}
