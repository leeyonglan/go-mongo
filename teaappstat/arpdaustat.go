package teaappstat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	adreportPath   = "ad_order.xlsx"
	androidPayPath = "androidPay.xlsx"
	iapconfigPath  = "gameconfig/IAP_package.json"
	iapconfg       map[string][]packItem
)

type userOrderPay struct {
	CreateTime    int                    `bson:"create_time"`
	GoodsId       string                 `bson:"goods_id"`
	OrderType     int8                   `bson:"order_type"`
	GoodInfo      map[string]interface{} `bson:"goodsInfo"`
	DayTotalPrice float32
}

type packItem struct {
	AndroidId string  `json:"android_id"`
	Price     float32 `json:"price_dollar"`
}

func Do() {
	Init()
	err := parseIAPConfig(iapconfigPath)
	if err != nil {
		return
	}

	var iOS plateInter = IOS{
		Name: "iOS",
	}
	ArpDauStat(iOS)

	var android plateInter = Android{
		IOS: IOS{
			Name: "Android",
		},
	}
	ArpDauStat(android)
}

func ArpDauStat(plat plateInter) {
	err, adList := plat.getAdList()
	if err != nil {
		fmt.Printf("get ad report list err %v", err)
	}

	var adTotal float32
	var iapTotal float32
	iapDict := plat.getPayList()
	for i, v := range adList {
		day := v[0]

		adrevenue, _ := strconv.ParseFloat(adList[i][4], 32)
		var totalrevenue float32 = float32(adrevenue)
		adTotal += totalrevenue
		var iap float32

		var totalrevenueStr string = adList[i][4]
		//判断是否有充值数据
		if pay, ok := iapDict[day]; ok {
			iap = pay.DayTotalPrice
			iapTotal += iap
			totalrevenue += pay.DayTotalPrice
			totalrevenueStr = strconv.FormatFloat(float64(totalrevenue), 'f', -1, 32)
			// adList[i][4] = totalrevenueStr
		}
		// arpdau
		dau, err := strconv.Atoi(adList[i][3])
		if err != nil {
			fmt.Printf("convert dau type err %v", err)
		}

		arpdau := totalrevenue / float32(dau)
		arpdaustr := strconv.FormatFloat(float64(arpdau), 'f', -1, 32)
		adList[i][5] = arpdaustr

		iapstr := strconv.FormatFloat(float64(iap), 'f', -1, 32)
		adList[i] = append(adList[i], iapstr)
		adList[i] = append(adList[i], totalrevenueStr)
	}

	fmt.Println(plat.getPlatName()+"_adTotal:", adTotal)
	fmt.Println(plat.getPlatName()+"_iapTotal:", iapTotal)

	// outFileName := plat.getPlatName() + "_arpdau.xlsx"
	// outputXlsx(outFileName, "arp", []string{"Day", "Application", "Platform", "DAU", "AD_Revenue", "ARPDAU", "IAP", "Total_Revenue"}, adList)
}

func parseIAPConfig(iapconfpath string) (err error) {
	data, err := ioutil.ReadFile(iapconfpath)

	if err != nil {
		fmt.Println("Failed to read JSON file:", err)
		return
	}
	err = json.Unmarshal(data, &iapconfg)
	if err != nil {
		fmt.Println("Failed to parse JSON file:", err)
	}
	return
}

func getpriceById(id string) float32 {
	var price float32
	for _, v := range iapconfg["diamond_pack"] {
		if v.AndroidId == id {
			price = v.Price
			break
		}
	}
	return price
}

func formatDate(date string) string {
	dateArr := strings.Split(date, "-")
	return "20" + dateArr[2] + "/" + dateArr[0] + "/" + dateArr[1]
}
