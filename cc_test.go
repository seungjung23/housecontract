package cc_test

import (
	"housecontract/cc"
	"testing"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
)

const (
	alice = `{"Id":"Alice"}`
	bob   = `{"Id":"Bob"}`

	aliceid     = `"Alice"`
	bobid       = `"Bob"`
	emptyOwners = "[]"
	oneOwners   = "[" + alice + "]"
	twoOwners   = "[" + alice + "," + bob + "]"

	timestamp = `"2018-01-01T12:34:56Z"`

	house1  = `{"Id":"1", "Address":"seoul", "OwnerId":"Alice","Price":"3000", "Timestamp":` + timestamp + `}`
	house1b = `{"Id":"1", "Address":"seoul", "OwnerId":"Bob","Price":"3000", "Timestamp":` + timestamp + `}`
	house2  = `{"Id":"2", "Address":"bucheon", "OwnerId":"Alice","Price":"2000", "Timestamp":` + timestamp + `}`

	oneHouses = "[" + house1 + "]"
	twoHouses = "[" + house1 + "," + house2 + "]"

	one = `"1"`
	two = `"2"`
)

func responseOK(res pb.Response) func() bool {
	return func() bool { return res.Status < shim.ERRORTHRESHOLD }
}

func responseFail(res pb.Response) func() bool {
	return func() bool { return res.Status >= shim.ERRORTHRESHOLD }
}

func getBytes(function string, args ...string) [][]byte {
	bytes := make([][]byte, 0, len(args)+1)
	bytes = append(bytes, []byte(function))
	for _, s := range args {
		bytes = append(bytes, []byte(s))
	}
	return bytes
}

// OK1: normal Init()
func TestInit_OK1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) {
		res := stub.MockInit(util.GenerateUUID(), nil)
		assert.Condition(t, responseOK(res))
	}
}

// NG1: unknown method Invoke()
func TestInvoke_NG1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("BadMethod"))
		assert.Condition(t, responseFail(res))
	}
}

// NG1: less arguments
func TestAddOwner_NG1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner"))
		assert.Condition(t, responseFail(res))
	}
}

// NG2: illegal JSON argument
func TestAddOwner_NG2(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", "bad"))
		assert.Condition(t, responseFail(res))
	}
}

// OK1: success
func TestAddOwner_OK(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, emptyOwners, string(res.Payload))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, oneOwners, string(res.Payload))
	}
}

// OK1: 1 Owner
func TestListOwners_OK1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		t.Logf("%s", res.Payload) //끝에 찍히는 로그
		assert.JSONEq(t, oneOwners, string(res.Payload))
	}
}

// OK2: 2 Owners
func TestListOwners_OK2(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, twoOwners, string(res.Payload))
	}
}

// OK1: a single House
func TestAddHouse_OK1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("GetHouse", one))
		assert.Condition(t, responseFail(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house1))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("ListHouses"))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, oneHouses, string(res.Payload))
		}
	}
}

// OK2: two Houses
func TestListHouses_OK2(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house1))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house2))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("ListHouses"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, twoHouses, string(res.Payload))
	}
}

// OK1: change owner from Alice to Bob
func TestUpdateHouse_OK1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house1))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("UpdateHouse", house1b))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("GetHouse", one))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, house1b, string(res.Payload))
		}
	}
}

// NG1: specified house does not exist
func TestUpdateHouse_NG1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("UpdateHouse", house1b))
		assert.Condition(t, responseFail(res))
	}
}

// OK2: transfer from Alice to Bob
func TestTransferHouse_OK1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house1))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("TransferHouse", one, bobid))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("GetHouse", one))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, house1b, string(res.Payload))
		}
	}
}

// NG1: specified House does not exist
func TestTransferHouse_NG1(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("TransferHouse", one, bobid))
		assert.Condition(t, responseFail(res))
	}
}

// NG2: new Owner not found
func TestTransferHouse_NG2(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("AddHouse", house1))

		res = stub.MockInvoke(util.GenerateUUID(), getBytes("TransferHouse", one, bobid))
		assert.Condition(t, responseFail(res))
	}
}

// NG3: less arguments
func TestTransferHouse_NG3(t *testing.T) {
	stub := shim.NewMockStub("housecontract", new(cc.HouseContractCC))
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit(util.GenerateUUID(), nil))) {
		res := stub.MockInvoke(util.GenerateUUID(), getBytes("AddOwner", alice))
		res = stub.MockInvoke(util.GenerateUUID(), getBytes("TransferHouse", one))
		assert.Condition(t, responseFail(res))
	}
}
