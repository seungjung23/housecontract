package cc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const prefixOwner = "Owner"
const prefixHouse = "House"

type Owner struct {
	Id string //식별자
}

type House struct {
	Id        string
	Address   string
	OwnerId   string
	Price     string
	Timestamp time.Time
}

type HouseContract interface {
	AddOwner(shim.ChaincodeStubInterface, *Owner) error
	CheckOwner(shim.ChaincodeStubInterface, string) (bool, error)
	ListOwners(shim.ChaincodeStubInterface) ([]*Owner, error)

	AddHouse(shim.ChaincodeStubInterface, *House) error
	CheckHouse(shim.ChaincodeStubInterface, string) (bool, error)
	ValidateHouse(shim.ChaincodeStubInterface, *House) (bool, error)
	GetHouse(shim.ChaincodeStubInterface, string) (*House, error)
	UpdateHouse(shim.ChaincodeStubInterface, *House) error
	ListHouses(shim.ChaincodeStubInterface) ([]*House, error)

	TransferHouse(shim.ChaincodeStubInterface, string, string) error
}

type HouseContractCC struct {
}

func checkLen(logger *shim.ChaincodeLogger, expected int, args []string) error {
	if len(args) < expected {
		mes := fmt.Sprintf(
			"not enough number of arguments: %d given, %d expected",
			len(args),
			expected,
		)
		logger.Warning(mes)
		return errors.New(mes)
	}
	return nil
}

func (t *HouseContractCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger := shim.NewLogger("housecontract")
	logger.Info("chaincode initialized")
	return shim.Success([]byte{})
}

