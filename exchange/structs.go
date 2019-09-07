package exchange

import (
	"errors"
	"strconv"
)

type JArray map[string]interface{}

type PriceStatus struct {
	Asks [][]float64 `"json:asks"`
	Bids [][]float64 `"json:bids"`
}

func (ps *PriceStatus) GetPair(depth int, pType PriceType) (PricePair, error) {
	switch pType {
	case Ask:
		return ps.getAskPricePair(depth)
	case Bid:
		return ps.getBidPricePair(depth)
	}
	return PricePair{}, errors.New("has no match PriceType")
}

func (ps *PriceStatus) getAskPricePair(depth int) (PricePair, error) {
	return GetPricePair(depth, ps.Asks)
}
func (ps *PriceStatus) getBidPricePair(depth int) (PricePair, error) {
	return GetPricePair(depth, ps.Bids)
}

func GetPricePair(depth int, prices [][]float64) (PricePair, error) {
	var res = PricePair{}
	size := len(prices)
	if depth > size {
		return res, errors.New("depth can't over size")
	}

	index := depth - 1
	res.Price = prices[index][0] // first prize, second volume
	res.Volume = prices[index][1]
	return res, nil
}

func (ps *PriceStatus) SetByJArray(json map[string]interface{}) {
	ps.Asks = transToFloatTwoArray(json["asks"].([]interface{}))
	ps.Bids = transToFloatTwoArray(json["bids"].([]interface{}))
}

func transToFloatTwoArray(askStrArrays []interface{}) [][]float64 {
	res := [][]float64{}
	for _, array := range askStrArrays {
		askFloatArray := []float64{}
		sArray := array.([]interface{})
		for _, s := range sArray {
			res, _ := strconv.ParseFloat(s.(string), 64)
			askFloatArray = append(askFloatArray, res)
		}
		res = append(res, askFloatArray)
	}
	return res
}
