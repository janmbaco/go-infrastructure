package configuration

import (
	"fmt"
	"io"
	"os"
)

// MaxSnapshotBytes limita la memoria utilizada al leer configuración operativa.
const MaxSnapshotBytes = 2 * 1024 * 1024

// ReadSnapshot lee y valida una configuración independiente sin crear, modificar ni vigilar archivos.
// Rechaza campos desconocidos y documentos concatenados. Cada llamada devuelve un valor nuevo.
// La validación es obligatoria; si falla, devuelve el valor cero y conserva la causa del error.
func ReadSnapshot[T any](path string, validate func(T) error) (T, error) {
	var empty T
	if validate == nil {
		return empty, fmt.Errorf("snapshot validator is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return empty, fmt.Errorf("open configuration snapshot: %w", err)
	}
	data, err := io.ReadAll(io.LimitReader(file, MaxSnapshotBytes+1))
	closeErr := file.Close()
	if err != nil {
		return empty, fmt.Errorf("read configuration snapshot: %w", err)
	}
	if closeErr != nil {
		return empty, fmt.Errorf("close configuration snapshot: %w", closeErr)
	}
	if len(data) > MaxSnapshotBytes {
		return empty, fmt.Errorf("configuration snapshot exceeds %d bytes", MaxSnapshotBytes)
	}
	value, err := decodeSnapshot[T](data)
	if err != nil {
		return empty, fmt.Errorf("decode configuration snapshot: %w", err)
	}
	if err := validate(value); err != nil {
		return empty, fmt.Errorf("validate configuration snapshot: %w", err)
	}
	return value, nil
}
