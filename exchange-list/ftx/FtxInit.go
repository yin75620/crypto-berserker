package ftx

import (
	"fmt"

	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/setting"
)

type FtxInit struct {
	ApiKey       string
	ApiSecretKey string
	SubAccount   string
}

func NewFtxInit() *FtxInit {
	return &FtxInit{}
}

func (fi *FtxInit) IniSetting(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return err
	}

	fi.IniSettingByFile(cfg)

	return nil
}

func (fi *FtxInit) IniSettingByFile(cfg *ini.File) {
	const section = "FTX"

	fi.SubAccount = cfg.Section(section).Key("SubAccount").String()
	fi.ApiKey = setting.FTX_KEY
	fi.ApiSecretKey = setting.FTX_API_SECRET_KEY
}
