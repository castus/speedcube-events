package printer

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(obj any) {
	fmt.Print(Stringify(obj))
}

func Stringify(obj any) string {
	j, _ := json.MarshalIndent(obj, "", "    ")
	return fmt.Sprintf("%s\n", j)
}
