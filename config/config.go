package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ilayaraja97/clipper/logger"
	"github.com/ilayaraja97/clipper/system"
	"github.com/spf13/viper"
)

type Config struct {
	ai     AiConfig
	user   UserConfig
	system *system.Analysis
}

type ConfigInput struct {
	Key         string
	Model       string
	BaseURL     string
	Proxy       string
	Temperature string
	MaxTokens   string
}

func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

func NewConfig() (*Config, error) {
	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(filepath.Join(system.GetHomeDirectory(), ".config"))

	logger.Log.Debug().Str("configPath", filepath.Join(system.GetHomeDirectory(), ".config", strings.ToLower(system.GetApplicationName()))).Msg("loading config")

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Debug().Err(err).Msg("failed to read config")
		return nil, err
	}

	logger.Log.Info().
		Str("model", viper.GetString(model)).
		Str("baseURL", viper.GetString(baseURL)).
		Msg("config loaded")

	return &Config{
		ai: AiConfig{
			key:         viper.GetString(key),
			model:       viper.GetString(model),
			baseURL:     viper.GetString(baseURL),
			proxy:       viper.GetString(proxy),
			temperature: viper.GetFloat64(temperature),
			maxTokens:   viper.GetInt(maxTokens),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

func WriteConfig(input ConfigInput, write bool) (*Config, error) {
	system := system.Analyse()

	keyVal := strings.TrimSpace(input.Key)
	if keyVal == "" {
		keyVal = DefaultKey
	}

	modelVal := strings.TrimSpace(input.Model)
	if modelVal == "" {
		modelVal = DefaultModel
	}

	baseURLVal := strings.TrimSpace(input.BaseURL)
	if baseURLVal == "" {
		baseURLVal = DefaultBaseURL
	}

	proxyVal := strings.TrimSpace(input.Proxy)

	temperatureVal := DefaultTemperature
	if value := strings.TrimSpace(input.Temperature); value != "" {
		parsedTemperature, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid temperature %q: %w", value, err)
		}
		temperatureVal = parsedTemperature
	}

	maxTokensVal := DefaultMaxTokens
	if value := strings.TrimSpace(input.MaxTokens); value != "" {
		parsedMaxTokens, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid max tokens %q: %w", value, err)
		}
		maxTokensVal = parsedMaxTokens
	}

	// ai defaults
	viper.Set(key, keyVal)
	viper.Set(model, modelVal)
	viper.Set(baseURL, baseURLVal)
	viper.Set(proxy, proxyVal)
	viper.Set(temperature, temperatureVal)
	viper.Set(maxTokens, maxTokensVal)

	// user defaults
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if write {
		cfgPath := system.GetConfigFile()
		logger.Log.Info().Str("configPath", cfgPath).Msg("writing config")
		if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
			logger.Log.Error().Err(err).Msg("failed to create config directory")
			return nil, err
		}
		if err := viper.SafeWriteConfigAs(cfgPath); err != nil {
			logger.Log.Error().Err(err).Msg("failed to write config file")
			return nil, err
		}
		logger.Log.Info().Msg("config written successfully")
	}

	return NewConfig()
}
