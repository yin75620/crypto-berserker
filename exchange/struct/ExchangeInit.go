package exchange

import (
	"fmt"

	"github.com/go-ini/ini"
)

type ExchangeInit struct {
	ApiKey       string
	ApiSecretKey string
	//SubAccount   string
	SectionName string
}

func NewExchangeInit(sectionName string) *ExchangeInit {
	ei := ExchangeInit{}
	ei.SectionName = sectionName
	return &ei
}

func (ei *ExchangeInit) SetKey(apiKey, apiSecretKey string) {
	ei.ApiKey = apiKey
	ei.ApiSecretKey = apiSecretKey
}

func (ei *ExchangeInit) IniSetting(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return err
	}

	ei.IniSettingByFile(cfg)

	return nil
}

func (ei *ExchangeInit) IniSettingByFile(cfg *ini.File) {
	section := ei.SectionName

	//ei.SubAccount = cfg.Section(section).Key("SubAccount").String()

	apiKey := cfg.Section(section).Key("ApiKey").String()
	if apiKey != "" {
		ei.ApiKey = apiKey
	}

	apiSecretKey := cfg.Section(section).Key("ApiSecretKey").String()
	if apiSecretKey != "" {
		ei.ApiSecretKey = apiSecretKey
	}
}
