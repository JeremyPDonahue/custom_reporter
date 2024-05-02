# custom_reporter

This project has been uploaded to GitHub for demonstration purposes only. This is a copy of a non-sensitive project that I built to be used internally at my company.
For that purpose I am removing any references specific to any internal links etc.


## Setup:
- Ensure you are running go1.19
- Make sure you have setup your AWS profiles using the [Setup Script]() (If you are getting keychain errors) 
- Set the following ENV variables:

    - `GITLAB_API=https://gitlab.organization.com/api/v4/projects/`
    
    - `GITLAB_API_TOKEN` (API Token for OLD Gitlab)

--------------
## Usage:
This tool currently generates one of two reports in CSV format.
1. Namespace Report
2. API Audit Report

To run this from the root of the ns-report-gen directory and produce a Namespace Report run:
```
go run main.go -nsr
```
Or to run the API Audit Report:
```
go run main.go -api
```
The CSV file will appear along side the main.go file when complete.
```
namespace_report_YYYY_MM_DD.csv or apiAudit_YYYY_MM_DD.csv
```

## Make it a command line tool:

**On Mac:**
```
GOOS=darwin go build custom_reporter .
```
You can move the resulting binary to /usr/local/bin. This way you can just type `custom_reporter` with `-nsr` or `-api` wherever you are on your system and it will generate the desired report in your current directory.

From within the root directory of the `ns-report-gen` run:
```
cp custom_reporter /usr/local/bin/custom_reporter
```

Once it's there you can alias it if you like by adding something like this to your rc file (e.g. `zshrc` etc.)
```
Example:
alias reporter="/usr/local/bin/custom_reporter"
```
That way you can just call it from the command line like this:
```
reporter -nsr or reporter -api
```
---

**On Windows:**

```
GOOS=windows go build custom_reporter .
```
If you are on a Windows system you can store the resulting `custom_reporter` binary in any directory
that is included in your system `PATH`. 

_Not Working?_
------------
- Do you need to be on a VPN?
- [Do you have Go installed?](https://go.dev/doc/install)
- Are you set up to query the clusters? (do you have your kubeconfigs setup under $HOME/.kube?)