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

/**
 * 下载图片，并将数据插入到数据库
 * @param  {[type]} this *PipelineMysql) insert(item map[string]string [description]
 * @return {[type]}      [description]
 */
func (this *PipelineMysql) insert(item map[string]string) {



    //查看之前的数据，并将数据存入快照
  //   `id` char(20) NOT NULL DEFAULT '',
  // `origin_id` char(20) DEFAULT NULL,
  // `origin_url` varchar(125) DEFAULT NULL,
  // `title` varchar(255) DEFAULT NULL,
  // `content` text,
  // `answer_num` int(11) DEFAULT NULL,
  // `tags` varchar(255) DEFAULT NULL,
	stmt, err := this.pMysql.Prepare("REPLACE INTO questions(id,origin_id,origin_url,title,content,answer_num,tags) VALUES(?,?,?,?,?,?,?)")
	defer stmt.Close()
 
	if err != nil {
		log.Println(err)
		return
	}

	stmt.Exec(item["hash"],item["origin_id"],item["origin_url"],item["title"],item["content"],item["answer_num"],item["tags"])
	
 
}

