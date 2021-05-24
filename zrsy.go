package main
 
import (
        "encoding/json"
        "fmt"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
        "strconv"
        "time"
)

// 定义一个对象，继承合约对象
type Logistics struct {
        contractapi.Contract
}

// 上链信息（对象）
type LogisticsInfo struct {
        Number          string  `json:"number"`
        Company         string  `json:"company"`
        Time            string  `json:"time"`
        Car             string  `json:"car"`
        Employee        string  `json:"employee"`
        Temperature     string  `json:"temperature"`
        Route           string  `json:"route"`
}
 
// QueryResult structure used for handling result of query
type QueryResult struct {
        Key    string `json:"Key"`
        Record *LogisticsInfo
}
type QueryHistoryResult struct {
        TxId            string  `json:"tx_id"`
        Value           string  `json:"value"`
        IsDel           string  `json:"is_del"`
        OnChainTime     string  `json:"on_chain_time"`
}
// 初始化账本
func (s *Logistics) InitLedger(ctx contractapi.TransactionContextInterface) error {
        LogisticsInfos := []LogisticsInfo{
                {Number: "01", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "02", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "03", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "04", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "05", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "06", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "07", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "08", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "09", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
                {Number: "10", Company: "张三", Time: "23", Car: "北京", Employee: "北京", Temperature: "北京", Route: "北京"},
        }
        for _, LogisticsInfo := range LogisticsInfos {
                LogisticsInfoAsBytes, _ := json.Marshal(LogisticsInfo)
                err := ctx.GetStub().PutState(LogisticsInfo.Number, LogisticsInfoAsBytes)
                if err != nil {
                        return fmt.Errorf("Failed to put to world state. %s", err.Error())
                }
        }
        return nil
}
 
// 上链学生信息
func (s *Logistics) CreateLogisticsInfo(ctx contractapi.TransactionContextInterface, number string, company string, time string, car string, employee string, temperature string, route string) error {
        LogisticsInfo := LogisticsInfo{
                Number:         number,
                Company:        company,
                Time:           time,
                Car:            car,
                Employee:       employee,
                Temperature:    temperature,
                Route:          route,
        }
        LogisticsInfoAsBytes, _ := json.Marshal(LogisticsInfo)
        return ctx.GetStub().PutState(LogisticsInfo.Number, LogisticsInfoAsBytes)
}
 
//查询学生信息
func (s *Logistics) QueryLogisticsInfo(ctx contractapi.TransactionContextInterface, LogisticsInfoNumber string) (*LogisticsInfo, error) {
        LogisticsInfoAsBytes, err := ctx.GetStub().GetState(LogisticsInfoNumber)
        if err != nil {
                return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
        }
        if LogisticsInfoAsBytes == nil {
                return nil, fmt.Errorf("%s does not exist", LogisticsInfoNumber)
        }
        logInfo := new(LogisticsInfo)
        //注意： Unmarshal(data []byte, v interface{})的第二个参数为指针类型（结构体地址）
        err = json.Unmarshal(LogisticsInfoAsBytes, logInfo) //logInfo := new(LogisticsInfo)，logInfo本身就是指针
        if err != nil {
                return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
        }
        return logInfo, nil
}

// 查询学生信息（查询的key末尾是数字，有对应的区间）
func (s *Logistics) QueryAllLogisticsInfos(ctx contractapi.TransactionContextInterface, startNum, endNum string) ([]QueryResult, error) {
        resultsIterator, err := ctx.GetStub().GetStateByRange(startNum, endNum)
        if err != nil {
                return nil, err
        }
        defer resultsIterator.Close()
        results := []QueryResult{}
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()

                if err != nil {
                        return nil, err
                }
                LogisticsInfo := new(LogisticsInfo)
                _ = json.Unmarshal(queryResponse.Value, LogisticsInfo)

                queryResult := QueryResult{Key: queryResponse.Key, Record: LogisticsInfo}
                results = append(results, queryResult)
        }
        return results, nil
}
 
// 修改学生信息
func (s *Logistics) ChangeLogisticsInfo(ctx contractapi.TransactionContextInterface, number string, company string, time string, car string, employee string, temperature string, route string) error {
        logInfo, err := s.QueryLogisticsInfo(ctx, number)
        if err != nil {
                return err
        }
        logInfo.Number = number
        logInfo.Company = company
        logInfo.Time = time
        logInfo.Car = car
        logInfo.Employee = employee
        logInfo.Temperature = temperature
        logInfo.Route = route
        LogisticsInfoAsBytes, _ := json.Marshal(logInfo)
        return ctx.GetStub().PutState(number, LogisticsInfoAsBytes)
}
 
//获取历史信息
func (s *Logistics) GetHistory(ctx contractapi.TransactionContextInterface, number string) ([]QueryHistoryResult, error) {
        resultsIterator, err := ctx.GetStub().GetHistoryForKey(number)
        if err != nil {
                return nil, err
        }
        defer resultsIterator.Close()
        //results := []QueryResult{}
        //results := make([]QueryResult, 0)
        results := make([]QueryHistoryResult, 0)
        for resultsIterator.HasNext() {
                if queryResponse, err := resultsIterator.Next();err==nil{
                        res := QueryHistoryResult{}
                        res.TxId=queryResponse.TxId
                        res.Value=string(queryResponse.Value)
                        res.IsDel=strconv.FormatBool(queryResponse.IsDelete)
                        res.OnChainTime=time.Unix(queryResponse.Timestamp.Seconds,0).Format("2006-01-02 15:04:05")
                        results= append(results, res)
                }
                if err!=nil {
                        return nil,err
                }
        }
        return results, nil
}
 
func main() {
        chaincode, err := contractapi.NewChaincode(new(Logistics))
        if err != nil {
                fmt.Printf("Error create fabLogisticsInfo chaincode: %s", err.Error())
                return
        }
        if err := chaincode.Start(); err != nil {
                fmt.Printf("Error starting fabLogisticsInfo chaincode: %s", err.Error())
        }
}