func (t *HouseContractCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	var (
		function string
		args     []string
	)
	function, args = stub.GetFunctionAndParameters()
	logger := shim.NewLogger("housecontract")
	logger.Infof("function name = %s", function)
	logger.Infof("args  = %s", args)

	switch function {
	case "AddOwner":
		if err := checkLen(logger, 1, args); err != nil {
			return shim.Error(err.Error())
		}

		goowner := new(Owner)
		err := json.Unmarshal([]byte(args[0]), goowner)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = t.AddOwner(stub, goowner)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success([]byte{})

	case "ListOwners":
		goowners, err := t.ListOwners(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		jsonowners, err := json.Marshal(goowners)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(jsonowners)

	case "AddHouse":
		if err := checkLen(logger, 1, args); err != nil {
			return shim.Error(err.Error())
		}

		gohouse := new(House)
		err := json.Unmarshal([]byte(args[0]), gohouse)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = t.AddHouse(stub, gohouse)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success([]byte{})

	case "ListHouses":
		gohouses, err := t.ListHouses(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		jsonhouses, err := json.Marshal(gohouses)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(jsonhouses)

	case "ListOwnerIdHouses":

		var ownerId string
		err := json.Unmarshal([]byte(args[0]), &ownerId)
		if err != nil {
			return shim.Error(err.Error())
		}

		gohouses, err := t.ListOwnerIdHouses(stub, ownerId)
		if err != nil {
			return shim.Error(err.Error())
		}

		jsonhouses, err := json.Marshal(gohouses)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(jsonhouses)

	case "GetHouse":
		if err := checkLen(logger, 1, args); err != nil {
			return shim.Error(err.Error())
		}

		var houseId string
		err := json.Unmarshal([]byte(args[0]), &houseId)
		if err != nil {
			return shim.Error(err.Error())
		}

		gohouse, err := t.GetHouse(stub, houseId)
		if err != nil {
			return shim.Error(err.Error())
		}

		jsonhouse, err := json.Marshal(gohouse)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(jsonhouse)

	case "UpdateHouse":
		if err := checkLen(logger, 1, args); err != nil {
			return shim.Error(err.Error())
		}

		gohouse := new(House)
		err := json.Unmarshal([]byte(args[0]), gohouse)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = t.UpdateHouse(stub, gohouse)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success([]byte{})

	case "TransferHouse":
		if err := checkLen(logger, 2, args); err != nil {
			return shim.Error(err.Error())
		}

		var houseId, newownerId string
		err := json.Unmarshal([]byte(args[0]), &houseId)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal([]byte(args[1]), &newownerId)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = t.TransferHouse(stub, houseId, newownerId)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte{})
	}

	mes := fmt.Sprintf("Unknown method: %s", function)
	logger.Warning(mes)
	return shim.Error(mes)
}

func (t *HouseContractCC) AddOwner(stub shim.ChaincodeStubInterface,
	goowner *Owner) error {
	logger := shim.NewLogger("AddOwner")
	logger.Infof("AddOwner:  Id = %s", goowner.Id)

	found, err := t.CheckOwner(stub, goowner.Id)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	if found {
		mes := fmt.Sprintf("an Owner with Id = %s alerady exists", goowner.Id)
		logger.Warning(mes)
		return errors.New(mes)
	}

	jsonowner, err := json.Marshal(goowner)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	key, err := stub.CreateCompositeKey(prefixOwner, []string{goowner.Id})
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	err = stub.PutState(key, jsonowner)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	return nil
}

func (t *HouseContractCC) CheckOwner(stub shim.ChaincodeStubInterface,
	id string) (bool, error) {
	logger := shim.NewLogger("CheckOwner")
	logger.Infof("CheckOwner:  Id = %s", id)

	// creates a composite key
	key, err := stub.CreateCompositeKey(prefixOwner, []string{id})
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	// loads from the State DB
	jsonBytes, err := stub.GetState(key)
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	// returns successfully
	return jsonBytes != nil, nil
}

// Lists Owners
func (t *HouseContractCC) ListOwners(stub shim.ChaincodeStubInterface) ([]*Owner,
	error) {
	logger := shim.NewLogger("ListOwners")
	logger.Info("ListOwners")

	// executes a range query, which returns an iterator
	iter, err := stub.GetStateByPartialCompositeKey(prefixOwner, []string{})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	// will close the iterator when returned from this method
	defer iter.Close()
	goowners := []*Owner{}

	// loops over the iterator
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		goowner := new(Owner)
		err = json.Unmarshal(kv.Value, goowner)
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		logger.Infof("Owner Id = %s", goowner.Id)
		goowners = append(goowners, goowner)
	}

	logger.Infof("%d %s found", len(goowners), "Owner")

	return goowners, nil
}

