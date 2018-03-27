package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessageToTelegramChannel(t *testing.T) {
	config, err := GetTelegramConfig("./.env")
	assert.NoError(t, err)
	err = SendMessageToTelegramChannel(config)
	assert.NoError(t, err)
}
