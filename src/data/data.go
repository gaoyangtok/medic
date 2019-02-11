package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"
)

type Sex string

type Foo struct {
	Name		string
	Phone		string
	Create		time.Time
	Update		time.Time
	Diagnosed	string
	Program		string
	AllFee		float64
	RealFee		float64
	PaidFee		float64
	Address		string
	Age			int
	Sex			Sex
	Index		int
	Checked 	bool
	Deleted		bool
}
const data="data.csv"

//Write 写入运行配置
func Write(dabs []*Foo) {
	//_ = os.Remove(data)
	f, err := os.OpenFile(data, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		fmt.Print("写文件失败！")
	}
	write:=csv.NewWriter(f)
	records:=make([][]string, len(dabs)+1)
	records[0]=[]string{"姓名","电话","登记时间","最新时间","病例诊断","治疗方案","就诊费用","实收费用","已付费用","住址","性别","年龄","是否删除"}
	for index,foo:=range dabs {
		var del string
		if foo.Deleted {
			del = "1";
		}else {
			del = "0"
		}
		records[index+1]=[]string{
			foo.Name,
			foo.Phone,
			foo.Create.Format("2006-01-02 15:04:05"),
			foo.Update.Format("2006-01-02 15:04:05"),
			foo.Diagnosed,
			foo.Program,
			strconv.FormatFloat(foo.AllFee, 'f', 1, 64),
			strconv.FormatFloat(foo.RealFee, 'f', 1, 64),
			strconv.FormatFloat(foo.PaidFee, 'f', 1, 64),
			foo.Address,string(foo.Sex),strconv.Itoa(foo.Age),
			del,
		}
	}
	_ = write.WriteAll(records)

	f.Close()
	return
}

//Read 读取运行配置
func Read() []*Foo {
	file, err := os.Open(data)
	if err != nil {
		Write([]*Foo{}) //如果不存在先创建一个空文件
		file,_ = os.Open(data)
	}
	defer file.Close()

	read := csv.NewReader(file)

	records ,err:=read.ReadAll()
	dabs := make([]*Foo, len(records)-1)
	for index:= range dabs{
		record:=records[index+1]
		allFee,_ := strconv.ParseFloat(record[6],64)
		realFee,_ := strconv.ParseFloat(record[7],64)
		paidFee,_ := strconv.ParseFloat(record[8],64)
		create, _ :=time.Parse("2006-01-02 15:04:05",record[2])
		update, _ :=time.Parse("2006-01-02 15:04:05",record[3])
		age, err := strconv.Atoi(record[11])
		if err!=nil{
			panic(err)
		}
		dabs[index] = &Foo{
			Name: record[0],
			Phone: record[1],
			Create: create,
			Update: update,
			Diagnosed: record[4],
			Program: record[5],
			AllFee: allFee,
			RealFee: realFee,
			PaidFee: paidFee,
			Address: record[9],
			Sex: Sex(record[10]),
			Age:age,
			Index: index,
			Deleted:len(record)>=13 && strings.Compare(record[12],"1")==0,
		}
	}

	file.Close()
	return dabs
}
//
//func main() {
//	fobs :=Read()
//	for index,_ := range fobs {
//		foo:= fobs[index]
//		fmt.Println(foo.Name,foo.Phone,foo.Create,foo.Update,foo.AllFee,foo.RealFee,foo.PaidFee,foo.Diagnosed,foo.Program,foo.Addr)
//	}
//	Write([]*Foo{{Name: "高阳", Phone: "18611754986",Create:"2018-01-01 12:00:00",Update:"2018-01-01 12:00:00",AllFee:"200",RealFee:"150",PaidFee:"50"}})
//	fobs =Read()
//	for index,_ := range fobs {
//		foo:= fobs[index]
//		fmt.Println(foo.Name,foo.Phone,foo.Create,foo.Update,foo.AllFee,foo.RealFee,foo.PaidFee,foo.Diagnosed,foo.Program,foo.Addr)
//	}
//}