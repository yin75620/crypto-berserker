package ftx

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/yin75620/crypto-berserker/setting"
)

func TestInfo(t *testing.T) {
	fmt.Println("TEST")

	var vftx = NewFtx(http.DefaultClient,
		FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			setting.FTX_SUBACCOUNT})

	vftx.GetAccountInfo()

	//ftx.GetMarkets()
}
