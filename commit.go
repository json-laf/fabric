package main
 
import (
        "encoding/json"
        "fmt"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
        "strconv"
        "time"
)

// 定义一个对象，继承合约对象
type Commit struct {
        contractapi.Contract
}

// 上链信息（对象）
type CommitInfo struct {
        ProjectName string `json:"projectName"`
        ProjectDir string `json:"projectDir"`
        Hash      string `json:"hash"`
        Tree    string `json:"tree"`
        Parent     string `json:"parent"`
        Author string `json:"author"`
        Committer   string `json:"committer"`
        Msg   string `json:"msg"`
}

// 查询结果
type QueryResult struct {
        Key    string `json:"Key"`
        Record *CommitInfo
}
type QueryHistoryResult struct {
        TxId string `json:"tx_id"`
        Value string `json:"value"`
        IsDel string `json:"is_del"`
        OnChainTime string `json:"on_chain_time"`
}

// 初始化账本
func (s *Commit) InitLedger(ctx contractapi.TransactionContextInterface) error {
        CommitInfos := []CommitInfo{
                {ProjectName: "test01", ProjectDir: "/home/json/git-test/test01", Hash: "43", Tree: "4312", Parent: "412", Author: "Tom", Committer: "Tom", Msg: "commit01"},
                {ProjectName: "test01", ProjectDir: "/home/json/git-test/test01", Hash: "43", Tree: "4312", Parent: "412", Author: "Tom", Committer: "Tom", Msg: "commit01"},
                {ProjectName: "test01", ProjectDir: "/home/json/git-test/test01", Hash: "43", Tree: "4312", Parent: "412", Author: "Tom", Committer: "Tom", Msg: "commit01"},
                {ProjectName: "test01", ProjectDir: "/home/json/git-test/test01", Hash: "43", Tree: "4312", Parent: "412", Author: "Tom", Committer: "Tom", Msg: "commit01"},
        }
        for _, CommitInfo := range CommitInfos {
                CommitInfoAsBytes, _ := json.Marshal(CommitInfo)
                err := ctx.GetStub().PutState(CommitInfo.Hash, CommitInfoAsBytes)
                if err != nil {
                        return fmt.Errorf("Failed to put to world state. %s", err.Error())
                }
        }
        return nil
}

//写入commit信息
func (s *Commit) CreateStudentInfo(ctx contractapi.TransactionContextInterface, projectName string, projectDir string, hash string, tree string, parent string, author string, committer string, msg string) error {
        CommitInfo := CommitInfo{
                ProjectName: projectName,
                ProjectDir: projectDir,
                Hash:      hash,
                Tree:    tree,
                Parent:     parent,
                Author: author,
                Committer:   committer,
                Msg: msg,
        }
        CommitInfoAsBytes, _ := json.Marshal(CommitInfo)
        return ctx.GetStub().PutState(CommitInfo.Hash, CommitInfoAsBytes)
}

//查询commit信息
func (s *Commit) QueryStudentInfo(ctx contractapi.TransactionContextInterface, CommitInfoHash string) (*CommitInfo, error) {
        CommitInfoAsBytes, err := ctx.GetStub().GetState(CommitInfoHash)
        if err != nil {
                return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
        }
        if CommitInfoAsBytes == nil {
                return nil, fmt.Errorf("%s does not exist", CommitInfoHash)
        }
        commitInfo := new(CommitInfo)
        //注意： Unmarshal(data []byte, v interface{})的第二个参数为指针类型（结构体地址）
        err = json.Unmarshal(CommitInfoAsBytes, commitInfo) //stuInfo := new(StudentInfo)，stuInfo本身就是指针
        if err != nil {
                return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
        }
        return commitInfo, nil //&取地址,*取指针;指针类型前面加上*号（前缀）来获取指针所指向的内容
}
 
// 查询commit信息（查询的key末尾是数字，有对应的区间）
func (s *Commit) QueryAllStudentInfos(ctx contractapi.TransactionContextInterface, startId, endId string) ([]QueryResult, error) {
        resultsIterator, err := ctx.GetStub().GetStateByRange(startId, endId)
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
                CommitInfo := new(CommitInfo)
                _ = json.Unmarshal(queryResponse.Value, CommitInfo)
 
                queryResult := QueryResult{Key: queryResponse.Key, Record: CommitInfo}
                results = append(results, queryResult)
        }
        return results, nil
}
 
// 修改学生信息
func (s *Commit) ChangeStudentInfo(ctx contractapi.TransactionContextInterface,  hash string, tree string, parent string, author string, committer string, msg string) error {
        commitInfo, err := s.QueryStudentInfo(ctx, hash)
        if err != nil {
                return err
        }
        commitInfo.Hash = hash
        commitInfo.Tree = tree
        commitInfo.Parent = parent
        commitInfo.Author = author
        commitInfo.Committer = committer
        CommitInfoAsBytes, _ := json.Marshal(commitInfo)
        return ctx.GetStub().PutState(hash, CommitInfoAsBytes)
}

//获取历史信息
func (s *Commit) GetHistory(ctx contractapi.TransactionContextInterface, hash string) ([]QueryHistoryResult, error) {
        resultsIterator, err := ctx.GetStub().GetHistoryForKey(hash)
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
                        res.OnChainTime=time.Unix(queryResponse.Timestamp.Seconds,0).Format("2020-01-27 15:04:05")
                        results= append(results, res)
                }
                if err!=nil {
                        return nil,err
                }
        }
        return results, nil
}
 
func main() {
        chaincode, err := contractapi.NewChaincode(new(Commit))
        if err != nil {
                fmt.Printf("Error create fabStudentInfo chaincode: %s", err.Error())
                return
        }
        if err := chaincode.Start(); err != nil {
                fmt.Printf("Error starting fabStudentInfo chaincode: %s", err.Error())
        }
}