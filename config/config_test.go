package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ilayaraja97/clipper/system"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("NewConfig", testNewConfig)
	t.Run("WriteConfig", testWriteConfig)
	t.Run("WriteConfigDefaults", testWriteConfigDefaults)
}

func setupViper(t *testing.T) {
	t.Helper()
	system := system.Analyse()
	tempDir := os.TempDir()

	viper.Reset()
	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(tempDir)
	viper.Set(key, "test_key")
	viper.Set(model, DefaultModel)
	viper.Set(baseURL, "https://openrouter.ai/api/v1")
	viper.Set(proxy, "test_proxy")
	viper.Set(temperature, 0.2)
	viper.Set(maxTokens, 2000)
	viper.Set(user_default_prompt_mode, "exec")
	viper.Set(user_preferences, "test_preferences")

	require.NoError(t, viper.SafeWriteConfigAs(filepath.Join(tempDir, "clipper.json")))
}

func cleanup(t *testing.T) {
	t.Helper()
	viper.Reset()
	require.NoError(t, os.Remove(filepath.Join(os.TempDir(), "clipper.json")))
}

func testNewConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, "test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, DefaultModel, cfg.GetAiConfig().GetModel())
	assert.Equal(t, "https://openrouter.ai/api/v1", cfg.GetAiConfig().GetBaseURL())
	assert.Equal(t, "test_proxy", cfg.GetAiConfig().GetProxy())
	assert.Equal(t, 0.2, cfg.GetAiConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAiConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())
}

func testWriteConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := WriteConfig(ConfigInput{
		Key:         "new_test_key",
		Model:       DefaultModel,
		BaseURL:     "https://openrouter.ai/api/v1",
		Proxy:       "test_proxy",
		Temperature: "0.2",
		MaxTokens:   "2000",
	}, false)
	require.NoError(t, err)

	assert.Equal(t, "new_test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, DefaultModel, cfg.GetAiConfig().GetModel())
	assert.Equal(t, "https://openrouter.ai/api/v1", cfg.GetAiConfig().GetBaseURL())
	assert.Equal(t, "test_proxy", cfg.GetAiConfig().GetProxy())
	assert.Equal(t, 0.2, cfg.GetAiConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAiConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())

	assert.Equal(t, "new_test_key", viper.GetString(key))
	assert.Equal(t, DefaultModel, viper.GetString(model))
	assert.Equal(t, "https://openrouter.ai/api/v1", viper.GetString(baseURL))
	assert.Equal(t, "test_proxy", viper.GetString(proxy))
	assert.Equal(t, 0.2, viper.GetFloat64(temperature))
	assert.Equal(t, 2000, viper.GetInt(maxTokens))
	assert.Equal(t, "exec", viper.GetString(user_default_prompt_mode))
	assert.Equal(t, "test_preferences", viper.GetString(user_preferences))
}

func testWriteConfigDefaults(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := WriteConfig(ConfigInput{}, false)
	require.NoError(t, err)

	assert.Equal(t, DefaultKey, cfg.GetAiConfig().GetKey())
	assert.Equal(t, DefaultModel, cfg.GetAiConfig().GetModel())
	assert.Equal(t, DefaultBaseURL, cfg.GetAiConfig().GetBaseURL())
	assert.Equal(t, DefaultProxy, cfg.GetAiConfig().GetProxy())
	assert.Equal(t, DefaultTemperature, cfg.GetAiConfig().GetTemperature())
	assert.Equal(t, DefaultMaxTokens, cfg.GetAiConfig().GetMaxTokens())
}
