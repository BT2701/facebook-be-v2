package service

import (
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WinResult struct {
	PaylineID   interface{} `json:"paylineId"`
	Symbol      string      `json:"symbol"`
	Occurrences int         `json:"occurrences"`
	OnWinline   bool        `json:"onWinline"`
	SymbolValue float64     `json:"symbolValue"`
	Positions   [][]int     `json:"positions"`
}

type SpinService interface {
	CalculateWinnings(finalSymbols [][]string, betAmount float64) (float64, []WinResult, error)
	BetOptions() []float64
}

type spinService struct {
	paylines PaylinesService
	symbols  SymbolsService
	configs  ConfigsService
}

func NewSpinService(paylines PaylinesService, symbols SymbolsService, configs ConfigsService) SpinService {
	return &spinService{paylines: paylines, symbols: symbols, configs: configs}
}

func defaultGameName() string {
	if name := os.Getenv("DEFAULT_GAME"); name != "" {
		return name
	}
	return "slot"
}

func (s *spinService) BetOptions() []float64 {
	config, err := s.configs.GetConfig(defaultGameName())
	if err == nil {
		if options := extractBetOptions(config.Data); len(options) > 0 {
			return options
		}
	}
	return []float64{0.5, 1, 1.5, 2, 5, 10}
}

func (s *spinService) CalculateWinnings(finalSymbols [][]string, betAmount float64) (float64, []WinResult, error) {
	paylinesDoc, err := s.paylines.GetPayline(defaultGameName())
	if err != nil {
		return 0, []WinResult{}, err
	}
	symbolsDoc, err := s.symbols.GetSymbol(defaultGameName())
	if err != nil {
		return 0, []WinResult{}, err
	}

	paylines := asSlice(unwrapField(paylinesDoc.Data, "paylines"))
	symbols := asSlice(unwrapField(symbolsDoc.Data, "symbols"))

	var winAmount float64
	results := make([]WinResult, 0)

	for _, raw := range paylines {
		payline, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		pattern := toPattern(payline["pattern"])
		if len(pattern) < 3 {
			continue
		}

		for _, size := range []int{5, 4, 3} {
			if size > len(pattern) {
				continue
			}
			subset := pattern[:size]
			symbol := symbolAt(finalSymbols, subset[0])
			if symbol == "" {
				continue
			}
			matched := true
			for _, pos := range subset {
				if symbolAt(finalSymbols, pos) != symbol {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
			value := symbolValue(symbols, symbol, size)
			winAmount += betAmount * value
			results = append(results, WinResult{
				PaylineID:   payline["id"],
				Symbol:      symbol,
				Occurrences: size,
				OnWinline:   true,
				SymbolValue: value,
				Positions:   subset,
			})
			break
		}
	}

	return winAmount, results, nil
}

func unwrapField(data interface{}, key string) interface{} {
	if data == nil {
		return nil
	}
	if m := asMap(data); m != nil {
		if inner, exists := m[key]; exists {
			return inner
		}
		if inner, exists := m["data"]; exists {
			return unwrapField(inner, key)
		}
	}
	return data
}

func asMap(value interface{}) map[string]interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return typed
	case bson.M:
		return typed
	case primitive.D:
		return typed.Map()
	default:
		return nil
	}
}

func asSlice(value interface{}) []interface{} {
	switch typed := value.(type) {
	case []interface{}:
		return typed
	case primitive.A:
		return []interface{}(typed)
	default:
		return nil
	}
}

func toPattern(value interface{}) [][]int {
	raw, ok := value.([]interface{})
	if !ok {
		return nil
	}
	pattern := make([][]int, 0, len(raw))
	for _, item := range raw {
		pair, ok := item.([]interface{})
		if !ok || len(pair) < 2 {
			continue
		}
		pattern = append(pattern, []int{toInt(pair[0]), toInt(pair[1])})
	}
	return pattern
}

func symbolAt(grid [][]string, pos []int) string {
	if len(pos) < 2 {
		return ""
	}
	col, row := pos[0], pos[1]
	if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[row]) {
		return ""
	}
	return grid[row][col]
}

func symbolValue(symbols []interface{}, name string, count int) float64 {
	for _, raw := range symbols {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if toString(item["name"]) != name {
			continue
		}
		if values, ok := item["values"].(map[string]interface{}); ok {
			if value, exists := values[strconv.Itoa(count)]; exists {
				return toFloat(value)
			}
		}
		if value, exists := item[strconv.Itoa(count)]; exists {
			return toFloat(value)
		}
	}
	return float64(count)
}

func extractBetOptions(data interface{}) []float64 {
	payload := unwrapField(data, "bet_options")
	raw, ok := payload.([]interface{})
	if !ok {
		if m, ok := data.(map[string]interface{}); ok {
			raw, _ = m["bets"].([]interface{})
		}
	}
	options := make([]float64, 0, len(raw))
	for _, item := range raw {
		options = append(options, toFloat(item))
	}
	return options
}

func toInt(value interface{}) int {
	return int(toFloat(value))
}

func toFloat(value interface{}) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case int32:
		return float64(typed)
	case int64:
		return float64(typed)
	case string:
		parsed, _ := strconv.ParseFloat(typed, 64)
		return parsed
	default:
		return 0
	}
}

func toString(value interface{}) string {
	if typed, ok := value.(string); ok {
		return typed
	}
	return ""
}
