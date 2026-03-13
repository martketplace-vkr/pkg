package crypto

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

const (
	merchantIDHeaderName = "Merchant"
	businessIDHeaderName = "Business"
)

func ParseMerchantID(request *fasthttp.Request) (int64, error) {
	merchantID, err := strconv.ParseInt(string(request.Header.Peek(merchantIDHeaderName)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse merchant id")
	}

	if merchantID == 0 {
		return 0, fmt.Errorf("merchant id is 0")
	}

	return merchantID, nil
}

func ParseBusinessID(request *fasthttp.Request) (int64, error) {
	businessID, err := strconv.ParseInt(string(request.Header.Peek(businessIDHeaderName)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse business id")
	}

	if businessID == 0 {
		return 0, fmt.Errorf("business id is 0")
	}

	return businessID, nil
}
