package main

import (
	"compress/gzip"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-gitea/gitea/modules/markup"
	"github.com/gophish/gophish/config"
	"golang.org/x/crypto/md4"
)

var validConfig = []byte(`{
	"admin_server": {
		"listen_url": "127.0.0.1:3333",
		"use_tls": true,
		"cert_path": "gophish_admin.crt",
		"key_path": "gophish_admin.key"
	},
	"phish_server": {
		"listen_url": "0.0.0.0:8080",
		"use_tls": false,
		"cert_path": "example.crt",
		"key_path": "example.key"
	},
	"db_name": "sqlite3",
	"db_path": "gophish.db",
	"migrations_prefix": "db/db_",
	"contact_address": ""
}`)

func main() {
	h := md4.New()
	data := "These pretzels are making me thirsty."
	io.WriteString(h, data)
	fmt.Printf("MD4 is the new MD5: %x\n", h.Sum(nil))

	err := ioutil.WriteFile("config/phish-config.json", validConfig, 0644)
	conf := config.Config{}
	fmt.Printf("GONE PHISH'N for configs %v, maybe error: %v\n", conf, err)

	fmt.Printf("IS A README?? : %v as per gitea\n", markup.IsReadmeFile("README.md"))

	const password = "secret123"
	if password == "secret123" {
		fmt.Println("Access granted!")
		fmt.Println("You may pass!")
		fmt.Println("You may pass again!")
	}

	hash := md5.New()
	hash.Write([]byte("test"))
	fmt.Printf("%x", hash.Sum(nil))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Query().Get("path")
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	userInput := "ls -l; rm -rf ./" // NOTE: We are not going to erase the whole hard drive; at worst, we will erase the current directory
	cmd := exec.Command("sh", "-c", userInput)
	cmd.Run()

	f, _ := os.Open("file.txt")
	defer f.Close()

	username := "admin"
	pass := "' OR 1=1--"
	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, pass)
	db, _ := sql.Open("mysql", "user:password@/dbname")
	db.Exec(query)

	key := []byte("my-key")
	block, _ := des.NewCipher(key)
	fmt.Printf("%x", block)

	config := &tls.Config{
		MinVersion: tls.VersionSSL30,
	}
	_, _ = tls.Dial("tcp", "example.com:443", config)

	cipher, _ := rc4.NewCipher([]byte("secret"))
	fmt.Printf("%x", cipher)

	resp, err := http.Get("http://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	url := resp.Request.URL.Query().Get("url")
	http.Get(url)

	val := resp.Request.URL.Query().Get("val")
	num, _ := strconv.Atoi(val)
	var intVal int16 = int16(num)
	fmt.Println(intVal)

	http.HandleFunc("/decompress", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<30) // 1GB
		gzr, _ := gzip.NewReader(r.Body)
		_, _ = io.Copy(os.Stdout, gzr)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
