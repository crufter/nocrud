Templates
---

### Public and private folders

This folder can contain two folders:
- /public (public templates)
- /private (private templates)

Public folders are meant to be never written once published.  
A host's private template reside in the folder /private/:hostname (eg: /private/example.com)

### Routing

If a request path identifies a resource-verb combination, then the loaded template file will be the "resource/verb", or "modulename/verb",
depending on which can be found. If it does not

Examples:

- **/users/:user-id** (or literally: /users/UIBQHm-Q6RBEAAAB) will be mapped to **"users/GetSingle.tpl"**, or **"skeleton/GetSingle.tpl"**
- **/contact-us** will be mapped to **contact-us.tpl**, in case no "contact-us" resource is defined.

### Meta access pathes

In the templates, if you want to load resources (js, css, etc files), you can use the following meta pathes:
- /template		Is the folder of your current template (so you dont have to care about the template name, neither if it is public or private).
- /uploads		Is the uploads folder of your site.

### The template engine

The templates use the language of the go template package http://golang.org/pkg/html/template  
For builtin functions see /frame/display/builtins.go

### Localized data in templates

You can load localized text by using pseudovariables {{.loc.filename.varname}}  
The system recognizes this, and loads the used files.
The variable **{{.loc.filename.varname}}** will result in the **/template/loc/filename.en** file to be loaded, where "en" is the actual language of the user.

The language of the user currently comes from two sources:
- It can either be defined explicitly in the database, under the field "languages",
- Or it is inherited from the Accept-Language header of the browser.

### Template hooks

You can call a hook by:
```
{{$h := hook "HookName"}}
{{$h.Fire}}
```

See an example in the module **fkid**.  
Local template variables (beginning with $ and not .) of the hook caller are not available in the templates displayed by the hook.  

#### For more information about the display layer, see /frame/display