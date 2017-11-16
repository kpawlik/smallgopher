package worker

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// Acp holds I/O buffer to communicate with Magik ACP
type Acp struct {
	Name string
	io   *bufio.ReadWriter
}

// NewAcp creates and init new Acp with name
func NewAcp(name string) *Acp {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	inOut := bufio.NewReadWriter(reader, writer)
	return &Acp{Name: name, io: inOut}
}

// Flush send buffer data
func (a *Acp) Flush() {
	a.io.Flush()
}

// Write writes buffer to Acp output
func (a *Acp) Write(buf []byte) {
	a.io.Write(buf)
	a.Flush()
}

// PutBool sends boolean value to Acp output
func (a *Acp) PutBool(b bool) {
	var iVal byte
	if b {
		iVal = 1
	}
	a.io.WriteByte(iVal)
	a.Flush()
}

// PutUbyte sends unsigned byte value to Acp output
func (a *Acp) PutUbyte(value uint8) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutByte sends  byte value to Acp output
func (a *Acp) PutByte(value int8) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutUShort sends unsigned short value to Acp output
func (a *Acp) PutUShort(value uint16) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutShort sends short value to Acp output
func (a *Acp) PutShort(value int16) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutUint sends int value to Acp output
func (a *Acp) PutUint(value uint32) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutInt sends int value to Acp output
func (a *Acp) PutInt(value int32) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutULong sends unsigned long value to Acp output
func (a *Acp) PutULong(value uint64) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutLong sends long value to Acp output
func (a *Acp) PutLong(value int64) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutShortFloat sends short float value to Acp output
func (a *Acp) PutShortFloat(value float32) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutFloat sends float value to Acp output
func (a *Acp) PutFloat(value float64) {
	binary.Write(a.io, binary.LittleEndian, value)
	a.Flush()
}

// PutString sends string value to Acp output
func (a *Acp) PutString(s string) {
	bytes := []byte(s)
	l := len(bytes)
	a.PutUShort(uint16(l))
	a.Write(bytes)
}

// read bytes from Acp input
func (a *Acp) readBytes(n int) (buf []byte, err error) {
	buf = make([]byte, 0, n)
	_, err = a.io.Read(buf[:cap(buf)])
	if err != nil && err != io.EOF {
		err = fmt.Errorf("Error reading short: %v\n", err)
		return
	}
	buf = buf[:n]
	return
}

// ReadNumber reads number from Acp input
func (a *Acp) ReadNumber(data interface{}) {
	if err := binary.Read(a.io, binary.LittleEndian, data); err != nil {
		log.Fatal(err)
	}
}

// GetBool reads boolean value from Acp input
func (a *Acp) GetBool() bool {
	b, err := a.io.ReadByte()
	if err != nil {
		log.Fatalf("Error reading boolean %v\n", err)
	}
	return b == 1
}

// GetUbyte reads unsigned byte from Acp input
func (a *Acp) GetUbyte() int {
	var res uint8
	a.ReadNumber(&res)
	return int(res)
}

// GetByte reads byte from Acp input
func (a *Acp) GetByte() int {
	var res int8
	a.ReadNumber(&res)
	return int(res)
}

// GetUShort reads unsigned short from Acp input
func (a *Acp) GetUShort() int {
	var res uint16
	a.ReadNumber(&res)
	return int(res)
}

// GetShort reads short from Acp input
func (a *Acp) GetShort() int {
	var res int16
	a.ReadNumber(&res)
	return int(res)
}

// GetUint reads unsigned int from Acp input
func (a *Acp) GetUint() int {
	var res uint32
	a.ReadNumber(&res)
	return int(res)
}

// GetInt reads unsigned int from Acp input
func (a *Acp) GetInt() int {
	var res int32
	a.ReadNumber(&res)
	return int(res)
}

// GetULong reads unsigned long from Acp input
func (a *Acp) GetULong() uint64 {
	var res uint64
	a.ReadNumber(&res)
	return res
}

// GetLong reads long from Acp input
func (a *Acp) GetLong() int64 {
	var res int64
	a.ReadNumber(&res)
	return res
}

// GetShortFloat read float32 from Acp input
func (a *Acp) GetShortFloat() float32 {
	var res float32
	a.ReadNumber(&res)
	return res
}

