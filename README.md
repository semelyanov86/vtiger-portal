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

You can also download executable file and run it on server.

## API Documentation
Swagger documentation is available on `/swagger/` endpoint
Also there are yaml swagger file in `documents` folder, available for import in Insomnia.
To import whole insomnia collection, there are Insomnia.yaml file in `documents` folder.

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
Note, that Vtiger does not provide an endpoint for storing images. If you want to allow users change images, you need to add following code in file `include/Webservices/Revise.php`, before line `$entity = $handler->revise($element);`:
```php
if ($element['imagecontent'] && $element['imagename'] && $element['imagetype']) {
            $decodedImageContent = base64_decode($element['imagecontent']);
            $tempFileName = tempnam(sys_get_temp_dir(), 'image_');
            file_put_contents($tempFileName, $decodedImageContent);
            $_FILES['uploaded_image'] = [
                'name' => $element['imagename'], // Provide the desired file name
                'type' => $element['imagetype'], // Provide the correct file type
                'tmp_name' => $tempFileName,
                'error' => 0,
                'size' => filesize($tempFileName)
            ];
            /** @var Contacts $contactModule */
            $contactModule = new $entityName;
            if (method_exists($contactModule, "insertIntoAttachment")) {
                $contactModule->insertIntoAttachment($idList[1], $entityName);
            }
        }
```
Now, you can pass imagename, imagetype and imagecontent params

Note, in old vtiger versions, we receive always empty attachment ids. In this case, you can not download documents via API. To fix this, add methods from vtiger 7.5 in `include/Webservices/DataTransform.php` and `vtlib/Vtiger/Functions.php`. sanitizeFileFieldsForIds and getAttachmentIds

Also for add_related endpoint currently we support GET, not POST request.

If vtiger does not provide file contents in response, fix file `include/Webservices/Utils.php`, function `vtws_filedetails`, change filenamewithpath variable:
```php
        $storedName = $fileData['storedname'];
        if (!$storedName) {
            $storedName = $fileName;
        }
        $filenamewithpath = $fileData['path'].$fileData['attachmentsid'].'_'.$storedName;
```

Currently for delete operation in vtiger we support GET request. To change REST api in Vtiger for this operation, change `vtiger_ws_operation` table, for delete operation name type from POST to GET.

## Adding new custom field to module
What if you created new custom field in Vtiger module and want to add it in Portal? Because we using golang type system in portal, you need to register it in our domain system.
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

## Support for custom modules
System has a support for custom modules. For example. what can you do if you have custom module in Vtiger with name SPPayments and you want to display it in customer portal? To enabme this feature, you need to insert configuration in portal.yaml:
```yaml
    customModules:
      SPPayments:
        - pay_no
        - pay_no
        - pay_type
```
Where SPPayments is the name of module. Below module name you need to insert minimum 2 fields. First field is default sorting field. Second field and more is fields, which support live search

## User Settings Route

This application supports managing of user settings, which is available under `user/settings` endpoint. As a result you will get list of boolean fields and their values. These fields are dynamic, you can set them in configuration, under `userSettingsFields` part.

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

## Features
- Contact registration via temporary code
- CRUD operations for helpdesk module
- Get all project, project tasks, invoices and sales orders. Change statuses of this modules.
- Support for custom modules
- Make a payment in invoice and sales order.
- Storing notifications from vtiger crm

## Deployment

To deploy this project run

```bash
  make production/deploy/api
```

Note, this script running as a service. So, you need to create tunnel between apache and this script in order to use https

Example apache configuration you can find in remote/production/apache.conf

To see the status of application run command:
```bash
sudo systemctl status portal
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Feedback
If you have any feedback, please reach out to us at info@itvolga.com


## Tech Stack

**Client:** React, Zod, Zustand, Bootstrap

**Server:** Golang

## Customizing the portal
This portal is built using Golang. You are welcome to fork and modify this repository for your own needs! Do not attempt to customize the files inside the bin folder, as they will be automatically overwritten during rebuilt. If you need help customizing this portal beyond what is currently available, we recommend Center of information Technologies (info@itvolga.com) as an experienced third party developer.