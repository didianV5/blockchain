package main

import(
	"fmt"
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func main() {

	// 第一种调用方法
	sum := sha256.Sum256([]byte("hello world\n"))
	fmt.Printf("%x\n", sum)

	// 第二种调用方法
	h := sha256.New()
	h.Write([]byte("hello world\n"))
	fmt.Printf("%x\n", h.Sum(nil))

	// 对文件加密
	f, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h = sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", h.Sum(nil))
}