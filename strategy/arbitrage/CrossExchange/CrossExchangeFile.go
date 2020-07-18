package CrossExchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func savePairMapToFile(positionCrossPairMap map[string][]CrossPair) {

	jmap := map[string][]CrossPairJson{}

	for k, pairArray := range positionCrossPairMap {
		for _, value := range pairArray {
			jmap[k] = append(jmap[k], value.toJson())
		}
	}

	b, err := json.Marshal(jmap)
	if err != nil {
		fmt.Println("savePairMapToFile json.Marshal fail", err)
		return
	}
	err = ioutil.WriteFile(PairMapFileName, b, 0644)
	if err != nil {
		fmt.Println("savePairMapToFile ioutil.WriteFile fail", err)
		return
	}
}

func loadPairMapFromFile(exchanges []exc.Exchange) (map[string][]CrossPair, error) {

	res := map[string][]CrossPair{}

	b, err := ioutil.ReadFile(PairMapFileName)
	if err != nil {
		//fmt.Println("loadPairMapFromFile ioutil.ReadFile fail", err)
		return res, nil
	}

	jmap := map[string][]CrossPairJson{}
	err = json.Unmarshal(b, &jmap)
	if err != nil {
		fmt.Println("loadPairMapFromFile json.Unmarshal", err)
		return res, err
	}

	res, err = jMapToStruct(jmap, exchanges)
	if err != nil {
		fmt.Println("loadPairMapFromFile jMapToStruct fail ", err)
		return res, err
	}

	return res, nil
}

func jMapToStruct(jmap map[string][]CrossPairJson, exchanges []exc.Exchange) (map[string][]CrossPair, error) {

	res := map[string][]CrossPair{}

	for k, pairArray := range jmap {
		for _, value := range pairArray {
			askExchange, err := findExchange(exchanges, value.AskExchangeName)
			if err != nil {
				return res, fmt.Errorf("ejMapToStruct findExchange ask. %g", err)
			}
			bidExchange, err := findExchange(exchanges, value.BidExchangeName)
			if err != nil {
				return res, fmt.Errorf("ejMapToStruct findExchange bid. %g", err)
			}

			cp := NewCrossPair(askExchange, bidExchange, value.AskPair, value.BidPair, value.Symbol)
			cp.orderVolume = value.OrderVolume

			res[k] = append(res[k], *cp)
		}
	}

	return res, nil
}

func findExchange(exchanges []exc.Exchange, excName string) (exc.Exchange, error) {
	for _, exchange := range exchanges {
		if excName == exchange.GetName() {
			return exchange, nil
		}
	}
	return nil, fmt.Errorf("Has no match exchange:%s, size:%d", excName, len(exchanges))
}
