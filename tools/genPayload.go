package tools

import (
	"fmt"
	"github.com/beevik/etree"
	"math/rand"
	"strings"
)

func genRandom4Num(source *rand.Rand) string {
	return fmt.Sprintf("%04v", source.Int31n(10000))
}

func genRandomStr(length int, lowercase bool, alphabet string) string{
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	result := string(b)
	if lowercase{
		return strings.ToLower(result)
	}
	return result
}

func genHeuristicCheckPayload() string {
	payload := genRandomStr(10,false,HEURISTIC_CHECK_ALPHABET)
	if strings.Count(payload,"\"") !=1 || strings.Count(payload,"'") !=1 {
		payload = genHeuristicCheckPayload()
	}
	return payload
}

// 解析xml的话 应该只执行一次 xml --> map
func genDbmsErrorFromXml() map[string]string {
	dbmsErrorKeywordList := make(map[string]string)
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("xml/errors.xml"); err != nil {
		panic(err)
	}
	root := doc.SelectElement("root")
	for _, dbms := range root.SelectElements("dbms") {
		for _, dbName := range dbms.Attr {
			for _ , e := range dbms.SelectElements("error"){
				for _ , errWord := range e.Attr{
					//log.Println(dbName.Value,errWord.Value)
					dbmsErrorKeywordList[errWord.Value] = dbName.Value
				}
			}
		}
	}
	return dbmsErrorKeywordList
}
