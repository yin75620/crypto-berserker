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
	return &FtxInit{
		//SubAccount: "",
	}
}

func (fi *FtxInit) IniSetting(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return err
	}

	const section = "FTX"

	fi.SubAccount = cfg.Section(section).Key("SubAccount").String()
	fi.ApiKey = setting.FTX_KEY
	fi.ApiSecretKey = setting.FTX_API_SECRET_KEY

	return nil
}
