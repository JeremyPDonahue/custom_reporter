package report

// ReportData contains the data that will populate the Excel spreadsheet
type ReportData struct {
	appID         string
	teamName      string
	opStatus      string
	appTier       string
	emailDL       string
	managerName   string
	managerEmail  string
	directorName  string
	directorEmail string
	vpName        string
	vpEmail       string
	clusters      []string
}

// KubectlCommand hold the piece of the kubectl command used to query NSK clusters
type KubectlCommand struct {
	kubectl          string
	kubeconfigFlag   string
	pathToKubeConfig string
	get              string
	ns               string
}

// JsonResponse hold the data returned by the ServiceNow API call
type JsonResponse struct {
	Result []struct {
		SysID        string `json:"sys_id"`
		UAppID       string `json:"u_app_id"`
		Name         string `json:"name"`
		SupportGroup struct {
			Value             string `json:"value"`
			DisplayValue      string `json:"displayValue"`
			UDistributionList string `json:"u_distribution_list"`
			UAdSecurityGroup  string `json:"u_ad_security_group"`
			Manager           struct {
				Value        string `json:"value"`
				DisplayValue string `json:"displayValue"`
				Email        string `json:"email"`
				Active       bool   `json:"active"`
			} `json:"manager"`
			UDirector struct {
				Value        string `json:"value"`
				DisplayValue string `json:"displayValue"`
				Email        string `json:"email"`
				Active       bool   `json:"active"`
			} `json:"u_director"`
			UVp struct {
				Value        string `json:"value"`
				DisplayValue string `json:"displayValue"`
				Email        string `json:"email"`
				Active       bool   `json:"active"`
			} `json:"u_vp"`
		} `json:"support_group"`
		UApplicationTier struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"u_application_tier"`
		DataClassification struct {
			Value        interface{} `json:"value"`
			DisplayValue string      `json:"displayValue"`
		} `json:"data_classification"`
		OperationalStatus struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"operational_status"`
		ULongName               string `json:"u_long_name"`
		ShortDescription        string `json:"short_description"`
		UHlq                    string `json:"u_hlq"`
		UPci                    bool   `json:"u_pci"`
		UPii                    bool   `json:"u_pii"`
		USox                    bool   `json:"u_sox"`
		USoc1                   bool   `json:"u_soc1"`
		USoc1Bank               bool   `json:"u_soc1_bank"`
		UHipaa                  bool   `json:"u_hipaa"`
		URequireApproval        bool   `json:"u_require_approval"`
		Comments                string `json:"comments"`
		LastUpdated             string `json:"last_updated"`
		LastUpdatedBy           string `json:"last_updated_by"`
		UMostRecentDesignReview struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"u_most_recent_design_review"`
		ULastSecurityAssessment struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"u_last_security_assessment"`
		ChangeControl struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"change_control"`
		AssignmentGroup struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"assignment_group"`
		UApplicationCategory struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"u_application_category"`
		Parent struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"parent"`
		UNordstromOrganization struct {
			Value        string `json:"value"`
			DisplayValue string `json:"displayValue"`
		} `json:"u_nordstrom_organization"`
		IncCount interface{} `json:"incCount"`
		ChgCount interface{} `json:"chgCount"`
		PrbCount interface{} `json:"prbCount"`
	} `json:"result"`
}
