package main

import (
	"crypto/sha256"
	"eth-relay/go-sdk/storeKey"
	"eth-relay/go-sdk/tool"
	"eth-relay/go-sdk/utils"
	"fmt"
	"git.huawei.com/poissonsearch/wienerchain/proto/common"
	"git.huawei.com/poissonsearch/wienerchain/wienerchain-go-sdk/client"
	t "git.huawei.com/poissonsearch/wienerchain/wienerchain-go-sdk/utils"
	"github.com/Yinlianlei/sql"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)


func groupSig()  {
	temp1,_ := new(big.Int).SetString("1515",10)

	temp2,_ := new(big.Int).SetString("151515",10)

	//temp1.GCD()
	temp2.Sub(temp2,temp1)

	fmt.Println(temp1)
	fmt.Println(temp2)
}

/*
	通过智能合约发送交易的函数
 */
func sendCopyTransaction(contract string,function string,copy *storeKey.CopyRight) (*common.Response,error,string) {
	configPath := storeKey.ConfigPath
	chainID := storeKey.ChainID//链名称
	//contract := "huaweichaincode"//智能合约（链码
	//function := "writeMarble"//调用的链码函数
	nodeName := storeKey.NodeName

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		return nil,err,""
	}

	txID, err := t.GenerateTxID()//获取交易ID
	copy.TxId = txID

	if err != nil {
		return nil,err,""
	}

	//返回RawMessage 合约交易需发送的消息
	var args []string
	s,err := tool.SpliceJson(copy)
	if err!=nil {
		return nil,err,""
	}
	//使用版权哈希作为Key
	args = append(args,copy.Hash,s)

	rawMsg, err := gatewayClient.ContractRawMessage.BuildInvokeMessage(chainID, txID, contract, function, args)
	if err != nil {
		return nil,err,""
	}

	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		return nil,err,""
	}

	invokeResponse, err := node.ContractAction.Invoke(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		return nil,err,""
	}

	var invokeResponses []*common.RawMessage
	invokeResponses = append(invokeResponses,invokeResponse)
	transactionRawMsg, err := gatewayClient.ContractRawMessage.BuildTransactionMessage(invokeResponses)
	if err != nil {
		return nil,err,""
	}

	transactionResponse, err := node.ContractAction.Transaction(transactionRawMsg)
	if err != nil {
		return nil,err,""
	}

	txResponse := &common.Response{}
	if err := proto.Unmarshal(transactionResponse.Payload, txResponse); err != nil {
		return nil,err,""
	}

	if txResponse.Status == common.Status_SUCCESS {
		//写入数据库
		return txResponse,nil,txID
	} else {
		return txResponse,errors.New(string(txResponse.Payload)),""
	}
}

func sendPurcTransaction(contract string,function string,pur *storeKey.Purchase) (*common.Response,error,string) {
	configPath := storeKey.ConfigPath
	chainID := storeKey.ChainID//链名称
	//contract := "huaweichaincode"//智能合约（链码
	//function := "writeMarble"//调用的链码函数
	nodeName := storeKey.NodeName

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		return nil,err,""
	}

	txID, err := t.GenerateTxID()//获取交易ID
	pur.TxId = txID

	if err != nil {
		return nil,err,""
	}

	//返回RawMessage 合约交易需发送的消息
	var args []string
	s,err := tool.SpliceJson(pur)
	if err!=nil {
		return nil,err,""
	}
	bs := []byte(s)
	hash := sha256.Sum256(bs)
	args = append(args,common2.BytesToHash(hash[:]).String(),s)

	rawMsg, err := gatewayClient.ContractRawMessage.BuildInvokeMessage(chainID, txID, contract, function, args)
	if err != nil {
		return nil,err,""
	}
	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		return nil,err,""
	}

	invokeResponse, err := node.ContractAction.Invoke(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		return nil,err,""
	}

	var invokeResponses []*common.RawMessage
	invokeResponses = append(invokeResponses,invokeResponse)
	transactionRawMsg, err := gatewayClient.ContractRawMessage.BuildTransactionMessage(invokeResponses)
	if err != nil {
		return nil,err,""
	}

	transactionResponse, err := node.ContractAction.Transaction(transactionRawMsg)
	if err != nil {
		return nil,err,""
	}

	txResponse := &common.Response{}
	if err := proto.Unmarshal(transactionResponse.Payload, txResponse); err != nil {
		return nil,err,""
	}

	if txResponse.Status == common.Status_SUCCESS {
		//写入数据库

		return txResponse,nil,txID
	} else {
		return txResponse,errors.New(string(txResponse.Payload)),""
	}
}

/*
	使用Key来查询账本内容
	需要合约名 contract
	函数名 function
	参数 args
	接受结果的数据结构 a
 */
