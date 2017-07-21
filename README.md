# NSGExport

Currently exporting an ARM template from Azure will fail when you try to export NSGs. The solution is to hit the API and then create an ARM template based on that data.

Sample usage:

`nsgexport > template.json`

Or if you want pretty formatting:

`nsgexport | jq . > template.json`


# Requirements

* An AAD application account. You need this in order to authenticate against the API

