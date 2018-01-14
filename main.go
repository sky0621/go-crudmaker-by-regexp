package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"regexp"

	"bufio"

	"strings"
)

// TODO 機能実現スピード最優先での実装なので要リファクタ

// ---------------------------------------------------------------------------------------------------------
// とりあえず第１弾は、controller層とbatch層をなめて、パス毎の使用サービス（及びメソッド）をマップ化するところまで
//
// 【第２弾以降でやること】
// ・とにかくちゃんと検討（データ構造やら汎用性やら）
// ・controller層とbatch層から直接service層を読んでいないケースを想定
// ・CRUD特定のための表現を正規表現化
// ・CRUD特定のための表現等の情報を設定ファイル保持化
// ・各環境用のバイナリを release 化
// ・Docker化
// ・goroutine使用
// ・パーサージェネレーター使用？
// ・エラーハンドリング
// ・テストコード
// ・ロギング（zap）
// ---------------------------------------------------------------------------------------------------------

var targetDir = filepath.Join("fuel", "app", "classes", "controller", ".*\\.php")
var targetDir2 = filepath.Join("fuel", "app", "classes", "batch", ".*\\.php")

const (
	Service_Aws_DynamoDB    = "Service_Aws_DynamoDB"
	Service_Aws_ElastiCache = "Service_Aws_ElastiCache"
	Service_Aws_Kms         = "Service_Aws_Kms"
	Service_Aws_S3          = "Service_Aws_S3"
	Service_Aws_Sqs         = "Service_Aws_Sqs"

	Service_RDB = "Service_"
	Service_AWS = "Service_Aws_"
)

type Result struct {
	Datetime     string
	TargetGroups []*TargetGroup
}

// controller層ないしbatch層
type TargetGroup struct {
	Name     string
	Services []*ServiceAws // DynamoDB/ElastiCache/Kms/S3/Sqs/RDS
}

type ServiceAws struct {
	Name   string // DynamoDB/ElastiCache/Kms/S3/Sqs/RDS
	Tables []*Table
}

type Table struct {
	Name  string
	CRUDs []*CRUD
}

type CRUD struct {
	Name string
}

var batchServices []*ServiceAws = []*ServiceAws{}
var controllerServices []*ServiceAws = []*ServiceAws{}

var targetGroups []*TargetGroup = []*TargetGroup{
	&TargetGroup{Name: "batch", Services: batchServices},
	&TargetGroup{Name: "controller", Services: controllerServices},
}

var outputBuf *bytes.Buffer

func main() {
	if len(os.Args) < 2 {
		fmt.Println("引数[ターゲットディレクトリのパス]が必要です")
		os.Exit(-1)
	}
	target := os.Args[1]

	outputBuf = &bytes.Buffer{}
	defer outputBuf.Reset()
	outputBuf.WriteString(fmt.Sprintf("# 管理画面及びPHPバッチのCRUD(%v 時点)\n\n", time.Now().Format("2006-01-02 15:04")))
	outputBuf.WriteString("#### ※ツール（ https://github.com/sky0621/go-crudmaker-by-regexp ）による自動生成\n\n")
	outputBuf.WriteString("#### ・「controller」層、「batch」層から直接「service」層を呼んでいるケースのみ想定\n\n")
	outputBuf.WriteString("#### ・CRUDの判定については「service」層のメソッドが「get〜〜」なら「READのR」、「insert」を含むなら「CREATEのC」といった恣意的なレベル\n\n")

	err := filepath.Walk(target, Apply)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// tmpl := template.Must(template.ParseFiles("./tmpl.md"))
	// buf := &bytes.Buffer{}
	// err = tmpl.Execute(buf, &Result{Datetime: time.Now().Format("2006-01-02 15:04"), TargetGroups: targetGroups})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(buf.String())
	fmt.Println(outputBuf.String())
}

// Apply ...
func Apply(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !filter(path, info) {
		return nil
	}

	dispPath := strings.Replace(path, "/home/sasaki/work/go/src/oden.dac.co.jp/dialogone/sally/fuel/app/classes/", "", -1)
	outputBuf.WriteString(fmt.Sprintf("##### %v\n\n", dispPath))
	// fmt.Println("####################################################")
	// fmt.Println(path)
	// fmt.Println("####################################################")

	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if fp != nil {
			fp.Close()
		}
	}()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		txt2 := strings.Trim(txt, "\t")

		// TODO 適当すぎ・・・

		if strings.Contains(txt2, Service_Aws_DynamoDB) {
			// if strings.Contains(path, "controller") {
			// 	isFirstSvc := true
			// 	for _, batchService := range batchServices {
			// 		if batchService.Name == "DynamoDB" {
			// 			isFirstSvc = false
			// 			for _, table := range batchService.Tables {
			// 				if strings.Contains(txt2, table.Name) {

			// 				}
			// 			}
			// 		}
			// 	}
			// 	if isFirstSvc {
			// 		batchServices = append(batchServices, &ServiceAws{Name: "DynamoDB"})
			// 	}
			// }
			// if strings.Contains(path, "batch") {

			// }
			ba := strings.Split(txt2, "::")
			ba2 := strings.Split(ba[0], " ")
			outputBuf.WriteString(fmt.Sprintf("[DynamoDB] %v\n\n", ba2[len(ba2)-1]))
			// fmt.Printf("[Service_Aws_DynamoDB] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_ElastiCache) {
			// fmt.Printf("[Service_Aws_ElastiCache] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_Kms) {
			// fmt.Printf("[Service_Aws_Kms] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_S3) {
			// fmt.Printf("[Service_Aws_S3] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_Sqs) {
			// fmt.Printf("[Service_Aws_Sqs] %s\n", txt2)
		}
		if !strings.Contains(txt2, Service_AWS) && strings.Contains(txt2, Service_RDB) {
			// fmt.Printf("[Service_RDB] %s\n", txt2)
		}
	}

	return nil
}

// TODO go用に適当に作ったものをとりあえず持ってきて多少いじっただけ
func filter(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	outDirExp2, err := regexp.Compile("\\.git")
	if err != nil {
		return false
	}
	if outDirExp2.MatchString(absPath) {
		return false
	}

	outFileExp, err := regexp.Compile(".*test.*")
	if err != nil {
		return false
	}
	if outFileExp.MatchString(path) {
		return false
	}

	inFileExp, err := regexp.Compile(targetDir)
	if err != nil {
		return false
	}
	inFileExp2, err := regexp.Compile(targetDir2)
	if err != nil {
		return false
	}
	if !inFileExp.MatchString(path) && !inFileExp2.MatchString(path) {
		return false
	}

	return true
}