func queryTransaction(contract string,function string,args []string,a interface{})  (error)  {
	configPath := storeKey.ConfigPath
	chainID := storeKey.ChainID//链名称
	//contract := "huaweichaincode"//智能合约（链码
	//function := "writeMarble"//调用的链码函数
	nodeName := storeKey.NodeName

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		return err
	}

	txID, err := t.GenerateTxID()//获取交易ID
	if err != nil {
		return err
	}

	//返回RawMessage 合约交易需发送的消息
	rawMsg, err := gatewayClient.ContractRawMessage.BuildInvokeMessage(chainID, txID, contract, function, args)
	if err != nil {
		return err
	}
	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		return err
	}

	invokeResponse, err := node.ContractAction.Invoke(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		return err
	}

	response := &common.Response{}
	if err := proto.Unmarshal(invokeResponse.Payload, response); err != nil {}
	if response.Status == common.Status_SUCCESS {
		tx := &common.Transaction{}
		if err := proto.Unmarshal(response.Payload, tx); err != nil {
			return err
		}
		txPayLoad := &common.TxPayload{}
		if err := proto.Unmarshal(tx.Payload, txPayLoad); err != nil {
			return err
		}
		txData := &common.CommonTxData{}
		if err := proto.Unmarshal(txPayLoad.Data, txData); err != nil {
			return err
		}

		if err:=tool.ParseJson(string(txData.Response.Payload),&a);err!=nil{
			return err
		}
		return nil
	} else {
		return errors.New(response.StatusInfo)
	}

}

/*
	查询交易相关信息
 */
func queryById(txID string) (error,string) {
	configPath := storeKey.ConfigPath
	chainID := storeKey.ChainID//链名称
	//contract := "huaweichaincode"//智能合约（链码
	//function := "writeMarble"//调用的链码函数
	nodeName := storeKey.NodeName

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		return err,""
	}

	rawMsg, err := gatewayClient.QueryRawMessage.BuildTxRawMessage(chainID,txID)
	if err != nil {
		return err,""
	}
	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		return err,""
	}

	invokeResponse, err := node.QueryAction.GetTxByTxID(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		return err,""
	}

	response := &common.Response{}
	if err := proto.Unmarshal(invokeResponse.Payload, response); err != nil {}
	if response.Status == common.Status_SUCCESS {
		tx := &common.Transaction{}
		if err := proto.Unmarshal(response.Payload, tx); err != nil {
			return err,""
		}
		txPayLoad := &common.TxPayload{}
		if err := proto.Unmarshal(tx.Payload, txPayLoad); err != nil {
			return err,""
		}
		//fmt.Println(txPayLoad.Header)
		//fmt.Println(string(txPayLoad.Data))
		txData := &common.CommonTxData{}
		if err := proto.Unmarshal(txPayLoad.Data, txData); err != nil {
			return err,""
		}

		/*
		c := &common.ContractInvocation{}
		if err := proto.Unmarshal(txData.ContractInvocation, c); err != nil {
			return err,""
		}
		fmt.Println(c.ContractName)
		for i := 0; i < len(c.Args); i++ {
			fmt.Println(string(c.Args[i]))
		}
		fmt.Println(c.FuncName) */

		return nil,string(txData.Response.Payload)
	} else {
		return errors.New(response.StatusInfo),""
	}
}

func queryResultById(txID string) (error,string) {
	configPath := storeKey.ConfigPath
	chainID := storeKey.ChainID//链名称
	//contract := "huaweichaincode"//智能合约（链码
	//function := "writeMarble"//调用的链码函数
	nodeName := storeKey.NodeName

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		return err,""
	}

	rawMsg, err := gatewayClient.QueryRawMessage.BuildTxRawMessage(chainID,txID)
	if err != nil {
		return err,""
	}
	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		return err,""
	}

	invokeResponse, err := node.QueryAction.GetVote(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		return err,""
	}



	response := &common.Response{}
	if err := proto.Unmarshal(invokeResponse.Payload, response); err != nil {}
	if response.Status == common.Status_SUCCESS {
		tx := &common.Transaction{}
		if err := proto.Unmarshal(response.Payload, tx); err != nil {
			return err,""
		}
		txPayLoad := &common.TxPayload{}
		if err := proto.Unmarshal(tx.Payload, txPayLoad); err != nil {
			return err,""
		}
		txData := &common.CommonTxData{}
		if err := proto.Unmarshal(txPayLoad.Data, txData); err != nil {
			return err,""
		}

		return nil,string(txData.Response.Payload)
	} else {
		return errors.New(response.StatusInfo),""
	}
}

/*
	接收微信的处理数据
 */
