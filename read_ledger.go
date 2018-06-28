package main

import (
//      "bytes"
        "encoding/json"
        "fmt"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)


// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        type AuditHistory struct {
                TxId    string   `json:"txId"`
                Value   Info   `json:"value"`
        }
        var history []AuditHistory;
        var info Info

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        infoId := args[0]
        fmt.Printf("- start getHistoryFor Customer Info: %s\n", infoId)

        // Get History
    // 마블 id 를 받아서 stub.GetHistoryForKey 메소드 호출,,,
        resultsIterator, err := stub.GetHistoryForKey(infoId)
        if err != nil {
                return shim.Error(err.Error())
        }
        // defer는 지연처리 하는 것,,, 현재 실행되는 func 이 종료될때 실행이 되게 함...
        // 현재 func 이 종료될때 resultsIterator 객체? 를 닫으라는 명령
        defer resultsIterator.Close()

        // resultsIterator 의 데이터를 하나씩 순차적으로 가져옴..
        // Iterator 의 역할은 저장된 데이터를 순차적으로 가져오거나 삭제할수 있는 기능 제공
        for resultsIterator.HasNext() {
                historyData, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }

                var tx AuditHistory
                tx.TxId = historyData.TxId                     //copy transaction id over
                json.Unmarshal(historyData.Value, &info)     //un stringify it aka JSON.parse()
                if historyData.Value == nil {                  //marble has been deleted
                        var emptyInfo Info
                        tx.Value = emptyInfo                 //copy nil marble
                } else {
                        json.Unmarshal(historyData.Value, &info) //un stringify it aka JSON.parse()
                        tx.Value = info                      //copy marble over
                }
                history = append(history, tx)              //add this tx to the list
        }
        fmt.Printf("- getHistoryFor Customer Info returning:\n%s", history)

        //change to array of bytes
        historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
        return shim.Success(historyAsBytes)
}
