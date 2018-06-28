package main

import (
        "encoding/json"
        "fmt"
        _ "strconv"
        _ "strings"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var key, value string
        var err error
        fmt.Println("starting write")

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
        }

        key = args[0]                                   //rename for funsies
        value = args[1]
        err = stub.PutState(key, []byte(value))         //write the variable into the ledger
        if err != nil {
                return shim.Error(err.Error())
        }

        fmt.Println("- end write")
        return shim.Success(nil)
}

// ============================================================================================================================
// delete_info() - remove a marble from state and from marble index
// ============================================================================================================================
func delete_info(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
        fmt.Println("starting delete_info")

        id := args[0]

        // get the info
        _ , err := get_info(stub, id) // _ 를 변수로 한 이유는 에러확인 유무만 하고 해당 변수를 사용하지 않기 때문이다.
        if err != nil{
                fmt.Println("Failed to find info by id " + id)
                return shim.Error(err.Error())
        }

        // remove the marble
        err = stub.DelState(id)                                                 //remove the key from chaincode state
        if err != nil {
                return shim.Error("Failed to delete state")
        }

        fmt.Println("- end delete_info")
        return shim.Success(nil)
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
// type Info struct {
//      Id       string         `json:"id"`
//      Name      string        `json:"name"`
//      Phone       string      `json:"phone"`
//      Address      string     `json:"address"`
// ============================================================================================================================
func init_info(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
        var err error
        fmt.Println("starting init_info")

        if len(args) != 4 {
                return shim.Error("Incorrect number of arguments. Expecting 4")
        }

        id := args[0]
        name := args[1]
        phone := args[2]
        address := args[3]

        //check if info id already exists
        info, err := get_info(stub, id)
        if err == nil {
                fmt.Println("This info already exists - " + id)
                fmt.Println(info)
                return shim.Error("This marble already exists - " + id)  // 정보가 존재하면 중지
        }

        //build the info  json string manually
        str := `{
                "id": "` + id + `",
                "name": "` + name + `",
                "phone": "` + phone + `",
                "address": "` + address + `"
        }`
        err = stub.PutState(id, []byte(str))
        if err != nil {
                return shim.Error(err.Error())
        }

        fmt.Println("- end init_info")
        return shim.Success(nil)
}

// ============================================================================================================================
// modify   //   id, 수정할 내용, 데이터 를 입력 받는다..   변수 3개를 받는다.
// ============================================================================================================================
func modify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        fmt.Println("starting modify")

        if len(args) != 3 {
                return shim.Error("Incorrect number of arguments. Expecting 3")
        }

        var id = args[0]
        var type_d = args[1]
        var data = args[2]

        fmt.Println( "Modify " + type_d + " Data -> " + data + " ...")

        // 현재 ID의 정보를 가져옴
        infoAsBytes, err := stub.GetState(id)
        if err != nil {
                return shim.Error("Failed to get Info")
        }
        res := Info{}
        json.Unmarshal(infoAsBytes, &res)           //un stringify it aka JSON.parse()

        // transfer
        if type_d == "name" {
                fmt.Println( "Current " + type_d + " : " + res.Name + " => " + data + "...")
                res.Name = data
        } else if type_d == "phone" {
                fmt.Println( "Current " + type_d + " : " + res.Phone + " => " + data + "...")
                res.Phone = data
        } else if type_d == "address" {
                fmt.Println( "Current " + type_d + " : " + res.Address + " => " + data + "...")
                res.Address = data
        }

        jsonAsBytes, _ := json.Marshal(res)           //convert to array of bytes
        err = stub.PutState(id, jsonAsBytes)
        if err != nil {
                return shim.Error(err.Error())
        }

        fmt.Println("- end set " + type_d + " data.")
        return shim.Success(nil)
}