func dealpost(c *gin.Context)  {
	json := make(map[string]string) //注意该结构接受的内容
	c.BindJSON(&json)

	//验证账号登录情况
	username := json["owner"]
	pwd := json["pwd"]

	err1 := sql.Sql_Wx_Login(username,pwd)

	if err1 != nil{
		c.JSON(400, gin.H{
			"info":err1.Error(),
		})
		return
	}

	var args storeKey.Args
	var err error

	ty := json["type"]
	fun := json["function"]
	arg := json["args"]
	temp := arg[1:len(arg)-1]
	args.A = strings.Split(temp,",")
	/*
	err := tool.ParseJson(arg,&args)//解析参数
	if err!=nil {
		c.JSON(400, gin.H{
			"info":err.Error(),
		})
		return
	}
	 */


	if ty=="query" {
		if fun=="认证信息" {
			var copy storeKey.CopyRight
			err := queryTransaction("huaweichaincode","getMarble",args.A,&copy)
			//返回交易信息
			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}
			str,err := tool.SpliceJson(copy)

			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"info":str,
			})
			return
		}else if fun=="购买信息" {
			var pur storeKey.Purchase
			err := queryTransaction("huaweichaincode","getMarble",args.A,&pur)
			//返回交易信息
			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}

			str,err := tool.SpliceJson(pur)

			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"info":str,
			})
			return
		}
	}else if ty=="send" {
		if fun=="认证" {
			if len(args.A)<4 {
				c.JSON(400, gin.H{
					"info":"参数个数不匹配",
				})
				return
			}

			var copy storeKey.CopyRight

			//处理文件哈希生成版权哈希

			str := args.A[0]

			copy.Time = time.Now().UnixNano()
			copy.Owner = args.A[1]

			str = str + copy.Owner + strconv.Itoa(int(copy.Time))
			bstr := []byte(str)
			hash := sha256.Sum256(bstr)
			copy.Hash = common2.BytesToHash(hash[:]).String()

			copy.Filename = args.A[2]
			copy.FileID = args.A[3]

			response,err,_ := sendCopyTransaction("huaweichaincode","writeMarble",&copy)

			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}

			if response.Status!= common.Status_SUCCESS{
				c.JSON(400, gin.H{
					"info":response.Status,
					"key":copy.Hash,
				})
				return
			}

			sql.Sql_Wx_ConfirmCopyright(copy)
			c.JSON(200, gin.H{
				"info":copy.TxId,
			})
			return

		}else if fun=="购买"{
			if len(args.A)<4 {
				c.JSON(400, gin.H{
					"info":"参数个数不匹配",
				})
				return
			}

			var pur storeKey.Purchase
			pur.Buyer = args.A[0]
			pur.Owner = args.A[1]
			pur.Hash = args.A[2]
			pur.Price,err = strconv.ParseFloat(args.A[3],64)
			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}
			response,err,_ := sendPurcTransaction("huaweichaincode","writeMarble",&pur)

			if err!=nil {
				c.JSON(400, gin.H{
					"info":err.Error(),
				})
				return
			}

			if response.Status!= common.Status_SUCCESS{
				c.JSON(400, gin.H{
					"info":response.Status,
					"key":pur.Hash,
				})
				return
			}

			sql.Sql_Wx_Purchase(pur)
			c.JSON(200, gin.H{
				"info":pur.TxId,
			})
			return

		}else {
			if err != nil{
				c.JSON(400, gin.H{
					"info":"函数错误",
				})
				return
			}
		}
	}else {
		if err != nil{
			c.JSON(400, gin.H{
				"info":"type错误",
			})
			return
		}
	}

}

/*
	处理上传的数据
	定制型服务
 */
func dealVipFile(c *gin.Context)  {
	//验证权限

	/*
		文件存储测试
	*/
	file,err := c.FormFile("file");

	if err!=nil {
		fmt.Println(err.Error())
	}else {
		filepath := "./files/"+file.Filename
		c.SaveUploadedFile(file,filepath)

		//读取文件
		bytes, err := ioutil.ReadFile(filepath)
		if err!=nil {
			c.JSON(400,gin.H{
				"info":err.Error(),
			})
		}
		//文件加密
		re,_,err := tool.CreateSm2Encrypt(bytes)
		if err!=nil {
			c.JSON(400,gin.H{
				"info":err.Error(),
			})
			return
		}
		strre,err := tool.Sm2Decrypt(re)
		if err!=nil {
			c.JSON(400,gin.H{
				"info":err.Error(),
			})
			return
		}
		//删除原有文件
		err = os.Remove(filepath)
		if err!=nil {
			c.JSON(400,gin.H{
				"info":err.Error(),
			})
			return
		}
		//输出文件
		err = ioutil.WriteFile(filepath,[]byte(strre),0777)
		if err!=nil {
			c.JSON(400,gin.H{
				"info":err.Error(),
			})
		}
		c.JSON(200,gin.H{
			"info":"操作成功",
		})

	}


}

/*
	校验权限
 */
