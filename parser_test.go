package dbcparser_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	p "github.com/rhabichl/parser"
)


func TestNew(t *testing.T){
    if p.New() == nil {
        t.Fatal("New() returns nil")
    }
}

func TestLoad(t *testing.T){
    
    testData := []struct{
        id int
        input []byte
    }{}

    f, _ :=os.ReadFile("./can_messages.log")
    a := strings.Split(string(f), "\n")
    for _, aa := range a {
        if len(aa) >0 {
            bb := strings.Split(aa, "#")
            i,_ := strconv.Atoi(bb[0])
            in, _ := hex.DecodeString(bb[1])
            testData = append(testData, struct{
                id int 
                input []byte
            }{
                id: i, 
                input: in,
            })
        }
    }
    parser, err := p.New().Load("./test/data/4.dbc")
    if err != nil {
        t.Fatal(err)
    }
    for _, td := range testData {
        t.Run("", func(t *testing.T) {
            msg, err := parser.Parse(td.id, td.input)
            if err != nil {
                t.Errorf("error: %v", err)
            }
            p, _ :=json.MarshalIndent(msg, "", "  ")
            fmt.Println(string(p))
        })
    }
}
