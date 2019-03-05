# NSGExport

Export just an NSG from an Azure resource group

# Requirements

* An AAD application account with subscription read access. You need this in order to authenticate against the API


# Sample usage

First create a `conf.json` file that lives in the same directoy as the nsgexport binary. From their plug in the values you'll be using.

`nsgexport > template.json`

Or if you want pretty formatting:

`nsgexport | jq . > template.json`

