package servicenow

type GetApplicationResponse struct {
	Result []Application `json:"result"`
}

// Application is the SNAPI version of a ServiceNow Application record
type Application struct {
	SysID                  string         `json:"sys_id"`
	AppID                  string         `json:"u_app_id"`
	Name                   string         `json:"name"`
	SupportGroup           Group          `json:"support_group"`
	ApplicationTier        DisplayedValue `json:"u_application_tier"`
	OperationalStatus      DisplayedValue `json:"operational_status"`
	LongName               string         `json:"u_long_name"`
	ShortDescription       string         `json:"short_description"`
	HLQ                    string         `json:"u_hlq"`
	PCI                    bool           `json:"u_pci"`
	PII                    bool           `json:"u_pii"`
	SOX                    bool           `json:"u_sox"`
	SOC1                   bool           `json:"u_soc1"`
	SOC1Bank               bool           `json:"u_soc1_bank"`
	HIPAA                  bool           `json:"u_hipaa"`
	RequireApproval        bool           `json:"u_require_approval"`
	Comments               string         `json:"comments"`
	LastUpdated            string         `json:"last_updated"`
	LastUpdatedBy          string         `json:"last_updated_by"`
	MostRecentDesignReview DisplayedValue `json:"u_most_recent_design_review"`
	LastSecurityAssessment DisplayedValue `json:"u_last_security_assessment"`
	ChangeControl          DisplayedValue `json:"change_control"`
	AssignmentGroup        DisplayedValue `json:"assignment_group"`
	ApplicationCategory    DisplayedValue `json:"u_application_category"`
	Parent                 DisplayedValue `json:"parent"`
	IncCount               string         `json:"incCount"`
	ChgCount               string         `json:"ChgCount"`
	PrbCount               string         `json:"PrbCount"`
	Entity                 string         `json:"entity"`
}

// Group is the SNAPI version of a group record
type Group struct {
	DisplayedValue
	DistributionList string `json:"u_distribution_list"`
	ADSecurityGroup  string `json:"u_ad_security_group"`
	Manager          User   `json:"manager"`
	Director         User   `json:"u_director"`
	VP               User   `json:"u_vp"`
}

// User is the SNAPI version of a user record
type User struct {
	DisplayedValue
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

// ChangeRequest is the old version of a ServiceNow Change Request record. This structure is deprecated and should be updated
type ChangeRequest struct {
	Type             string `json:"type"`
	Category         string `json:"category"`
	Risk             string `json:"risk"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	AssignmentGroup  string `json:"assignment_group"`
	CmdbCi           string `json:"cmdb_ci"`
	UVersion         string `json:"u_version"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	BackoutPlan      string `json:"backout_plan"`
	ProductionSystem string `json:"production_system"`
}

// DisplayedValue is a type that containes a value and string for human readable output
type DisplayedValue struct {
	Value        string `json:"value"`
	DisplayValue string `json:"displayValue"`
}

type responseEnvelope struct {
	Result []Application `json:"result"`
	Error  *struct {
		Message string `json:"Message"`
		Detail  string `json:"Detail"`
	} `json:"error"`
	Status string `json:"status"`
}