// GetFloat read float64 from Acp input
func (a *Acp) GetFloat() float64 {
	var res float64
	a.ReadNumber(&res)
	return res
}

// GetString reads string from Acp input
func (a *Acp) GetString() string {
	b := a.GetUShort()
	buf, err := a.readBytes(b)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

// GetCoord return coordinate
func (a *Acp) GetCoord() [2]float64 {
	var (
		res [2]float64
	)
	res[0] = a.GetFloat()
	res[1] = a.GetFloat()
	return res

}

// VerifyConnection verify Acp process name
func (a *Acp) VerifyConnection(name string) bool {
	acpName := a.GetString()
	res := acpName == name
	log.Printf("Name: %s; From SW: %s; Res: %v\n", name, acpName, res)
	a.PutBool(!res)
	return res
}

// EstablishProtocol checks Acp protocol
func (a *Acp) EstablishProtocol(minProtocol, maxProtocol int) bool {
	min, max := minProtocol, maxProtocol
	for {
		in := a.GetUShort()
		log.Printf("Protocol from SW: %d\n", in)
		if in < min || in > max {
			a.PutBool(true)
			a.PutUShort(uint16(min))
			a.PutUShort(uint16(max))
			return false
		}
		break
	}
	a.PutBool(false)
	return true
}

// Connect verify connection and protocol to Acp
func (a *Acp) Connect(processName string, protocolMin, protocolMax int) (err error) {
	log.Printf("ACP started name: %s\n", processName)
	if res := a.VerifyConnection(processName); !res {
		err = errors.New("Error verify connection")
		return
	}
	log.Println("Connection verified")
	if res := a.EstablishProtocol(protocolMin, protocolMax); !res {
		err = errors.New("Error establish protocol")
		return
	}
	log.Println("Protocol established")
	log.Println("Connected")
	return
}

// Put convert value to dataType and send this value to ACP
func (a *Acp) Put(dataType string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	switch dataType {
	case "boolean":
		iVal := value.(bool)
		a.PutBool(iVal)
	case "unsigned_byte":
		iVal := value.(uint8)
		a.PutUbyte(iVal)
	case "signed_byte":
		iVal := value.(int8)
		a.PutByte(iVal)
	case "unsigned_short":
		iVal := value.(uint16)
		a.PutUShort(iVal)
	case "signed_short":
		iVal := value.(int16)
		a.PutShort(iVal)
	case "unsigned_int":
		iVal := value.(uint32)
		a.PutUint(iVal)
	case "signed_int":
		iVal := value.(int32)
		a.PutInt(iVal)
	case "unsigned_long":
		iVal := value.(uint64)
		a.PutULong(iVal)
	case "signed_long":
		iVal := value.(int64)
		a.PutLong(iVal)
	case "short_float":
		iVal := value.(float32)
		a.PutShortFloat(iVal)
	case "float":
		iVal := value.(float64)
		a.PutFloat(iVal)
	case "chars":
		iVal := value.(string)
		a.PutString(iVal)
	default:
		return fmt.Errorf("Unsupported data type '%s' in Put method", dataType)
	}
	return nil

}

// Get method reads dataType value from ACP
func (a *Acp) Get(dataType string) (value interface{}, err *AcpErr) {
	defer func() {
		if r := recover(); r != nil {
			err = NewAcpErr(fmt.Sprint(r.(error)))
			return
		}
	}()
	switch dataType {
	case "boolean":
		return a.GetBool(), nil
	case "unsigned_byte":
		return a.GetUbyte(), nil
	case "signed_byte":
		return a.GetByte(), nil
	case "unsigned_short":
		return a.GetUShort(), nil
	case "signed_short":
		return a.GetShort(), nil
	case "unsigned_int":
		return a.GetUint(), nil
	case "signed_int":
		return a.GetInt(), nil
	case "unsigned_long":
		return a.GetULong(), nil
	case "signed_long":
		return a.GetLong(), nil
	case "short_float":
		return a.GetShortFloat(), nil
	case "float":
		return a.GetFloat(), nil
	case "chars":
		return a.GetString(), nil
	default:
		return nil, NewAcpErr(fmt.Sprintf("Unsupported data type '%s' in Get method", dataType))
	}
}
