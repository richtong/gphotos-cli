// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/trinhdrew1418/gphotos-cli/utils"
	"github.com/trinhdrew1418/gphotos-cli/utils/expobackoff"
	"google.golang.org/api/photoslibrary/v1"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	DownloadDir string
)

type Date struct {
	day   int
	month int
	year  int
}

type DownloadTask struct {
	url      string
	location string
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsDir(DownloadDir) {
			println("You provided an invalid download directory")
			os.Exit(1)
		}

		_, gphotoService := getClientService(photoslibrary.PhotoslibraryScope)
		var noDate bool
		var noCat bool

		startDate, endDate := GetDates(&noDate)
		println()
		categories := GetCategories(&noCat)
		println()

		filt := photoslibrary.Filters{}
		if !noCat {
			filt.ContentFilter = &photoslibrary.ContentFilter{IncludedContentCategories: categories}
		}

		if !noDate {
			filt.DateFilter = &photoslibrary.DateFilter{
				Ranges: []*photoslibrary.DateRange{
					{
						StartDate: &photoslibrary.Date{
							Day:   int64(startDate.day),
							Month: int64(startDate.month),
							Year:  int64(startDate.year),
						},
						EndDate: &photoslibrary.Date{
							Day:   int64(endDate.day),
							Month: int64(endDate.month),
							Year:  int64(endDate.year),
						},
					},
				},
			}
		}

		resp, err := gphotoService.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{Filters: &filt}).Do()
		if err != nil {
			log.Fatal("Failed media search")
		}

		dTaskFeed := make(chan DownloadTask)
		feedPage(resp, dTaskFeed)
		for resp.NextPageToken != "" {
			req := gphotoService.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{PageToken: resp.NextPageToken}).Do
			resp, err = req()

			if err != nil || resp.HTTPStatusCode == 429 {
				durations := expobackoff.Calculate(expobackoff.NUM_RETRIES)
				for _, sleepDur := range durations {
					duration := time.Duration(sleepDur)
					time.Sleep(duration)
					resp, err = req()
					if err == nil && resp.HTTPStatusCode == 200 {
						break
					}
				}
				println("Unable to get next page")
				os.Exit(1)
			}

			feedPage(resp, dTaskFeed)
		}
	},
}

func guaranteeDirectory(mItem *photoslibrary.MediaItem) (string, error) {
	creationParts := strings.Split(mItem.MediaMetadata.CreationTime, "-")
	year := creationParts[0]
	month := creationParts[1]
	loc := path.Join(DownloadDir, year, month)
	err := os.MkdirAll(loc, os.ModePerm)
	return loc, err
}

func feedPage(resp *photoslibrary.SearchMediaItemsResponse, dTaskFeed chan DownloadTask) {
	for _, mItem := range resp.MediaItems {
		dest, err := guaranteeDirectory(mItem)
		if err != nil {
			log.Fatal(err)
		}

		dTaskFeed <- DownloadTask{path.Join(mItem.BaseUrl, "-d"), dest}
	}
}

func initDownloads(dTaskFeed chan DownloadTask) {
	for i := 0; i < MAX_WORKERS; i++ {
		go Downloader(dTaskFeed)
	}
}

func Downloader(dTaskFeed chan DownloadTask) {
	defer close(dTaskFeed)

	for task := range dTaskFeed {
		resp, err := http.Get(task.url)
		for err != nil {
			println("Download failed", resp.StatusCode)
			os.Exit(1)
		}
		defer resp.Body.Close()

	}
}

func GetCategories(noCat *bool) []string {
	allCategories := map[string]bool{
		"ANIMALS":      true,
		"LANDMARKS":    true,
		"PETS":         true,
		"UTILITY":      true,
		"BIRTHDAYS":    true,
		"LANDSCAPES":   true,
		"RECEIPTS":     true,
		"WEDDINGS":     true,
		"CITYSCAPES":   true,
		"NIGHT":        true,
		"SCREENSHOTS":  true,
		"WHITEBOARDS":  true,
		"DOCUMENTS":    true,
		"PEOPLE":       true,
		"SELFIES":      true,
		"FOOD":         true,
		"PERFORMANCES": true,
		"SPORT":        true}

	println("Here are the available categories: ")
	for _, cat := range allCategories {
		println(" - ", cat)
	}
	println(" - ALL")

	var parseString string
	print("Select up to 10 categories (ALL if you want any) [capital or lowercase, separate by spaces]: ")
	_, err := fmt.Scan(&parseString)
	if err != nil {
		log.Fatal("Unable to obtain categories")
	}

	categories := strings.Split(parseString, " ")
	if len(categories) > 10 {
		log.Fatal("Too many categories")
	}

	for i, str := range categories {
		str = strings.ToUpper(str)

		if str == "ALL" {
			*noCat = true
			return make([]string, 0)
		}

		if !allCategories[str] {
			log.Fatal("Invalid category")
		}

		categories[i] = str
	}

	return categories
}

func GetDates(noDate *bool) (Date, Date) {
	dateOptions := map[string]int{
		"Today":                       0,
		"Past week":                   1,
		"Past month":                  2,
		"Specific number of days ago": 3,
		"Specific date range":         4,
		"Any":                         5,
	}

	var startDate Date
	var endDate Date

	pmpt := promptui.Select{
		Label: "Select a time frame: ",
		Items: []string{"Today", "Past week", "Past month", "Specific number of days ago",
			"Specific date range", "Any"},
	}

	_, resp, err := pmpt.Run()
	if err != nil {
		log.Fatal("Prompt failed to open")
	}

	switch dateOptions[resp] {
	case 0:
		startDate = getSomeDaysAgo(0)
		endDate = startDate
	case 1:
		endDate = getSomeDaysAgo(0)
		startDate = getSomeDaysAgo(7)
	case 2:
		endDate = getSomeDaysAgo(0)
		startDate = getSomeDaysAgo(30)
	case 3:
		var numDays int
		print("Number of days: ")

		_, err := fmt.Scan(&numDays)
		if err != nil {
			log.Fatal("Could not read the amount of days")
		}

		endDate = getSomeDaysAgo(0)
		startDate = getSomeDaysAgo(numDays)

	case 4:
		var sDate string
		var eDate string

		print("Start Date [MM-DD-YYYY]:")
		_, err := fmt.Scan(&sDate)
		if err != nil {
			log.Fatal("Unable to read  response")
		}

		stringDate := strings.Split(sDate, "-")
		startDate.month, _ = strconv.Atoi(stringDate[0])
		startDate.day, _ = strconv.Atoi(stringDate[1])
		startDate.year, _ = strconv.Atoi(stringDate[2])

		_, err = fmt.Scan(&sDate)
		if err != nil {
			log.Fatal("Unable to read response")
		}

		print("Start End Date [MM-DD-YYYY]:")
		stringDate = strings.Split(eDate, "-")

		endDate.month, _ = strconv.Atoi(stringDate[0])
		endDate.day, _ = strconv.Atoi(stringDate[1])
		endDate.year, _ = strconv.Atoi(stringDate[2])

		print("The following range will be listed: ", sDate, " to ", eDate)
	case 5:
		*noDate = true
	}

	return startDate, endDate
}

func getSomeDaysAgo(num int) Date {
	Months := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}

	var retDate Date
	year, month, day := time.Now().AddDate(0, 0, -num).Date()
	retDate.year = year
	retDate.month = Months[month.String()]
	retDate.month = Months[month.String()]
	retDate.day = day

	return retDate
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	pullCmd.PersistentFlags().StringVar(&DownloadDir, "download directory path", "./", "Define the directory you want to download your files to")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
