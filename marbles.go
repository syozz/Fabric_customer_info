package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"time"
	//      "github.com/hyperledger/fabric/common/util"
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
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	Payment_plan   string    `json:"payment_plan"`   //  07.02  데이터 추가
	Extra_plan     string    `json:"extra_plan"`     //  07.02  데이터 추가
	Method_payment string    `json:“method_payment"` // 결제방식 [ 카드, 계좌이체, 무통장 ]
	Grade          string    `json:"grade"`          //  07.02  데이터 추가
	Modified_time  time.Time `json:"modified_time"`  //  07.02  데이터 추가

}
type Plan_Info struct {
	Plan_name string `json:"plan_name"`
	Fee       string `json:"fee"`
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

	fmt.Println("Ready for action") //self-test pass
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
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "write" { //generic writes to ledger
		return write(stub, args)
	} else if function == "delete_info" { //deletes a info from state
		return delete_info(stub, args)
	} else if function == "init_info" { //create a new info
		return init_info(stub, args)
	} else if function == "modify" { //change owner of a info
		return modify(stub, args)
	} else if function == "getHistory" { //read history of a info (audit)
		return getHistory(stub, args)
	} else if function == "read" {
		return read(stub, args)
	} else if function == "query_fee" {
		// 전달할 매개변수
		// 체인코드명, 체인코드변수, 채널명
		info, err := get_info(stub, args[0])
		if err == nil {
			fmt.Println("This info exists - " + args[0])
		}

		var plan_name = info.Payment_plan
		trans := [][]byte{[]byte("read"), []byte(plan_name)}
		response := stub.InvokeChaincode("mycplan", trans, "mycplan")

		if response.Status != shim.OK {
			return shim.Error(response.Message)
		}

		res := Plan_Info{}
		json.Unmarshal(response.Payload, &res)

		fmt.Println("Current Plan : " + res.Plan_name + "...")
		fmt.Println("Current Fee : " + res.Fee + "...")

		return shim.Success([]byte(response.Payload))
	} else if function == "payment" {
		// 전달할 매개변수
		// 체인코드명, 체인코드변수, 채널명
		cusid := args[0]
		info, err := get_info(stub, cusid)
		if err == nil {
			fmt.Println("This info exists - " + cusid)
		}

		var payment_plan = info.Payment_plan
		// var extra_plan = info.Extra_plan
		var method_payment = info.Method_payment
		// 요금제 요금 조회
		trans := [][]byte{[]byte("read"), []byte(payment_plan)}
		response := stub.InvokeChaincode("mycplan", trans, "mycplan")

		if response.Status != shim.OK {
			return shim.Error(response.Message)
		}

		res := Plan_Info{}
		json.Unmarshal(response.Payload, &res)
		fee := res.Fee

		fmt.Println("Current Plan : " + res.Plan_name + "...")
		fmt.Println("Current Fee : " + res.Fee + "...")
		fmt.Println(" ## START PAYMENT ## ")

		// 결제를 위한 값 입력
		// type PaymentInfo struct {
		//         CusId                string   `json:"cusid"`  // 고객아이디
		//         Month                string   `json:"month"`  // 결제 월
		//         Payment_plan         string   `json:"month"`  // 가입요금제
		//         Extra_plan           string   `json:“extra_plan＂`   //  부가서비스
		//         Amount_payment       string   `json:"amount_payment"`  // 결제금액
		//         Method_payment       string   `json:“method_payment"`  // 결제방식 [ 카드, 계좌이체, 무통장 ]
		//         Result               string   `json:“result"`  // 결제 결과 [ 성공, 실패 ]
		// }
		t := time.Now()
		t.Format(time.RFC3339)
		y := t.Year()
		mon := t.Month()
		d := t.Day()
		h := t.Hour()
		m := t.Minute()
		s := t.Second()
		n := t.Nanosecond()

		fmt.Println("Year   :", y)
		fmt.Println("Month   :", mon)
		fmt.Println("Day   :", d)
		fmt.Println("Hour   :", h)
		fmt.Println("Minute :", m)
		fmt.Println("Second :", s)
		fmt.Println("Nanosec:", n)

		trans_payment := [][]byte{[]byte("payment"),
			[]byte(cusid),
			[]byte(strconv.Itoa(y)),
			[]byte(string(mon)),
			[]byte(payment_plan),
			[]byte("NONE"),
			[]byte(fee),
			[]byte(method_payment)}

		payment_response := stub.InvokeChaincode("mycpayment", trans_payment, "mycpayment")
		if payment_response.Status != shim.OK {
			return shim.Error(payment_response.Message)
		}

		return shim.Success([]byte(payment_response.Payload))
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
