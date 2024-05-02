package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"internal.gitlab/projectID/Operations/scripts/custom_reporter/servicenow"
)

// GoGrabServiceNowData makes API requests to ServiceNow using appid/ns
func GoGrabServiceNowData(clusterNamespaces [][]string) []ReportData {

	fmt.Println("Making requests to SN and parsing the responses... (This may take a minute)")

	var customerData []ReportData

	creds, err := servicenow.GetCredentialsFromEnvironment()
	if err != nil {
		e := fmt.Errorf("could not retrieve creds from environment")
		fmt.Println(e.Error())
		log.Fatal(err)
	}

	client := servicenow.NewDefaultClient(creds)

	for _, slice := range clusterNamespaces {

		for _, appid := range slice {

			if a := isItACluster.FindString(appid); a != "" {
				continue
			} else {
				fmt.Print(".")
				var _, _, jsonString, err = client.GetApplication(appid)

				if err != nil {
					e := fmt.Errorf("error retrieving data from ServiceNow")
					fmt.Println(e.Error())
					log.Fatal(err)
				}

				r := PopulateReportStruct(jsonString, clusterNamespaces)

				customerData = append(customerData, *r)

			}
		}
	}

	return customerData

}

// WhichClusters is used within PopulateReportStruct to help figure out which NSK clusters an AppID/Namespace lives on
func WhichClusters(clusterNamespaces [][]string, searchingFor string) []string {

	var whichClusters []string

	// looping through each slice within a slice of slices
Search:
	for _, slice := range clusterNamespaces {
		// looping through each appid within a slice
		for _, appid := range slice {
			if appid == searchingFor {
				whichClusters = append(whichClusters, slice[0])
				continue Search
			}
		}
	}

	return whichClusters
}

// PopulateReportStruct uses the data gathered by GoGrabServiceNowData to populate the ReportData struct
func PopulateReportStruct(resp string, clusterNamespace [][]string) *ReportData {

	var JS JsonResponse

	if err := json.Unmarshal([]byte(resp), &JS); err != nil {
		fmt.Println("Could not unmarshal the json data")
		log.Fatal(err)
	}

	result := JS.Result[0]

	a := strings.ToLower(result.UAppID)
	clusters := WhichClusters(clusterNamespace, a)

	return &ReportData{
		appID:         result.UAppID,
		teamName:      result.Name,
		opStatus:      result.OperationalStatus.DisplayValue,
		appTier:       result.UApplicationTier.DisplayValue,
		emailDL:       result.SupportGroup.UDistributionList,
		managerName:   result.SupportGroup.Manager.DisplayValue,
		managerEmail:  result.SupportGroup.Manager.Email,
		directorName:  result.SupportGroup.UDirector.DisplayValue,
		directorEmail: result.SupportGroup.UDirector.Email,
		vpName:        result.SupportGroup.UVp.DisplayValue,
		vpEmail:       result.SupportGroup.UVp.Email,
		clusters:      clusters,
	}

}

// CreateCsvFile creates a CSV file with all of the Namespace data that this application gathers
func CreateCsvFile(r []ReportData, clusterList []string) {

	now := time.Now()
	timeString := now.Format("2006-01-02")
	fileName := fmt.Sprintf("namespace_report_%s.csv", timeString)

	csvFile, err := os.Create(fileName)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	headerRow := []string{
		"Namespace", "Team", "Ops Status",
		"Tier", "Distribution List", "Manager", "Manager Email",
		"Director", "Director Email", "VP", "VP Email",
	}

	fmt.Println("\nGenerating report...")

	for _, cluster := range clusterList {
		headerRow = append(headerRow, strings.Trim(cluster, ".yaml"))
	}

	_ = csvwriter.Write(headerRow)

	var used []string

	for _, entry := range r {

		if Contains(used, entry.appID) {
			continue
		} else {

			s := []string{
				entry.appID,
				entry.teamName,
				entry.opStatus,
				entry.appTier,
				entry.emailDL,
				entry.managerName,
				entry.managerEmail,
				entry.directorName,
				entry.directorEmail,
				entry.vpName,
				entry.vpEmail,
			}

			i := 0

			for _, cluster := range clusterList {

				if Contains(entry.clusters, clusterList[i]) == true {
					cluster = "x"
					s = append(s, cluster)
				} else {
					s = append(s, "")
				}

				i++
			}

			err = csvwriter.Write(s)
			if err != nil {
				e := fmt.Errorf("error writing to CSV file")
				fmt.Println(e.Error())
				log.Fatal(err)
			}
		}

		used = append(used, entry.appID)
	}

	csvwriter.Flush()
	csvFile.Close()

}
