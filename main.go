package main

import (
	"bufio"
	"dns/config"
	"dns/db"
	"dns/models"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/jasonlvhit/gocron"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize main database
	db.Db = db.MariaDBInit()

	gocron.Every(30).Seconds().Do(LogDns)

	<-gocron.Start()
}

func LogDns() {
	log.Println("RUNNING CRON DNS")
	var files []string
	root := config.DnsLog
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		log.Println(err)
	} else {
		for a, getFile := range files {
			if a != 0 {
				fileName := config.DnsLog + filepath.Base(getFile)
				log.Println("open file: ", filepath.Base(getFile))
				var extension = filepath.Ext(fileName)
				if extension != ".filepart" {
					file, err := os.Open(fileName) // For read access.
					partsFilename := strings.Split(file.Name(), "_")
					partsDate := strings.Split(partsFilename[2], "-")
					//get date
					date := partsDate[2] + "-" + partsDate[1] + "-" + partsDate[0]
					checkQuestion := "question for '"
					if err != nil {
						log.Println(err)
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
					var hasSubDomain string
					var time string
					var dateTime string
					params := make(map[string]string)
					for _, eachline := range txtlines {
						if strings.Contains(eachline, checkQuestion) {
							splitTime := strings.Split(eachline, " ")
							time = splitTime[2]
							result := strings.SplitAfter(eachline, checkQuestion)
							for i := range result {
								if i == 1 {
									dateTime = date + " " + time
									currentTime := time.Now()
									dateNow := currentTime.Format("2006-01-02")
									if date == dateNow {
										var partDomain, partAfterDomain string
										if i := strings.Index(result[1], ".|"); i >= 0 {
											partDomain, partAfterDomain = result[1][:i], result[1][i:]
										}
										partIP := strings.Split(partAfterDomain, "from ")
										rep := partDomain + "|" + partIP[1]
										parts := strings.Split(rep, "|")
										// get ip addreess & domain & basedomain
										for i := range parts {
											if i == 0 {
												domain = parts[0]
												baseDomain = domainutil.Domain(domain)
												if baseDomain == "" {
													baseDomain = domain
												}
												if domainutil.HasSubdomain(domain) == true {
													hasSubDomain = "1"
												} else {
													hasSubDomain = "0"
												}
											}
											if i == 0 {
												ipAddress = parts[1]
											}
										}
										params["domain"] = domain
										params["base_domain"] = baseDomain
										params["ip_address"] = ipAddress
										params["has_subdomain"] = hasSubDomain
										params["log_datetime"] = dateTime
										statusResponse, err := models.CreateDns(params)
										if statusResponse != 200 {
											log.Println(err)
										}
									} else {
										return
									}
								}
							}

						}
					}
					e := os.Remove(fileName)
					if e != nil {
						log.Println(e)
					}
					return
				}
			}
		}
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		*files = append(*files, path)
		return nil
	}
}
