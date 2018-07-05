package main

import (
        "fmt"
        "time"
//      "strconv"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Asset Definitions - The ledger will store marbles and owners
// ============================================================================================================================

// ----- Marbles ----- //
//구조체로 관리,...
type Info struct {
        Id              string  `json:"id"`
        Name            string  `json:"name"`
        Phone           string  `json:"phone"`
        Address         string  `json:"address"`
        Payment_plan    string  `json:"payment_plan"`   //  07.02  데이터 추가
        Grade           string  `json:"grade"`          //  07.02  데이터 추가
        Modified_time   time.Time       `json:"modified_time"`  //  07.02  데이터 추가

}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
        err := shim.Start(new(SimpleChaincode))
        if err != nil {
                fmt.Printf("Error starting Simple chaincode - %s", err)
        }
}


// ============================================================================================================================
// Init - initialize the chaincode   /// 초기화 인데 자세히  무슨역할을 하는 부분인지 좀 더 확인이 필요할듯...
//
// Marbles does not require initialization, so let's run a simple test instead.
//
// Shows off PutState() and how to pass an input argument to chaincode.
// Shows off GetFunctionAndParameters() and GetStringArgs()
// Shows off GetTxID() to get the transaction ID of the proposal
//
// Inputs - Array of strings
//  ["314"]
//
// Returns - shim.Success or error
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

        var err error

        fmt.Println("Cus Info Starting Up")

        err = stub.PutState("admin", []byte("Init"))
        if err != nil {
                return shim.Error(err.Error())
        }

        fmt.Println("Ready for action")                          //self-test pass
        return shim.Success(nil)
}


// ============================================================================================================================
// Invoke - Our entry point for Invocations 실제 명령을 날리는 부분..
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
        fmt.Println(" ")
        fmt.Println("starting invoke, for - " + function)

        // Handle different functions
        if function == "init" {                    //initialize the chaincode state, used as reset
                return t.Init(stub)
        } else if function == "write" {            //generic writes to ledger
                return write(stub, args)
        } else if function == "delete_info" {    //deletes a info from state
                return delete_info(stub, args)
        } else if function == "init_info" {      //create a new info
                return init_info(stub, args)
        } else if function == "modify" {        //change owner of a info
                return modify(stub, args)
        } else if function == "getHistory"{        //read history of a info (audit)
                return getHistory(stub, args)
        } else if function == "read"{
                return read(stub, args)
        }

        // error out
        fmt.Println("Received unknown invoke function name - " + function)
        return shim.Error("Received unknown invoke function name - '" + function + "'")
}


// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
        return shim.Error("Unknown supported call - Query()")
}
