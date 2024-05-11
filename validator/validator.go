package validator

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetRegisterValidator(ec *ethclient.Client, ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{} = make(map[string]interface{})
	err := ec.Client().CallContext(ctx, &result, "bsc_getRegisterValidators")
	if nil != err {
		return nil, err
	}

	return result, nil
}
