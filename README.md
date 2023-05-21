# Vtiger Customer Portal [Backend GO Application] ![GO][go-badge]

[go-badge]: https://img.shields.io/github/go-mod/go-version/p12s/furniture-store?style=plastic
[go-url]: https://github.com/semelyanov86/vtiger-portal/blob/main/go.mod

Learn More about Vtiger [here](https://vtiger.com)

## Build & Run (Locally)
### Prerequisites
- go 1.20
- [staticcheck](https://staticcheck.io) (<i>optional</i>, for code static checking)
- [swag](https://github.com/swaggo/swag) (<i>optional</i>, used to re-generate swagger documentation)

To fill database, first install migrate tool:
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 $GOPATH/bin/migrate
```

Before you continue, please check that it’s available and working on your machine by trying to execute the migrate binary with the -version flag. It should output the current version number similar to this:
```bash
$ migrate -version
4.14.1
```

To run migration, execute following command:
```bash
make migrate
```

Copy file .envrc.example to .envrc

Use `make run` to build&run project, `make lint` to check code with linter.

## Vtiger Setup
By default you do not need to do any setup for Vtiger crm. Make sure, you have latest vtiger 7.5 version installed.
Also make sure that you have FilesRetrieve webservice function installed, in order you want to receive files and images from vtiger. If not, do the following:
1. Download [FilesRetrieve](https://code.vtiger.com/vtiger/vtigercrm/-/blob/master/include/Webservices/FileRetrieve.php) file and paste it in appropriate folder.
2. Register new webservice by executing sql query:
```sql
INSERT INTO `vtiger_ws_operation`(
    `operationid`,
    `name`,
    `handler_path`,
    `handler_method`,
    `type`,
    `prelogin`
) VALUES (
    '38',
    'files_retrieve',
    'include/Webservices/FileRetrieve.php',
    'vtws_file_retrieve',
    'GET',
    '0'
 );
```
3. Register parameters of webservice by executing sql query:
```sql
INSERT INTO `vtiger_ws_operation_parameters`(
    `operationid`,
    `name`,
    `type`,
    `sequence`
) VALUES (
    '38',
    'id',
    'string',
    '1'
 );
```

## Adding new custom field to module
What if you created new custom field in Vtiger module and want to add it in Portal? Because we usong golang type system in portal, you need to register it in our domain system.
For example, you created field 'cf_543' in HelpDesk module. Here is three steps, how you can register this field in portal:
1. Open file `internal/domain/help_desk.go` and find type HelpDesk struct. This is block, where we store all fields for module. At the end of this struct you can add following line:
```go
Cf543 string `json:"cf_543"`
```
2. We also need to tell a system. how to decode vtiger field to our struct. For this purposes open the same file and find function ConvertMapToHelpDesk. There you need to add new case:
```go
case "cf_543":
			helpDesk.Cf543 = v.(string)
```
3. Recompile a project by running `make build` command

## Command line arguments

You can run executable script with following arguments:

* `version` - Display script version and exit.

## Configuration file
There is 2 configuration example files.
`.envrc` - for storing environment variables
`portal.yaml` - for storing app configuration

To create them, use .envrc.example and mail.yaml files.
Put your config `portal.yaml` file in `~/.config` directory

## Running Tests

To run tests, run the following command

```bash
  make audit
```


## Deployment

To deploy this project run

```bash
  production/deploy/api