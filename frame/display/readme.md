Display
---

The biggest gotcha regarding this module is the way the template files are read: the module tries to load every file first from the /template, then
if it can not be found, from /modules.

Examples:  
":module/:action" => /template/:module/:action.tpl,  /modules/:module/tpl/:action.tpl  
"skeleton/GetSingle" => /template/skeleton/GetSingle.tpl,  /modules/skeleton/tpl/GetSingle.tpl

If multiple files are specified:  
```go
// ...
err := display.New(ctx).Do([]string{"articles/GetSingle", "skeleton/GetSingle"})
// ...
```
The first found will be executed, so put the more specific ones first.
