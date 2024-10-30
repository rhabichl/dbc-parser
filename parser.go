package dbcparser

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CanMsg struct {
    Id int
    Name string
    Length uint8
    Sender string
    Signals []CanSignal
}

type CanSignal struct {
    Name string
    BitStart uint8
    Length uint8
    LittleEndian bool
    Signed bool
    Scale float64
    Offset float64
    Min int
    Max int
    Unit string
    Reciever string
    Value []byte
}

type Parser interface{
    Parse(int, []byte) (map[string]int, error)
    Load(string) (Parser, error)
}

type ProtoParser struct{
    Msgs []CanMsg
}

func (p ProtoParser) Parse(id int ,frame []byte) (map[string]int, error){
    result := make(map[string]int)
    for _, i := range p.Msgs {
        if i.Id == id  {    
            for _,s := range i.Signals {
                // if length is 8 bits
                if s.Length == 8 {
                    result[s.Name] = int(frame[s.BitStart / 8])
                }else {
                    first := s.BitStart/8
                    last := (s.BitStart/8 + s.Length/8 )
                    fmt.Printf("first: %v last:%v lenghtOfFrame:%v\n", first, last, len(frame))
                    if s.LittleEndian {
                        result[s.Name] = int(binary.LittleEndian.Uint16(frame[first:last]))
                    } else {
                       result[s.Name] = int(binary.BigEndian.Uint16(frame[first:last]))
                    }
                }
            }
        }
    }
    return result, nil
}

func (p *ProtoParser) parseDbcFile(content string) {
    lines :=strings.Split(content, "\n")
    for i := 0; i <= len(lines) - 1; i++ {
        line := lines[i]
        if strings.HasPrefix(line, "BO_") {

            v := strings.Split(line, " ")[1:]
            id, _ := strconv.Atoi(v[0])
            length, _ := strconv.Atoi(v[2])
            tmpMsg := CanMsg{
                Name: v[1] ,
                Id: id,
                Sender: v[3],
                Length: uint8(length),
            }
            j := 1 
            for {
                if !strings.HasPrefix(lines[i + j], " SG_") {
                    break
                }
                tmpMsg.Signals = append(tmpMsg.Signals, parseSG(lines[i+j]))
                j++
            }
            p.Msgs = append(p.Msgs, tmpMsg)
        }
    }
}

func parseSG(s string) CanSignal {
	line := strings.TrimSpace(s)
	parts := strings.Split(line, " ")
    
    min := strings.TrimPrefix(strings.Split(parts[6], "|")[0], "[")
    max := strings.TrimSuffix(strings.Split(parts[6], "|")[1], "]")
    mi, _ := strconv.Atoi(min)
    ma, _ := strconv.Atoi(max)

    scale := strings.TrimPrefix(strings.Split(parts[5], ",")[0], "(")
    offset := strings.TrimPrefix(strings.Split(parts[5], ",")[1], ")")
    sc, _ := strconv.Atoi(scale)
    of, _ := strconv.Atoi(offset)

    bitStart := strings.Split(parts[4], "|")[0]
    length := strings.Split(strings.Split(parts[4], "|")[1], "@")[0]
    bS, _ := strconv.Atoi(bitStart)
    l, _ := strconv.Atoi(length) 
    

    littleEndian := strings.Split(strings.Split(parts[4], "|")[1], "@")[1]
    lE := strings.HasPrefix(littleEndian, "1")
    sig := strings.HasSuffix(littleEndian, "+")
    return CanSignal{
        Name: parts[1],
        Reciever: parts[8],
        Unit: parts[7],
        Min: mi,
        Max: ma,
        Scale: float64(sc),
        Offset: float64(of),
        BitStart: uint8(bS),
        Length: uint8(l),
        LittleEndian: lE,
        Signed: sig,
    }
}

func (p *ProtoParser) Load(f string) (Parser, error){
    data, err := os.ReadFile(f)
    if err != nil {
        return nil, err
    }
    p.parseDbcFile(string(data)) 
    return p, nil
}
 
func New() Parser {
    return &ProtoParser{}
}

