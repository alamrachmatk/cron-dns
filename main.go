package main

import (
	"bufio"
	"dns/config"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var files []string
	root := config.DnsLog
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		log.Println(err)
	}
	for a, getFile := range files {
		if a != 0 {
			file, err := os.Open(config.DnsLog + filepath.Base(getFile)) // For read access.
			partsFilename := strings.Split(file.Name(), "_")
			partsDate := strings.Split(partsFilename[2], "-")
			//get date
			date := partsDate[2] + "-" + partsDate[1] + "-" + partsDate[0]
			checkQuestion := "question for '"
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			var txtlines []string

			for scanner.Scan() {
				txtlines = append(txtlines, scanner.Text())
			}

			file.Close()

			var ipAddress string
			var domain string
			var baseDomain string
			var time string
			var dateTime string
			for _, eachline := range txtlines {
				if strings.Contains(eachline, checkQuestion) {
					splitTime := strings.Split(eachline, " ")
					time = splitTime[2]
					result := strings.SplitAfter(eachline, checkQuestion)
					for i := range result {
						if i == 1 {
							dateTime = date + " " + time
							rep := strings.ReplaceAll(result[1], ".|A' from ", "|")
							rep2 := strings.ReplaceAll(rep, ".|AAAA' from ", "|")
							rep3 := strings.ReplaceAll(rep2, ".|TYPE65' from ", "|")
							parts := strings.Split(rep3, "|")
							// get ip addreess & domain & basedomain
							for i := range parts {
								if i == 0 {
									domain = parts[0]
									baseDomain = domainutil.Domain(domain)
									log.Println("HasSubdomain: ", domainutil.HasSubdomain(domain))
								}
								if i == 0 {
									ipAddress = parts[1]
								}
							}
							log.Println("IP Address: ", ipAddress)
							log.Println("Domain: ", domain)
							log.Println("Base Domain: ", baseDomain)
							log.Println("Date Time: ", dateTime)
						}
					}

				}
			}
		}
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}