func (t *HouseContractCC) AddHouse(stub shim.ChaincodeStubInterface,
	gohouse *House) error {
	logger := shim.NewLogger("AddHouse")
	logger.Infof("AddHouse:  Id = %s", gohouse.Id)

	key, err := stub.CreateCompositeKey(prefixHouse, []string{gohouse.Id})
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	found, err := t.CheckHouse(stub, gohouse.Id)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	if found {
		mes := fmt.Sprintf("House with Id = %s already exists", gohouse.Id)
		logger.Warning(mes)
		return errors.New(mes)
	}

	ok, err := t.ValidateHouse(stub, gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	if !ok {
		mes := "Validation of the House failed"
		logger.Warning(mes)
		return errors.New(mes)
	}

	jsonhouse, err := json.Marshal(gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	err = stub.PutState(key, jsonhouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	return nil
}

func (t *HouseContractCC) CheckHouse(stub shim.ChaincodeStubInterface, id string) (bool,
	error) {
	logger := shim.NewLogger("CheckHouse")
	logger.Infof("CheckHouse: Id = %s", id)

	key, err := stub.CreateCompositeKey(prefixHouse, []string{id})
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	jsonBytes, err := stub.GetState(key)
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	return jsonBytes != nil, nil
}

func (t *HouseContractCC) ValidateHouse(stub shim.ChaincodeStubInterface,
	gohouse *House) (bool, error) {
	logger := shim.NewLogger("ValidateHouse")
	logger.Infof("ValidateHouse: Id = %s", gohouse.Id)

	found, err := t.CheckOwner(stub, gohouse.OwnerId)
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	return found, nil
}

func (t *HouseContractCC) GetHouse(stub shim.ChaincodeStubInterface,
	id string) (*House, error) {
	logger := shim.NewLogger("GetHouse")

	key, err := stub.CreateCompositeKey(prefixHouse, []string{id})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	jsonBytes, err := stub.GetState(key)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}
	if jsonBytes == nil {
		mes := fmt.Sprintf("House with Id = %s was not found", id)
		logger.Warning(mes)
		return nil, errors.New(mes)
	}

	gohouse := new(House)
	err = json.Unmarshal(jsonBytes, gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	logger.Infof("House Id = %s, OwnerId = %s", gohouse.Id, gohouse.OwnerId)
	return gohouse, nil
}

func (t *HouseContractCC) UpdateHouse(stub shim.ChaincodeStubInterface,
	gohouse *House) error {
	logger := shim.NewLogger("UpdateHouse")
	logger.Infof("UpdateHouse: house = %+v", gohouse)

	found, err := t.CheckHouse(stub, gohouse.Id)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	if !found {
		mes := fmt.Sprintf("House with Id = %s does not exist", gohouse.Id)
		logger.Warning(mes)
		return errors.New(mes)
	}

	ok, err := t.ValidateHouse(stub, gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	if !ok {
		mes := "Validation of the House failed"
		logger.Warning(mes)
		return errors.New(mes)
	}

	key, err := stub.CreateCompositeKey(prefixHouse, []string{gohouse.Id})
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	jsonhouse, err := json.Marshal(gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	err = stub.PutState(key, jsonhouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	return nil
}

func (t *HouseContractCC) ListHouses(stub shim.ChaincodeStubInterface) ([]*House,
	error) {
	logger := shim.NewLogger("ListHouses")
	logger.Info("ListHouses")

	iter, err := stub.GetStateByPartialCompositeKey(prefixHouse, []string{})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	defer iter.Close()

	gohouses := []*House{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		gohouse := new(House)
		err = json.Unmarshal(kv.Value, gohouse)
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		gohouses = append(gohouses, gohouse)
	}

	logger.Infof("%d %s found", len(gohouses), "House")
	return gohouses, nil
}

func (t *HouseContractCC) ListOwnerIdHouses(stub shim.ChaincodeStubInterface, ownerId string) ([]*House,
	error) {
	logger := shim.NewLogger("ListOwnerIdHouses")
	logger.Info("ListOwnerIdHouses")

	// executes a range query, which returns an iterator
	iter, err := stub.GetStateByPartialCompositeKey(prefixHouse, []string{})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	// will close the iterator when returned from this method
	defer iter.Close()

	// loops over the iterator
	gohouses := []*House{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		gohouse := new(House)
		err = json.Unmarshal(kv.Value, gohouse)
		if err != nil {
			logger.Warning(err.Error())
			return nil, err
		}
		if strings.Index(ownerId, "admin") != -1 {
			gohouses = append(gohouses, gohouse)
		} else {
			if gohouse.OwnerId == ownerId {
				gohouses = append(gohouses, gohouse)
			}
		}
	}

	logger.Infof("%d %s found", len(gohouses), "House")
	return gohouses, nil
}

func (t *HouseContractCC) TransferHouse(stub shim.ChaincodeStubInterface, houseId string,
	newownerId string) error {
	logger := shim.NewLogger("TransferHouse")
	logger.Infof("TransferHouse:  House Id = %s, new Owner Id = %s", houseId, newownerId)

	gohouse, err := t.GetHouse(stub, houseId)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}

	gohouse.OwnerId = newownerId

	err = t.UpdateHouse(stub, gohouse)
	if err != nil {
		logger.Warning(err.Error())
		return err
	}
	return nil
}
