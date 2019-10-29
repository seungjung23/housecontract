package main

import (
	"fmt"
	"housecontract/cc"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	//인터페이스 체킹
	var _ cc.HouseContract = (*cc.HouseContractCC)(nil)

	err := shim.Start(new(cc.HouseContractCC))

	if err != nil {
		fmt.Printf("Error in chaincode process: %s", err)
	}
}
