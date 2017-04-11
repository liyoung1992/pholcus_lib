package pholcus_lib

// 基础包
import (
	"fmt"
	"strconv"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
)

func init() {
	//添加自身到网易菜单
	Blog.Register()
}

var Blog = &Spider{
	Name:        "上海链家",
	Description: "上海链家 【http://sh.lianjia.com/】",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "http://sh.lianjia.com/ershoufang", Rule: "上海链家二手房主页"})
		},

		Trunk: map[string]*Rule{

			"上海链家二手房主页": {
				ParseFunc: func(ctx *Context) {
					i := 1
					for i <= 100 {
						url := "http://sh.lianjia.com/ershoufang/d" + strconv.Itoa(i)
						ctx.AddQueue(&request.Request{Url: url, Rule: "链接二手房"})
						i++
					}
				},
			},

			"链接二手房": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".list-wrap ul li .info-panel h2 a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							str := "http://sh.lianjia.com"
							str = str + (url)
							fmt.Println(str)

							ctx.AddQueue(&request.Request{Url: str, Rule: "链接二手房详细"})
						}
					})

				},
			},
			"链接二手房详细": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"title",      //标题
					"totalPrice", //总价
					"area",       //面积
					"houseType",  //户型

					"address",      //地址
					"areaEllipsis", //区域

				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					downPayment := ""
					monthPayment := ""
					address := ""
					areaEllipsis := ""
					//获取标题
					title := query.Find(".title h1").Text()
					//获取总价
					totalPrice := query.Find(".price .mainInfo").Text()

					//获取面积
					area := query.Find(".area .mainInfo").Text()

					//户型
					houseType := query.Find(".room .mainInfo").Text()

					//首付

					query.Find(".aroundInfo tr").Each(func(i int, s *goquery.Selection) {

						if i == 4 {
							areaEllipsis = s.Find("td .f1 .areaEllipsis").Text()
						} else if i == 5 {
							address = s.Find("td .addrEllipsis").Text()
						}
					})

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: title,
						1: totalPrice,
						2: area,
						3: houseType,

						4: address,
						5: areaEllipsis,
					})
				},
			},
		},
	},
}
