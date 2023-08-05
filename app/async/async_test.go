package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tkudlicka/portflux-api/config"
)

// TestNew_Ok checks that New creates a new async struct with the expected values
func TestNew_Ok(t *testing.T) {
	// Arrange
	expectedConfig := config.Config{}

	// Act
	async := New(expectedConfig)

	// Assert
	assert.Equal(t, expectedConfig, async.config)
}

// TestRun_ContextCancelled checks that Run finishes and returns the expected error when the context gets cancelled
func TestRun_ContextCancelled(t *testing.T) {
	// Arrange
	async := &async{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	expectedError := "async process stopped"

	// Act
	errFunc := async.Run(ctx, cancel)

	// Assert
	assert.Equal(t, expectedError, errFunc().Error())
}
