# custom_reporter

>Ensure you are running go version go1.19.4

TODO: Change the name of this tool to something like custom_reporter,
because it does more than just produce a NS report.

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
- Are you on the VPN?
- [Do you have Go installed?](https://go.dev/doc/install)
- Have you logged into AWS? (aws-okta login)
- [Are you set up to query the clusters?](https://gitlab.nordstrom.com/dev-compute/cluster-registry/-/merge_requests/42/diffs)
- What are the 'user' fields set to in your ~/.kube/config file (legacy clusters only) --- e.g. barcelona_sudo