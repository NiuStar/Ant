package Ant

import (
	"fmt"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"nqc.cn/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"encoding/json"
	"io/ioutil"
	"nqc.cn/log"
)


/*
{
    "url": "http://company.zhaopin.com/zhengzhou/210500/",
    "ant": [
        {
            "parend": {
                "class": "fleft checkjobs width280"
            },
            "grandpa": {
                "class": "jobs-list-box"
            },
            "curr": {},
            "currAtom": "A",
            "attr": [
                "href"
            ],
            "utf8": true
        }
    ]
}


*/

func GoAnt(c *gin.Context) {

	ant, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Write(err)
	}
	c.Request.Body.Close()

	j3 := make(map[string]interface{})
	j2 := make(map[string]interface{})
	err = json.Unmarshal(ant,&j3)
	if err != nil {
		fmt.Println("start ANT false2")
		c.String(http.StatusOK,`{"state":false,"msg":"抓取失败"}`)
		return
	}
	url := j3["url"].(string)

	var antList []AntInfo

	var list []interface{} = j3["ant"].([]interface{})

	for _,value_s := range list {
		fmt.Println("a")
		//j1 := make(map[string]interface{})
		j1 := value_s.(map[string]interface{})
		fmt.Println("123456")

		var p map[string]string = make(map[string]string)
		var pp map[string]string = make(map[string]string)
		var curr map[string]string = make(map[string]string)
		var currAtom string
		var attr []string
		var utf8 bool
		if j1["grandpa"] != nil {
			fmt.Println("grandpa")
			pp1 := j1["grandpa"].(map[string]interface{})

			for key,value := range pp1 {

				pp[key] = value.(string)
			}


		}

		if j1["parend"] != nil {
			pp1 := j1["parend"].(map[string]interface{})

			for key,value := range pp1 {

				p[key] = value.(string)
			}
		}
		if j1["curr"] != nil {
			pp1 := j1["curr"].(map[string]interface{})

			for key,value := range pp1 {

				curr[key] = value.(string)
			}
		}
		if j1["currAtom"] != nil {
			currAtom = j1["currAtom"].(string)
		}

		if j1["attr"] != nil {
			//fmt.Println("attr",j1["attr"])
			pp1 := j1["attr"].([]interface{})
			//fmt.Println("attr  1")
			for key,value := range pp1 {
				fmt.Println("key = ",key,"value = ",value)
				attr = append(attr,value.(string))
			}
			//attr = j1["attr"].([]string)
			//fmt.Println("attr end")
		}
		if j1["utf8"] != nil {
			utf8 = j1["utf8"].(bool)
		}
		a := AntInfo{p:p,pp:pp,curr:curr,currAtom:currAtom,attr:attr,utf8:utf8}
		antList = append(antList,a)

	}

	//fmt.Println("start ANT")

	result := GetInfo(url,antList)
	j2["state"] = true
	j2["msg"] = "抓取成功"
	j2["result"] = result

	body , err := json.Marshal(j2)

	if err != nil {
		c.String(http.StatusOK,`{"state":false,"msg":"抓取失败"}`)
		return
	}
	c.String(http.StatusOK,string(body))
	return

}

type AntInfo struct {
	p map[string]string
	pp map[string]string
	curr map[string]string
	currAtom string //当前节点的标签
	attr []string //获取更多结果值的类型 是text , href 或者其它
	utf8 bool
}


func GetInfo(url string,infoList []AntInfo) []interface{} {
	// request and parse the front page
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	// define a matcher

	var results []interface{}

	for _,info := range infoList {
		matcher := getMatcher(info)
		articles := scrape.FindAll(root, matcher)
		var result []interface{}

		for _, article := range articles {
			j1 := make(map[string]interface{})
			if !info.utf8 {
				j1["text"] = utils.ConvertGBK(scrape.Text(article))
			} else {
				j1["text"] = scrape.Text(article)
			}
			for _,value := range info.attr {
				j1[value] = scrape.Attr(article,value)
			}

			//fmt.Println(scrape.Text(article))

			result = append(result,j1)
		}
		results = append(results,result)
	}

	return results
}

func getMatcher(info AntInfo) scrape.Matcher {
	matcher := func(n *html.Node) bool {

		//return true
		// must check for nil values
		if len(info.currAtom) > 0 {
			if n.DataAtom != utils.GetAtomByString(info.currAtom) {
				return false
			}
		}
		if len(info.pp) > 0 {
			if n.Parent == nil || n.Parent.Parent == nil {
				return false
			}


		} else if len(info.p) > 0 {
			if n.Parent == nil {
				return false
			}

		}
		for key,value := range info.pp {
			if scrape.Attr(n.Parent.Parent, key) != value {
				return false
			}
		}
		for key,value := range info.p {
			if scrape.Attr(n.Parent, key) != value {
				return false
			}
		}
		for key,value := range info.curr {
			if scrape.Attr(n, key) != value {
				return false
			}
		}

		return true
	}
	return matcher
}