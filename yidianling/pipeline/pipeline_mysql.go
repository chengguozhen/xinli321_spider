package pipeline

import (
    "github.com/hu17889/go_spider/core/common/com_interfaces"
    "github.com/hu17889/go_spider/core/common/page_items"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"bytes"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "github.com/nfnt/resize"
    "image/png"
    "image/jpeg"
    "regexp"
    "image"
    "strconv"
    "time"
    "github.com/usual2970/xinli321/config"
)

type PipelineMysql struct {
	pMysql *sql.DB

}

func NewPipelineMysql() *PipelineMysql {
	db, err := sql.Open("mysql", "qdm194418114:123456Aa@tcp(qdm194418114.my3w.com:3306)/qdm194418114_db?allowOldPasswords=1")
	if err != nil {
		panic("Open database error")
	}
    return &PipelineMysql{pMysql:db}
}

func (this *PipelineMysql) Process(items *page_items.PageItems, t com_interfaces.Task) {
    println("----------------------------------------------------------------------------------------------")
    println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
    if tag:=items.GetRequest().GetUrlTag();tag=="list"{
    	return
    }
    this.insert(items.GetAll())
}

/**
 * 下载图片，并将数据插入到数据库
 * @param  {[type]} this *PipelineMysql) insert(item map[string]string [description]
 * @return {[type]}      [description]
 */
func (this *PipelineMysql) insert(item map[string]string) {

	imgDic:=getImg(item["head_url"],item["hash"])

    resizeImg(imgDic,item["hash"])


    //查看之前的数据，并将数据存入快照
    var oarticle_num,oanswer_num,opraise_num,oconsult_num,ioftf_price int
    var oftf_price string
    oerr:=this.pMysql.QueryRow("SELECT article_num,answer_num,praise_num,consult_num,ftf_price FROM experts where id=?",item["hash"]).Scan(&oarticle_num,&oanswer_num,&opraise_num,&oconsult_num,&oftf_price)

    var article_num,answer_num,praise_num,consult_num =0,0,0,0
    var carticle_num,canswer_num,cpraise_num,cconsult_num =0,0,0,0
    var ftf_price,cftf_price int
    if oerr == nil{
       carticle_num,_=strconv.Atoi(item["article_num"])
       canswer_num,_=strconv.Atoi(item["answer_num"])
       cpraise_num,_=strconv.Atoi(item["praise_num"])
       cconsult_num,_=strconv.Atoi(item["consult_num"])
       cftf_price,_=strconv.Atoi(item["price"])
        temp:=strings.Split(oftf_price,".")
        ioftf_price,_=strconv.Atoi(temp[0])
       article_num = carticle_num - oarticle_num
       answer_num = canswer_num - oanswer_num
       praise_num = cpraise_num - opraise_num
       consult_num = cconsult_num - oconsult_num
       println("diff:",article_num,answer_num,praise_num,consult_num)
       ftf_price = cftf_price -ioftf_price

    }

    snap_at:=time.Now().Format("2006-01-02")

    snap_stmt,snap_err:= this.pMysql.Prepare("insert into snapshots(experts_id,ftf_price,article_num,praise_num,consult_num,answer_num,snap_at) VALUES (?,?,?,?,?,?,?)");
    defer snap_stmt.Close()
    if(snap_err!=nil){
        log.Println(snap_err)
        return
    }
    snap_stmt.Exec(item["hash"],ftf_price,article_num,praise_num,consult_num,answer_num,snap_at)

	stmt, err := this.pMysql.Prepare("REPLACE INTO experts(name, city,title,brief,`desc`,good_at,article_num,answer_num,praise_num,consult_num,site_id,origin_url,origin_id,id,ftf_price,head_url) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
 
	if err != nil {
		log.Println(err)
		return
	}

	stmt.Exec(item["name"],item["city"],item["title"],item["brief"],item["desc"],item["good_at"],item["article_num"],item["answer_num"],item["praise_num"],item["consult_num"],item["site_id"],item["origin_url"],item["origin_id"],item["hash"],item["price"],item["head_url"])
	
 
}

func getImg(url,name string) string{
	
    path := strings.Split(url, ".")

    tempData:=strings.Split(url,"!")
    newUrl:=tempData[0]+"!s800x800"

    picFile:=config.ImagePath+"experts/2/origin/"+name+"."+path[len(path)-1]
    out, _ := os.Create(picFile)
    defer out.Close()
    client := &http.Client{
		CheckRedirect: nil,
	}


	req, _ := http.NewRequest("GET", newUrl, nil)
	req.Header.Add("Referer","http://www.yidianling.com/experts/")
	req.Header.Add("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	resp, _ := client.Do(req)

    defer resp.Body.Close()
    pix, _ := ioutil.ReadAll(resp.Body)
    io.Copy(out, bytes.NewReader(pix))
    return config.ImagePath+"experts/2/origin/"+name+"."+path[len(path)-1]

}

func resizeImg(imgDic,name string){
    file,_:= os.Open(imgDic)
    defer file.Close()
    pngR:=regexp.MustCompile("\\.png")

    var img image.Image

    if pngR.MatchString(imgDic){
        img, _ = png.Decode(file)

    }else{
        img, _ = jpeg.Decode(file)
     
    }

    file160,_:=os.Create(config.ImagePath+"experts/2/160/"+name+".jpg")
    defer file160.Close()
    file80,_:=os.Create(config.ImagePath+"experts/2/80/"+name+".jpg")
    defer file80.Close()

    dstImage160 := resize.Resize(160, 160, img, resize.Lanczos3)

    dstImage80 := resize.Resize(80, 80, img, resize.Lanczos3)

    


    jpeg.Encode(file160, dstImage160, nil)

    jpeg.Encode(file80, dstImage80, nil)

}