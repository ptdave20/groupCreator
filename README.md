# Group Creator
## What does it do?
This application will take a CSV file and create a series of groups in a GSuite Apps Domain. The CSV must be in the
following format (header included):

```
E-Mail Address, Group Title, Member Of
allstaff@example.org, All Staff,
admingroup@example.org, Administration, allstaff@example.org
buildinga@example.org, Building A, allstaff@example.org
```

After, it will modify the settings of the group to match the requested. Change it to how you want and re-compile.

## Requirement
You must provide a config.json in the directory you are executing from. You can create one from the Google Developer Console
found at: https://console.developers.google.com

Grant it permissions to Admin Directory V1, download the credentials file to the project directory and name it config.json

## Building
In order to build this, you need to Go development tools. They are available at https://golang.org/dl/.
Follow the instructions for the appropriate operating systemyou are on. After, you can download and install by running:

```
go get "github.com/ptdave20/groupCreator"
go install "github.com/ptdave20/groupCreator"
```

This will install the binary to your default go binary directory, which is normally your <user directory>/go/bin

## Running
Execute the application with the csv files as arguments, you can run through multiple.
During the first run, it will prompt to authenticate it, and should automatically open the browser for you.

You can check the reports in the admin page for GSuite to determine any changes.