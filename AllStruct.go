package storeKey

//过往认证记录
type CertificateHistory struct {
	FileName string
	Hash   string
}

//过往购买记录
type PurchaseHistory struct {
	FileName string
	TranAdd   string
}

//两者历史记录的json数组
type History struct {
	CertificateHistorys []CertificateHistory
	PurchaseHistorys []PurchaseHistory
}

//商品tag
type Tags struct {
	Tag []string
}

type Wings struct{ //用于标识此学生学习了哪个高校的课程
	//学前学后fid和学号等会出现变化
	//结构体hash得到账本id
	//hash得到fid
	Uid			string//学前标识符
	Id 			string//学号-待定是什么学号？大学的还是其他的？
	Name 		string //姓名
	University	string//学校
	School		string//学院
	Course		string//课程
	Time		string//开始课程时间
	Evidence	string//记录
}

type Args struct {
	A []string
}

type CopyRight struct {
	Hash string//版权哈希
	Owner string//拥有者
	Filename string//文件名
	FileID string//云储存文件下载标识
	TxId string//交易ID
	Time int64//认证时间
}

type Purchase struct {
	Buyer string//购买者
	Owner string//拥有者
	Hash string//版权哈希
	Price float64//金额
	//Id int//单号
	TxId string//交易ID
}