func authIdentification() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		//处理逻辑
		if true {
			//c.Next()
			c.JSON(200,gin.H{
				"info":"ss",
			})
		}
	}
}

func invoke()  {



	/*
	//处理返回的交易数据
	response := &common.Response{}
	if err := proto.Unmarshal(invokeResponse.Payload, response); err != nil {
		fmt.Printf("unmarshal invoke response error： %v", err)
	}
	//检查交易是否成功
	if response.Status == common.Status_SUCCESS {
		tx := &common.Transaction{}
		if err := proto.Unmarshal(response.Payload, tx); err != nil {
			fmt.Printf("unmarshal transaction error: %v\n", err)
			return
		}
		txPayLoad := &common.TxPayload{}
		if err := proto.Unmarshal(tx.Payload, txPayLoad); err != nil {
			fmt.Printf("unmarshal tx payload error: %v\n", err)
			return
		}
		txData := &common.CommonTxData{}
		if err := proto.Unmarshal(txPayLoad.Data, txData); err != nil {
			fmt.Printf("unmarshal common tx data error: %v\n", err)
			return
		}
		utils.PrintResponse(response, nodeName, string(txData.Response.Payload))//输出结果
	} else {
		utils.PrintResponse(response, nodeName, "")//失败的结果反馈
	}
	 */
}

func createGenesis()  {
	
}

func queryChain()  {
	configPath := "C:\\Goproject\\src\\eth-relay\\go-sdk\\config\\bcs-kvhp3j-huaweichain-sdk-config.yaml"
	chainID := "huaweichain"//链名称

	nodeName := "node-3b5d2a9952882e9ed9f93db69617145a72829f59-0"

	gatewayClient,err := client.NewGatewayClient(configPath)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	//返回RawMessage 合约交易需发送的消息
	rawMsg, err := gatewayClient.ChainRawMessage.BuildQueryChainRawMessage(chainID)
	if err != nil {
		fmt.Printf("建立交易数据失败，请检查格式： %v", err)
		return
	}

	nodeMap := gatewayClient.Nodes//获取所有可用node

	node, ok := nodeMap[nodeName]//判断输入的node是否可用
	if !ok {
		fmt.Printf("不存在节点： %v\n", nodeName)
		return
	}

	invokeResponse, err := node.ChainAction.QueryChain(rawMsg)//提交交易请求 接受交易请求结果
	if err != nil {
		fmt.Printf("查询单个链配置失败： %v", err)
		return
	}

	//处理返回的配置数据
	response := &common.Response{}
	if err := proto.Unmarshal(invokeResponse.Payload, response); err != nil {
		fmt.Printf("unmarshal invoke response error： %v", err)
	}
	//检查交易是否成功
	if response.Status == common.Status_SUCCESS {
		tx := &common.Transaction{}
		if err := proto.Unmarshal(response.Payload, tx); err != nil {
			fmt.Printf("unmarshal transaction error: %v\n", err)
			return
		}
		txPayLoad := &common.TxPayload{}
		if err := proto.Unmarshal(tx.Payload, txPayLoad); err != nil {
			fmt.Printf("unmarshal tx payload error: %v\n", err)
			return
		}
		txData := &common.CommonTxData{}
		if err := proto.Unmarshal(txPayLoad.Data, txData); err != nil {
			fmt.Printf("unmarshal common tx data error: %v\n", err)
			return
		}
		utils.PrintResponse(response, nodeName, string(txData.Response.Payload))//输出结果
	} else {
		utils.PrintResponse(response, nodeName, "")//失败的结果反馈
	}
}

func queryVote()  {
	
}

func main()  {
	/*
	r := gin.Default()
	//r.StaticFS("/files", http.Dir("files"))
	r.POST("/wxdata", dealpost)
	fileGroup :=r.Group("/files",authIdentification())//权限控制
	fileGroup.StaticFS("",http.Dir("files"))
	r.POST("/securefile", dealVipFile)
	err := http.ListenAndServeTLS(":80", "../ssl/blockchaintest.club/IIS/blockchaintest.club.pem",
		"../ssl/blockchaintest.club/IIS/blockchaintest.club.key", r)
	if err != nil {
		fmt.Println(err.Error())
	}
	err=r.RunTLS(":80", "../ssl/blockchaintest.club/IIS/blockchaintest.club.pem", "../ssl/blockchaintest.club/IIS/blockchaintest.club.key")
	if err!=nil {
		fmt.Println(err.Error())
	}*/



	/*
	filepath := "C:\\Goproject\\src\\eth-relay\\go-sdk\\text.txt"
	bytees, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("Failed to read file: " + filepath)
	}
	sign,pu,err := tool.CreateSm2Sig(bytees)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("验证结果为：")
	fmt.Println(tool.VerSm2Sig(pu,bytees,sign))
	*/


}
