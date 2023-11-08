package scenario

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

const scenTemplete = `# auto generate by octl
name: %s
description: "A simple description about scenarios %s..."
applications:
`

const aliasTemplete = `# auto generate by octl
group1:
- 'node3'
- 'node4'
group2:
- '@group1'
- 'node1'
- 'node2'
`

const appTemplate = `-
  name: "%s"
  scriptpath: "%s/scripts/"
  sourcepath: "%s/src/"
  description: "A simple description about application %s..."
  nodes:
    - 'node1'
    - '@group1'
  script:
    -
      target: prepare
      file: "prepare.sh"
      order: 1
    -
      target: start
      file: "start.sh"
      order: 1
    -
      target: stop
      file: "stop.sh"
      order: 1
    -
      target: purge
      file: "purge.sh"
      order: 1
    -
      target: user-defined
      file: "user-defined.sh"
      order: 1
`

const readmeTemplate = `- README of %s -
`

const prepareTemplate = `#!/bin/bash
# auto generate by octl
echo "prepare" >> $OCTOPODA_OUTPUT
`

const startTemplate = `#!/bin/bash
# auto generate by octl
echo "start" >> $OCTOPODA_OUTPUT
`

const stopTemplate = `#!/bin/bash
# auto generate by octl
echo "stop" >> $OCTOPODA_OUTPUT
`

const purgeTemplate = `#!/bin/bash
# auto generate by octl
echo "purge" >> $OCTOPODA_OUTPUT
`

const userdefTemplate = `#!/bin/bash
# auto generate by octl
echo "user-defined target" >> $OCTOPODA_OUTPUT
`

const resultTemplate = `Success to create config folder of %s!
    To deploy the scenario: octl apply %s prepare -m 'your message'.
    You can also host this scenario (%s/%s/%s) config with git service.`

func Create(scenario string, apps ...string) (string, *errs.OctlError) {
	output.PrintInfoln("Generating project...")
	err := os.Mkdir(scenario, os.ModePerm)
	if err != nil {
		emsg := fmt.Sprintf("Can not create folder %s: %s\n", scenario, err.Error())
		output.PrintFatalf(emsg, err)
		return emsg, errs.New(errs.OctlFileOperationError, emsg)
	}
	basePath, err := filepath.Abs(scenario)
	if err != nil {
		emsg := fmt.Sprintf("Can not calc basePath %s: %s\n", scenario, err.Error())
		output.PrintFatalf(emsg, err)
		return emsg, errs.New(errs.OctlFileOperationError, emsg)
	}
	configFile := fmt.Sprintf(scenTemplete, scenario, scenario)
	for _, app := range apps {
		configFile += fmt.Sprintf(appTemplate, app, app, app, app)
		appfolder := basePath + "/" + app
		err = os.Mkdir(appfolder, os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create folder %s: %s\n", appfolder, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		srcfolder := appfolder + "/src"
		err = os.Mkdir(srcfolder, os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create folder %s: %s\n", srcfolder, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		scriptfolder := appfolder + "/scripts"
		err = os.Mkdir(scriptfolder, os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create folder %s: %s\n", scriptfolder, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		readmefile := srcfolder + "/README"
		err = os.WriteFile(readmefile, []byte(fmt.Sprintf(readmeTemplate, app)), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", readmefile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		preparefile := scriptfolder + "/prepare.sh"
		err = os.WriteFile(preparefile, []byte(prepareTemplate), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", preparefile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		startfile := scriptfolder + "/start.sh"
		err = os.WriteFile(startfile, []byte(startTemplate), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", startfile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		stopfile := scriptfolder + "/stop.sh"
		err = os.WriteFile(stopfile, []byte(stopTemplate), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", stopfile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		purgefile := scriptfolder + "/purge.sh"
		err = os.WriteFile(purgefile, []byte(purgeTemplate), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", purgefile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
		userdeffile := scriptfolder + "/user-defined.sh"
		err = os.WriteFile(userdeffile, []byte(userdefTemplate), os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("Can not create file %s: %s\n", userdeffile, err.Error())
			output.PrintFatalf(emsg, err)
			return emsg, errs.New(errs.OctlFileOperationError, emsg)
		}
	}
	deployfile := basePath + "/deployment.yaml"
	err = os.WriteFile(deployfile, []byte(configFile), os.ModePerm)
	if err != nil {
		emsg := fmt.Sprintf("Can not create deployment file %s: %s\n", deployfile, err.Error())
		output.PrintFatalf(emsg, err)
		return emsg, errs.New(errs.OctlFileOperationError, emsg)
	}

	aliasfile := basePath + "/alias.yaml"
	err = os.WriteFile(aliasfile, []byte(aliasTemplete), os.ModePerm)
	if err != nil {
		emsg := fmt.Sprintf("Can not create alias file %s: %s\n", aliasfile, err.Error())
		output.PrintFatalf(emsg, err)
		return emsg, errs.New(errs.OctlFileOperationError, emsg)
	}

	result := fmt.Sprintf(resultTemplate, scenario, scenario,
		config.GlobalConfig.Gitinfo.ServeUrl, config.GlobalConfig.Gitinfo.Username, scenario)
	output.PrintInfoln(result)
	return result, nil
}
