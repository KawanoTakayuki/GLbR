# [Grouping Logs by Request](https://godoc.org/cloud.google.com/go/logging#hdr-Grouping_Logs_by_Request)

## How to use

* simple

```golang
package main

import (
    "github.com/KawanoTakayuki/glbr"
)

log, err := glbr.NewLogging(context.Background(), "ProjectID", "LogID")
if err != nil {
    panic(err.Error())
}
defer log.Close()

glbr.Debugf(log.Context(), "log")
```

* grouping

```golang
package main

import (
    "github.com/KawanoTakayuki/glbr"
)

log, err := glbr.NewLogging(context.Background(), "ProjectID", "LogID")
if err != nil {
    panic(err.Error())
}
defer log.Close()

glbr.GroupingBy(&http.Request{}, "ParentLogID", func(c context.Context) (int, int64) {
    glbr.Debugf(c, "log")
    return 0,0
})
```
