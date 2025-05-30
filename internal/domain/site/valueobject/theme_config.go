package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// ThemeConfig 主题配置值对象
type ThemeConfig struct {
	primaryColor    string
	secondaryColor  string
	accentColor     string
	textColor       string
	backgroundColor string
	fontFamily      string
	headerStyle     string
	borderRadius    string
	customCSS       string
}

// NewThemeConfig 创建主题配置
func NewThemeConfig(
	primaryColor, secondaryColor, accentColor, textColor, backgroundColor,
	fontFamily, headerStyle, borderRadius, customCSS string,
) (ThemeConfig, error) {
	config := ThemeConfig{
		primaryColor:    primaryColor,
		secondaryColor:  secondaryColor,
		accentColor:     accentColor,
		textColor:       textColor,
		backgroundColor: backgroundColor,
		fontFamily:      fontFamily,
		headerStyle:     headerStyle,
		borderRadius:    borderRadius,
		customCSS:       customCSS,
	}
	
	if err := config.validate(); err != nil {
		return ThemeConfig{}, err
	}
	
	return config, nil
}

// NewDefaultThemeConfig 创建默认主题配置
func NewDefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		primaryColor:    "#007bff",
		secondaryColor:  "#6c757d",
		accentColor:     "#28a745",
		textColor:       "#333333",
		backgroundColor: "#ffffff",
		fontFamily:      "Arial, sans-serif",
		headerStyle:     "standard",
		borderRadius:    "medium",
		customCSS:       "",
	}
}

// Getters
func (tc ThemeConfig) PrimaryColor() string {
	return tc.primaryColor
}

func (tc ThemeConfig) SecondaryColor() string {
	return tc.secondaryColor
}

func (tc ThemeConfig) AccentColor() string {
	return tc.accentColor
}

func (tc ThemeConfig) TextColor() string {
	return tc.textColor
}

func (tc ThemeConfig) BackgroundColor() string {
	return tc.backgroundColor
}

func (tc ThemeConfig) FontFamily() string {
	return tc.fontFamily
}

func (tc ThemeConfig) HeaderStyle() string {
	return tc.headerStyle
}

func (tc ThemeConfig) BorderRadius() string {
	return tc.borderRadius
}

func (tc ThemeConfig) CustomCSS() string {
	return tc.customCSS
}

// UpdatePrimaryColor 更新主色
func (tc ThemeConfig) UpdatePrimaryColor(color string) (ThemeConfig, error) {
	if err := validateColor(color); err != nil {
		return tc, err
	}
	tc.primaryColor = color
	return tc, nil
}

// UpdateSecondaryColor 更新次色
func (tc ThemeConfig) UpdateSecondaryColor(color string) (ThemeConfig, error) {
	if err := validateColor(color); err != nil {
		return tc, err
	}
	tc.secondaryColor = color
	return tc, nil
}

// UpdateCustomCSS 更新自定义CSS
func (tc ThemeConfig) UpdateCustomCSS(css string) (ThemeConfig, error) {
	if len(css) > 10000 {
		return tc, errors.New("自定义CSS长度不能超过10000个字符")
	}
	tc.customCSS = css
	return tc, nil
}

// Equals 判断两个主题配置是否相等
func (tc ThemeConfig) Equals(other ThemeConfig) bool {
	return tc.primaryColor == other.primaryColor &&
		tc.secondaryColor == other.secondaryColor &&
		tc.accentColor == other.accentColor &&
		tc.textColor == other.textColor &&
		tc.backgroundColor == other.backgroundColor &&
		tc.fontFamily == other.fontFamily &&
		tc.headerStyle == other.headerStyle &&
		tc.borderRadius == other.borderRadius &&
		tc.customCSS == other.customCSS
}

// validate 验证主题配置
func (tc ThemeConfig) validate() error {
	if err := validateColor(tc.primaryColor); err != nil {
		return errors.New("主色无效: " + err.Error())
	}
	if err := validateColor(tc.secondaryColor); err != nil {
		return errors.New("次色无效: " + err.Error())
	}
	if err := validateColor(tc.accentColor); err != nil {
		return errors.New("强调色无效: " + err.Error())
	}
	if err := validateColor(tc.textColor); err != nil {
		return errors.New("文本色无效: " + err.Error())
	}
	if err := validateColor(tc.backgroundColor); err != nil {
		return errors.New("背景色无效: " + err.Error())
	}
	if err := validateHeaderStyle(tc.headerStyle); err != nil {
		return err
	}
	if err := validateBorderRadius(tc.borderRadius); err != nil {
		return err
	}
	if len(tc.customCSS) > 10000 {
		return errors.New("自定义CSS长度不能超过10000个字符")
	}
	return nil
}

// validateColor 验证颜色值（支持hex和rgb格式）
func validateColor(color string) error {
	if color == "" {
		return errors.New("颜色不能为空")
	}
	
	// 验证hex格式 (#ffffff)
	hexRegex := regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	if hexRegex.MatchString(color) {
		return nil
	}
	
	// 验证rgb格式 (rgb(255,255,255))
	rgbRegex := regexp.MustCompile(`^rgb\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)$`)
	if rgbRegex.MatchString(color) {
		return nil
	}
	
	// 验证rgba格式
	rgbaRegex := regexp.MustCompile(`^rgba\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*[01]?(\.\d+)?\s*\)$`)
	if rgbaRegex.MatchString(color) {
		return nil
	}
	
	return errors.New("颜色格式无效，请使用hex或rgb格式")
}

// validateHeaderStyle 验证头部样式
func validateHeaderStyle(style string) error {
	validStyles := []string{"standard", "centered", "minimal"}
	for _, validStyle := range validStyles {
		if style == validStyle {
			return nil
		}
	}
	return errors.New("无效的头部样式，支持: " + strings.Join(validStyles, ", "))
}

// validateBorderRadius 验证边框圆角
func validateBorderRadius(radius string) error {
	validRadius := []string{"none", "small", "medium", "large"}
	for _, validR := range validRadius {
		if radius == validR {
			return nil
		}
	}
	return errors.New("无效的边框圆角，支持: " + strings.Join(validRadius, ", "))
} 