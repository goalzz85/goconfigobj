*Package goconfigobj* is Golang Version for configobj in Python.
Python Version: [https://github.com/DiffSK/configobj](https://github.com/DiffSK/configobj)

Just Parse UTF-8 no BOM text, not parse list type, be careful.

```Golang
package main

import "fmt"
import "github.com/kaixinmao/goconfigobj"
import "strings"

func main() {
	demoStr := `
    outside_key=val
    [test]
       [[demofff]]
            multi_line="""
            multi
            line
            data!
            """
            key=value
    [test2]
      test2key=value
    `

	reader := strings.NewReader(demoStr)
	configObj := goconfigobj.NewConfigObj(reader)
	fmt.Println(configObj.Value("outside_key"))
	fmt.Println(configObj.Section("test").Section("demofff").Value("multi_line"))
}
```