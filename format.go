package money

import (
	"strings"
	"sync/atomic"
	"unicode/utf8"
)

// SymbolPosition controls where the symbol appears relative to the amount.
// Example: SymbolSuffix yields "10.50 USD" when SymbolKind uses the code.
type SymbolPosition int32

// SymbolKind chooses the symbol source.
// Example: SymbolUseCurrencyCode yields "USD" for Currency{Code:"USD"}.
type SymbolKind int32

const (
	// SymbolPrefix places the symbol before the amount.
	SymbolPrefix SymbolPosition = iota
	// SymbolSuffix places the symbol after the amount.
	SymbolSuffix
)

const (
	// SymbolUseCurrencySymbol uses Currency.Symbol.
	SymbolUseCurrencySymbol SymbolKind = iota
	// SymbolUseCurrencyCode uses Currency.Code.
	SymbolUseCurrencyCode
	// SymbolUseCustom uses FormatConfig.CustomSymbol.
	SymbolUseCustom
)

// FormatConfig defines formatting behavior for Money rendering.
// Example: DecimalSeparator="," and ThousandsSeparator="." yields "1.234,56".
type FormatConfig struct {
	DecimalSeparator   string
	ThousandsSeparator string
	SymbolPosition     SymbolPosition
	SymbolKind         SymbolKind
	CustomSymbol       string
	Space              bool
}

var formatConfig atomic.Value

func init() {
	formatConfig.Store(FormatConfig{
		DecimalSeparator:   ".",
		ThousandsSeparator: "",
		SymbolPosition:     SymbolPrefix,
		SymbolKind:         SymbolUseCurrencySymbol,
		CustomSymbol:       "",
		Space:              false,
	})
}

// SetFormat sets the global default formatting configuration.
// Example: SetFormat(FormatConfig{DecimalSeparator:",", SymbolPosition:SymbolSuffix}).
func SetFormat(cfg FormatConfig) error {
	if err := validateFormat(cfg); err != nil {
		return err
	}
	formatConfig.Store(cfg)
	return nil
}

// DefaultFormat returns the current global format configuration.
// Example: DefaultFormat().DecimalSeparator -> ".".
func DefaultFormat() FormatConfig {
	return formatConfig.Load().(FormatConfig)
}

// Format renders Money using a local (per-call) configuration.
// Example: m.Format(FormatConfig{SymbolKind:SymbolUseCurrencyCode}) -> "10.50 USD".
func (m Money) Format(cfg FormatConfig) (string, error) {
	if err := validateFormat(cfg); err != nil {
		return "", err
	}
	return formatWithConfig(m, cfg)
}

func formatWithConfig(m Money, cfg FormatConfig) (string, error) {
	absDigits := absInt64String(m.amount)
	intPart, fracPart := splitAmount(absDigits, m.currency.Scale)
	if cfg.ThousandsSeparator != "" {
		intPart = groupThousands(intPart, cfg.ThousandsSeparator)
	}
	amount := intPart
	if fracPart != "" {
		amount = amount + cfg.DecimalSeparator + fracPart
	}

	symbol, err := formatSymbol(m.currency, cfg)
	if err != nil {
		return "", err
	}

	sep := ""
	if cfg.Space {
		sep = " "
	}
	if symbol == "" {
		sep = ""
	}

	if cfg.SymbolPosition == SymbolSuffix {
		return signPrefix(m.amount) + amount + sep + symbol, nil
	}
	return signPrefix(m.amount) + symbol + sep + amount, nil
}

func formatSymbol(currency Currency, cfg FormatConfig) (string, error) {
	switch cfg.SymbolKind {
	case SymbolUseCurrencySymbol:
		return currency.Symbol, nil
	case SymbolUseCurrencyCode:
		return currency.Code, nil
	case SymbolUseCustom:
		if cfg.CustomSymbol == "" {
			return "", ErrInvalidOperation
		}
		return cfg.CustomSymbol, nil
	default:
		return "", ErrInvalidOperation
	}
}

func validateFormat(cfg FormatConfig) error {
	if cfg.DecimalSeparator == "" {
		return ErrInvalidOperation
	}
	if utf8.RuneCountInString(cfg.DecimalSeparator) != 1 {
		return ErrInvalidOperation
	}
	if cfg.ThousandsSeparator != "" && utf8.RuneCountInString(cfg.ThousandsSeparator) != 1 {
		return ErrInvalidOperation
	}
	if cfg.ThousandsSeparator != "" && cfg.ThousandsSeparator == cfg.DecimalSeparator {
		return ErrInvalidOperation
	}
	if cfg.SymbolKind == SymbolUseCustom && cfg.CustomSymbol == "" {
		return ErrInvalidOperation
	}
	switch cfg.SymbolPosition {
	case SymbolPrefix, SymbolSuffix:
	default:
		return ErrInvalidOperation
	}
	switch cfg.SymbolKind {
	case SymbolUseCurrencySymbol, SymbolUseCurrencyCode, SymbolUseCustom:
	default:
		return ErrInvalidOperation
	}
	return nil
}

func splitAmount(absDigits string, scale int32) (string, string) {
	if scale <= 0 {
		return absDigits, ""
	}
	scaleInt := int(scale)
	if len(absDigits) <= scaleInt {
		absDigits = strings.Repeat("0", scaleInt-len(absDigits)+1) + absDigits
	}
	intPart := absDigits[:len(absDigits)-scaleInt]
	fracPart := absDigits[len(absDigits)-scaleInt:]
	return intPart, fracPart
}

func groupThousands(intPart, sep string) string {
	if len(intPart) <= 3 {
		return intPart
	}
	groups := (len(intPart) - 1) / 3
	out := make([]byte, 0, len(intPart)+groups*len(sep))
	start := len(intPart) % 3
	if start == 0 {
		start = 3
	}
	out = append(out, intPart[:start]...)
	for i := start; i < len(intPart); i += 3 {
		out = append(out, sep...)
		out = append(out, intPart[i:i+3]...)
	}
	return string(out)
}
