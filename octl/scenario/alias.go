package scenario

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var aliasMap map[string][]string

func parseAliasFile(filename string) error {
	aliasMap = map[string][]string{}
	aliasFile, err := os.Open(filename)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}
	defer aliasFile.Close()

	raw, err := io.ReadAll(aliasFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, &aliasMap)
	if err != nil {
		return err 
	}

	err = expandSelf()
	if err != nil {
		return err 
	}
	return nil 
}

func expandSelf() error {
	expandMap := map[string][]string{}

	var expand func (alias string, seen map[string]struct{}) ([]string, error)
	expand = func (alias string, seen map[string]struct{}) ([]string, error) {
		// avoid looping
		if _, hasSeen := seen[alias]; hasSeen {
			return nil, fmt.Errorf("alias loop reference: %s", alias)
		}
		seen[alias] = struct{}{}
		defer delete(seen, alias)

		// quick return
		if res, found := expandMap[alias]; found {
			return res, nil 
		}

		// dfs expand
		if res, found := aliasMap[alias]; found {
			output := []string{}
			for _, item := range res {
				if len(item) == 0 {
					continue 
				}
				if item[0] == '@' {
					exp, err := expand(item[1:], seen)
					if err != nil {
						return nil, err 
					}
					output = append(output, exp...)
				} else {
					output = append(output, item)
				}
			}
			expandMap[alias] = output
			return output, nil
		} 
		return nil, fmt.Errorf("invalid alias name: %s", alias)
	}

	for key, namelist := range aliasMap {
		if _, hasExpanded := expandMap[key]; hasExpanded {
			continue
		}
		seen := map[string]struct{}{}
		distinctNames := map[string]struct{}{}
		for _, name := range namelist {
			if len(name) == 0 {
				continue 
			}
			if name[0] == '@' {
				expandlist, err := expand(name[1:], seen)
				if err != nil {
					return err 
				}
				for _, expandName := range expandlist {
					distinctNames[expandName] = struct{}{}
				}
			} else {
				distinctNames[name] = struct{}{}
			}
		}

		outputNameList := []string{}
		for name := range distinctNames {
			outputNameList = append(outputNameList, name)
		}
		expandMap[key] = outputNameList
	}

	aliasMap = expandMap
	return nil
}

func expandAlias(input []string) ([]string, error) {
	output := []string{}
	distinctNames := map[string]struct{}{}
	for _, name := range input {
		if len(name) == 0 {
			continue 
		}
		if name[0] == '@' {
			if namelist, found := aliasMap[name[1:]]; found {
				for _, expandName := range namelist {
					distinctNames[expandName] = struct{}{}
				}
			} else {
				return nil, fmt.Errorf("invalid Alias: %s", name[1:])
			}
		} else {
			distinctNames[name] = struct{}{}
		}
	}
	for distinctName := range distinctNames {
		output = append(output, distinctName)
	}
	return output, nil
}
