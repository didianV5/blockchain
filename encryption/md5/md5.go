package main

import(
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

func main(){
	data := []byte("hello world")
    s := fmt.Sprintf("%x", md5.Sum(data))
    fmt.Println(s)

    // 也可以用这种方式
    h := md5.New()
    h.Write(data)
    s = hex.EncodeToString(h.Sum(nil))
    fmt.Println(s)
}

