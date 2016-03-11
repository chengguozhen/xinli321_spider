// The example gets stock newses from site sina.com (http://live.sina.com.cn/zt/f/v/finance/globalnews1).
// The spider is continuous service.
// The stock api returns json result.
// It fetchs news at regular intervals that has been set in the config file.
// The result is saved in a file by PipelineFile.
package main

import (
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/usual2970/xinli321_spider/xinli001/pipeline"
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
    if tag=="list"{
        query := p.GetHtmlParser()
        var urls []string
        query.Find("li").Each(func(i int,s *goquery.Selection){

            url,_:=s.Find(".img a").Attr("href")
            urls=append(urls,url)
        })
        p.AddTargetRequests(urls, "html")
        return

    }

    query := p.GetHtmlParser()

    name:=strings.TrimSpace(query.Find("div.desc-edit .fs16").Text())

    city:=strings.TrimSpace(query.Find("div.desc-edit .city-edit .content").Text())

    title:=strings.TrimSpace(query.Find(".introduce-edit .content").Text())

    brief:=strings.TrimSpace(query.Find(".brief-edit .content").Text())

    desc,_:=query.Find(".jg-jj .desc").Html()
    desc=strings.TrimSpace(desc)

    
    good_at,_:=query.Find(".jg-desc p").Html()
    good_at=strings.TrimSpace(good_at)

    //头像
    head_url,_:=query.Find(".jg-view .img img").Attr("src")



    article_num:=strings.TrimSpace(query.Find(".jg-view ul h4").Eq(0).Text())
    answer_num:=strings.TrimSpace(query.Find(".jg-view ul h4").Eq(1).Text())

    praise_num:=strings.TrimSpace(query.Find(".jg-view ul h4").Eq(2).Text())

    consult_num:=strings.TrimSpace(query.Find(".jg-view ul h4").Eq(3).Text())

    origin_url:=p.GetRequest().GetUrl()


    r := regexp.MustCompile("\\/(\\d+)")
    matches:=r.FindStringSubmatch(origin_url)
    origin_id:=matches[1]
    origin_idInt,_:=strconv.Atoi(origin_id)
    hash:=strconv.Itoa(1000000000000+origin_idInt)


    //计算price
    priceR:=regexp.MustCompile("(\\d+?)元")
    price:=1000000
    query.Find(".jg-desc dl span").Each(func(i int,s *goquery.Selection){
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

    p.AddField("article_num",article_num)

    p.AddField("answer_num",answer_num)

    p.AddField("praise_num",praise_num)

    p.AddField("consult_num",consult_num)

    p.AddField("site_id","1")

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
    var reqs []*request.Request
    for i:=1;i<=130;i++ {

        url:="http://www.xinli001.com/ajax/teacher-list.json?slug=&tag=&flag=date&city=&name=&page="+strconv.Itoa(i)
        req := request.NewRequest(url, "html", "list", "GET", "", nil, nil, nil, nil)
        reqs=append(reqs,req)
    }
    spider.NewSpider(NewMyPageProcesser(), "xinli001").
        AddRequests(reqs). // start url, html is the responce type ("html" or "json" or "jsonp" or "text")
        AddPipeline(pipeline.NewPipelineMysql()).                                                                                   // Print result to std output                                                                      // Print result in file
        OpenFileLog("/tmp"). 
        SetThreadnum(5).                                                                                                         // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".                                                                                        // Sleep time between 1s and 3s.
        Run()

    //AddPipeline(pipeline.NewPipelineFile("/tmp/tmpfile")). // print result in file
}
