package helpers

import (
	"log"
	"slices"
)

func IsEmpty(v string) bool {
	ret := false
	if len(v) == 0 {
		ret = true
	}
	return ret
}

func IsEmptyLog(varCheck, varText string, isRequired bool) {
	if IsEmpty(varCheck) {
		if isRequired {
			log.Fatalf("Entry %s are Required %s", varText, ICON_CANCEL_RED)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	} else {
		log.Printf(varText+" %s", ICON_OK)
	}
}

func IsNotListLog(varCheck string, dataMaster []string, varText string, isRequired bool) {
	err := slices.Contains(dataMaster, varCheck)
	if err {
		log.Printf(varText+" %s", ICON_OK)
	} else {
		if isRequired {
			log.Fatalf("Entry %s are Required %s", varText, ICON_CANCEL_RED)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	}
}

func IsError(err error, varText string, isRequired bool) {
	if err != nil {
		if isRequired {
			log.Fatalf("%s are Required %s", varText, ICON_CANCEL_RED)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	} else {
		log.Printf(varText+" %s", ICON_OK)
	}
}
