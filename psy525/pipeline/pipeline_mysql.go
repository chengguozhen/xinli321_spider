package pipeline

import (
    "github.com/hu17889/go_spider/core/common/com_interfaces"
    "github.com/hu17889/go_spider/core/common/page_items"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
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


func (this *PipelineMysql) insert(item map[string]string) {


	stmt, err := this.pMysql.Prepare("INSERT INTO experts(name, city,title,brief,`desc`,good_at,article_num,answer_num,praise_num,consult_num,site_id,origin_url,origin_id,hash,ftf_price,head_url) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
 
	if err != nil {
		log.Println(err)
		return
	}

	stmt.Exec(item["name"],item["city"],item["title"],item["brief"],item["desc"],item["good_at"],item["article_num"],item["answer_num"],item["praise_num"],item["consult_num"],item["site_id"],item["origin_url"],item["origin_id"],item["hash"],item["price"],item["head_url"])
	
 
}