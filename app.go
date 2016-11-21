package main

import (
	"encoding/binary"
	"fmt"

	"log"

	"github.com/baabeetaa/glogchain/blog"
	"github.com/baabeetaa/glogchain/protocol"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/types"
)

//structure for holding golgchain app types
type GlogChainApp struct {
	hashCount int
	txCount   int
}

func NewGlogChainApp() *GlogChainApp {
	return &GlogChainApp{}
}

func (app *GlogChainApp) Info() string {
	return Fmt("hashes:%v, txs:%v", app.hashCount, app.txCount)
}

func (app *GlogChainApp) SetOption(key string, value string) (log string) {
	return ""
}

func (app *GlogChainApp) AppendTx(tx []byte) types.Result {
	// tx is json string, need to convert to text and then parse into json object
	// yes, and the conversions are crazy-annoying in go.  Am I right?
	jsonstring := string(tx[:])

	obj, err := protocol.UnMarshal(jsonstring)

	if err != nil {
		log.Fatal(err)
		return types.ErrEncodingError
	}

	switch v := obj.(type) {
	case protocol.PostOperation:
		var objPostOperation protocol.PostOperation

		objPostOperation = v
		//fmt.Println("Title=" + objPostOperation.Title)
		//fmt.Println("Body=" + objPostOperation.Body)
		//fmt.Println("Author=" + objPostOperation.Author)

		blog.CreatePost(&objPostOperation)

	default:
	}

	return types.OK
}

func (app *GlogChainApp) CheckTx(tx []byte) types.Result {
	return types.OK
}

func (app *GlogChainApp) Commit() types.Result {
	app.hashCount += 1

	if app.txCount == 0 {
		return types.OK
	} else {
		hash := make([]byte, 8)
		binary.BigEndian.PutUint64(hash, uint64(app.txCount))
		return types.NewResultOK(hash, "")
	}
}

//Query tells the user that the query they are requesting isn't supported.
func (app *GlogChainApp) Query(query []byte) types.Result {
	return types.NewResultOK(nil, fmt.Sprintf("Query is not supported"))
}
