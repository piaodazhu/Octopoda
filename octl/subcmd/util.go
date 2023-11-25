package subcmd

import "strconv"

func extractArgString(inputlist []string, flagShortName string,
	flagLongName string, defaultValue string) (arg string, output []string) {
	arg, output = defaultValue, nil
	for i := 0; i < len(inputlist); i++ {
		if value, match := isKeyValuePair(flagShortName, inputlist[i]); match {
			arg = value
			continue
		}

		if value, match := isKeyValuePair(flagLongName, inputlist[i]); match {
			arg = value
			continue
		}

		if (inputlist[i] == flagShortName || inputlist[i] == flagLongName) && i+1 < len(inputlist) {
			arg = inputlist[i+1]
			i++
			continue
		}
		output = append(output, inputlist[i])
	}
	return
}

func extractArgInt(inputlist []string, flagShortName string,
	flagLongName string, defaultValue int) (arg int, output []string) {
	arg, output = defaultValue, nil
	for i := 0; i < len(inputlist); i++ {
		if value, match := isKeyValuePair(flagShortName, inputlist[i]); match {
			arg = toInt(value, defaultValue)
			continue
		}

		if value, match := isKeyValuePair(flagLongName, inputlist[i]); match {
			arg = toInt(value, defaultValue)
			continue
		}

		if (inputlist[i] == flagShortName || inputlist[i] == flagLongName) && i+1 < len(inputlist) {
			arg = toInt(inputlist[i+1], defaultValue)
			i++
			continue
		}
		output = append(output, inputlist[i])
	}
	return
}

func extractArgBool(inputlist []string, flagShortName string,
	flagLongName string, defaultValue bool) (arg bool, output []string) {
	arg, output = defaultValue, nil
	for i := 0; i < len(inputlist); i++ {
		if value, match := isKeyValuePair(flagShortName, inputlist[i]); match {
			arg = toBool(value)
			continue
		}

		if value, match := isKeyValuePair(flagLongName, inputlist[i]); match {
			arg = toBool(value)
			continue
		}

		if inputlist[i] == flagShortName || inputlist[i] == flagLongName {
			arg = true
			continue
		}
		output = append(output, inputlist[i])
	}
	return
}

func toBool(s string) bool {
	if s == "true" || s == "True" || s == "" || s == "1" || s == "t" {
		return true
	} 
	return false
}

func toInt(s string, defaultValue int) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return num
}

func isKeyValuePair(key, target string) (string, bool) {
	prefix := key + "="
	if len(prefix) >= len(target) {
		return "", false
	}
	if target[:len(prefix)] == prefix {
		return target[len(prefix):], true
	}
	return "", false
}
