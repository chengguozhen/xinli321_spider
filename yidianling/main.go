// The example gets stock newses from site sina.com (http://live.sina.com.cn/zt/f/v/finance/globalnews1).
// The spider is continuous service.
// The stock api returns json result.
// It fetchs news at regular intervals that has been set in the config file.
// The result is saved in a file by PipelineFile.
package main

import (
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/usual2970/xinli321_spider/yidianling/pipeline"
    //"github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/spider"
    "github.com/PuerkitoBio/goquery"
    "github.com/hu17889/go_spider/core/common/request"
    "fmt"
    "strings"
    "strconv"
    "regexp"

)

type MyPageProcesser struct {
    startNewsId int
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package simplejson (https://github.com/bitly/go-simplejson) is used to parse data of json.
func (this *MyPageProcesser) Process(p *page.Page) {
    if !p.IsSucc() {
        println(p.Errormsg())
        return
    }
    tag:=p.GetUrlTag()
    query := p.GetHtmlParser()
    //解析列表
    if tag=="list"{
        
        var urls []string
        query.Find(".expertsList_items .item").Each(func(i int,s *goquery.Selection){

            url,_:=s.Children().Eq(0).Attr("href")
            url="http://www.yidianling.com"+url
            urls=append(urls,url)
        })
        p.AddTargetRequests(urls, "html")
        return

    }
    //解析页码
    if tag=="pages"{
        totoalPageStr:=strings.TrimSpace(query.Find("li.totle").Text())
        pageR:=regexp.MustCompile("[^\\d]+?(\\d+?)[^\\d]+")
        matches:=pageR.FindStringSubmatch(totoalPageStr)

        totalPage,_:=strconv.Atoi(matches[1]) 
        var reqs[]*request.Request
        for i:=1;i<=totalPage;i++{
            url:="http://www.yidianling.com/experts?page="+strconv.Itoa(i)
            req := request.NewRequest(url, "html", "list", "GET", "", nil, nil, nil, nil)
            reqs=append(reqs,req)
        }
        p.AddTargetRequestsWithParams(reqs)
        return

    }

    //解析详情页
    name:=strings.TrimSpace(query.Find("div.e-info .i-right h1.f-24").Text())

    titleR:=regexp.MustCompile("([^\\s]+?)\\s+?([^\\s]+)")
    city:=strings.TrimSpace(query.Find("div.e-info .i-right p.mt-15").Text())
    matches:=titleR.FindStringSubmatch(city)

    title:=matches[1]
    city=strings.TrimSpace(matches[2])

    brief,_:=query.Find("div.e-summary .ctt p").Html()
    brief=strings.TrimSpace(brief)

    desc,_:=query.Find("div.e-summary .ctt .card-line .mt-5").Html()
    desc=strings.TrimSpace(desc)

    
    good_at,_:=query.Find("div.e-info .i-right .desc td").Eq(1).Find("span").Html()
    good_at=strings.TrimSpace(good_at)

    // //头像
    head_url,_:=query.Find("div.e-info .i-left img").Attr("src")



    // article_num:=strings.TrimSpace(query.Find(".jg-view ul h4").Eq(0).Text())
    answer_num:=strings.TrimSpace(query.Find("div.e-answers .title span.txt").Text())
    answerR:=regexp.MustCompile("(\\d+?)[^\\d]");
    answerMatcheds:=answerR.FindStringSubmatch(answer_num)
    answer_num=answerMatcheds[1]

    praise_num:=strings.TrimSpace(query.Find("div.e-info .i-right .data .item").Eq(1).Find(".num").Text())

    consult_num:=strings.TrimSpace(query.Find("div.e-info .i-right .data .item").Eq(2).Find(".num").Text())

    origin_url:=p.GetRequest().GetUrl()


    r := regexp.MustCompile("\\/(\\d+)")
    idMatches:=r.FindStringSubmatch(origin_url)
    origin_id:=idMatches[1]

    tempId,_:=strconv.Atoi(origin_id)
    hash:=strconv.Itoa(2000000000000+tempId)


    // //计算price
    priceR:=regexp.MustCompile("(\\d+?)元")
    price:=1000000
    query.Find(".expert-right .consult-box p .cl-orange").Each(func(i int,s *goquery.Selection){
        priceStr:=strings.TrimSpace(s.Text())
        matches:=priceR.FindStringSubmatch(priceStr)
        tempPrice,_:=strconv.Atoi(matches[1])
        if tempPrice < price{
            price=tempPrice
        }

    })
    p.AddField("price",strconv.Itoa(price))

    p.AddField("head_url",head_url)

    p.AddField("name",name)

    p.AddField("city",city)

    p.AddField("title",title)

    p.AddField("brief",brief)

    p.AddField("desc",desc)

    p.AddField("good_at",good_at)

    p.AddField("article_num","0")

    p.AddField("answer_num",answer_num)

    p.AddField("praise_num",praise_num)

    p.AddField("consult_num",consult_num)

    p.AddField("site_id","2")

    p.AddField("origin_url",origin_url)

    p.AddField("origin_id",origin_id)

    p.AddField("hash",hash)
    
}

func (this *MyPageProcesser) Finish() {
    fmt.Printf("TODO:before end spider \r\n")
}

func main() {
    // spider input:
    //  PageProcesser ;
    //  task name used in Pipeline for record;
    //  
    //  
    req := request.NewRequest("http://www.yidianling.com/experts", "html", "pages", "GET", "", nil, nil, nil, nil)
    spider.NewSpider(NewMyPageProcesser(), "yidianling").
        AddRequest(req). // start url, html is the responce type ("html" or "json" or "jsonp" or "text")
        AddPipeline(pipeline.NewPipelineMysql()).                                                                                   // Print result to std output                                                                      // Print result in file
        OpenFileLog("/tmp"). 
        SetThreadnum(5).                                                                                                         // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".                                                                                        // Sleep time between 1s and 3s.
        Run()

    //AddPipeline(pipeline.NewPipelineFile("/tmp/tmpfile")). // print result in file
}
