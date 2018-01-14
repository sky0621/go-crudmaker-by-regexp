package main

import (
	"fmt"
	"os"
	"path/filepath"

	"regexp"

	"bufio"

	"strings"

	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------------------------------------
// とりあえず第１弾は、controller層をなめて、パス毎の使用サービス（及びメソッド）をマップ化するところまで
// ---------------------------------------------------------------------------------------------------------

var targetDir = filepath.Join("fuel", "app", "classes", "controller", ".*\\.php")

const (
	Service_Aws_DynamoDB    = "Service_Aws_DynamoDB"
	Service_Aws_ElastiCache = "Service_Aws_ElastiCache"
	Service_Aws_Kms         = "Service_Aws_Kms"
	Service_Aws_S3          = "Service_Aws_S3"
	Service_Aws_Sqs         = "Service_Aws_Sqs"

	Service_RDB = "Service_"
	Service_AWS = "Service_Aws_"
)

var regexpserviceRDB *regexp.Regexp = nil

type CRUD struct {
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if len(os.Args) < 2 {
		logger.Error("引数[ターゲットディレクトリのパス]が必要です")
		os.Exit(-1)
	}
	target := os.Args[1]

	regexpserviceRDB, err = regexp.Compile("Service\\_\\A")
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(target, Apply)
	if err != nil {
		logger.Error("", zap.String("error", err.Error()))
		os.Exit(-1)
	}
}

// Apply ...
func Apply(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !filter(path, info) {
		return nil
	}

	fmt.Println("####################################################")
	fmt.Println(path)
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

	inComment := false

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		txt2 := strings.Trim(txt, "\t")

		if strings.HasPrefix(txt2, "/*") && strings.HasSuffix(txt2, "*/") {
			inComment = false
			continue
		}

		if strings.HasPrefix(txt2, "/*") {
			inComment = true
			continue
		}

		if strings.HasPrefix(txt2, "*/") {
			inComment = false
			continue
		}

		if inComment {
			continue
		}

		if strings.HasPrefix(txt2, "//") || strings.HasPrefix(txt2, "*") {
			continue
		}

		if txt2 == "" {
			continue
		}

		if strings.Contains(txt2, Service_Aws_DynamoDB) {
			fmt.Printf("[Service_Aws_DynamoDB] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_ElastiCache) {
			fmt.Printf("[Service_Aws_ElastiCache] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_Kms) {
			fmt.Printf("[Service_Aws_Kms] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_S3) {
			fmt.Printf("[Service_Aws_S3] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_Aws_Sqs) {
			fmt.Printf("[Service_Aws_Sqs] %s\n", txt2)
		}
		if strings.Contains(txt2, Service_RDB) && !strings.Contains(txt2, Service_AWS) {
			fmt.Printf("[Service_RDB] %s\n", txt2)
		}
		// if regexpserviceRDB.MatchString(txt2) {
		// 	fmt.Printf("[regexpserviceRDB] %s\n", txt2)
		// }
		// fmt.Println(txt2)
	}

	return nil
}

func filter(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	outDirExp, err := regexp.Compile("vendor")
	if err != nil {
		return false
	}
	if outDirExp.MatchString(absPath) {
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
	if !inFileExp.MatchString(path) {
		return false
	}

	return true
}
