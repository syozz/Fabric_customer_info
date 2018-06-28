package main

import (
        "encoding/json"
        "errors"
//      "strconv"

        "github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Get Info - get a info asset from ledger  정보를 가져온다...
// ============================================================================================================================
func get_info(stub shim.ChaincodeStubInterface, id string) (Info, error) {
        var info Info
        InfoAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
        if err != nil {                                          //this seems to always succeed, even if key didn't exist
                return info, errors.New("Failed to find info - " + id)
        }
        // JSON 으로 인코딩 된 데이터를 디코딩 한다, 첫 파라미터에는 JSON 데이터
        // 두번째 파라미터에는 출력할 구조체를 포인터로 지정.....
        json.Unmarshal(InfoAsBytes, &info)                   //un stringify it aka JSON.parse()

        if info.Id != id {                                     //test if info is actually here or just nil
                return info, errors.New("Info does not exist - " + id)
        }

        return info, nil
}



