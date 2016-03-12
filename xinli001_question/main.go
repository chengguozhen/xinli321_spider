// The example gets stock newses from site sina.com (http://live.sina.com.cn/zt/f/v/finance/globalnews1).
// The spider is continuous service.
// The stock api returns json result.
// It fetchs news at regular intervals that has been set in the config file.
// The result is saved in a file by PipelineFile.
package main

import (
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/usual2970/xinli321_spider/xinli001_question/pipeline"
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
    if tag=="list"{
        
        var urls []string
        query.Find(".ask_lists .items .ask_show").Each(func(i int,s *goquery.Selection){

            url,_:=s.Find("a").Attr("href")
            urls=append(urls,url)
        })
        p.AddTargetRequests(urls, "html")
        return

    }

    title:=strings.TrimSpace(query.Find(".infos-wrap .show_ask h2").Text())
    
    origin_url:=p.GetRequest().GetUrl()


    r := regexp.MustCompile("\\/(\\d+)")
    matches:=r.FindStringSubmatch(origin_url)
    origin_id:=matches[1]

    origin_idInt,_:=strconv.Atoi(origin_id)
    hash:=strconv.Itoa(1000000000000+origin_idInt)

    content:=strings.TrimSpace(query.Find(".infos-wrap .descs").Text())
    contentR:=regexp.MustCompile("[^\\n]+?\\n(.+)")
    contentMatches:=contentR.FindStringSubmatch(content)
    content=strings.TrimSpace(contentMatches[1])
    tags:=strings.TrimSpace(query.Find(".infos-wrap .descs span a").Text())

    answer_num:=query.Find(".answer_list").Children().Length()

    
    p.AddField("title",title)
    p.AddField("content",content)
    p.AddField("origin_url",origin_url)
    p.AddField("origin_id",origin_id)
    p.AddField("hash",hash)
    p.AddField("tags",tags)
    p.AddField("answer_num",strconv.Itoa(answer_num))

    
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

        url:="http://qa.xinli001.com/done/p"+strconv.Itoa(i)
        req := request.NewRequest(url, "html", "list", "GET", "", nil, nil, nil, nil)
        reqs=append(reqs,req)
    }
    spider.NewSpider(NewMyPageProcesser(), "xinli001").
        AddRequests(reqs). // start url, html is the responce type ("html" or "json" or "jsonp" or "text")
        AddPipeline(pipeline.NewPipelineMysql()).                                                                                   // Print result to std output                                                                      // Print result in file
        OpenFileLog("/tmp"). 
        SetThreadnum(5).                                                                                                         // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".                                                                                        // Sleep time between 1s and 3s.
        Run()

}
