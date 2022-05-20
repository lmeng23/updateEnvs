本程序的作用是将源青龙面板的JD_COOKIE环境变量（列表形式）更新到目标青龙面板的JD_COOKIE环境变量中（列表转字符串，以&拼接）。

有以下6个变量需要填写。

```js
sourceUrl = ""                           // 源青龙面板的URL
sourceClientID = ""                      // 源青龙面板的client_id
sourceClientSecret = ""                  // 源青龙面板的client_secret

distUrl = ""                             // 目标青龙面板的URL
distClientID = ""                        // 目标青龙面板的client_id
distClientSecret = ""                    // 目标青龙面板的client_secret
```
