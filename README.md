# Ulule Api Client

Api Client for Ulule's API ([developers.ulule.com](http://developers.ulule.com)), written in Go. 

```go
import (
	ulule "github.com/aduermael/ulule-api-client"
)

func main() {
	client := ulule.New(username, userid, apikey)
	projects, err := client.GetProjects(ulule.ProjectFilterSupported)
}
```