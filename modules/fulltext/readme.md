fulltext
---

Fulltext module helps to add fulltext search and indexing capability to noCrud.  

To generate a fulltext index for "articles":

```js
{
	"Hooks": {
		"articlesInserted": [ [ "fulltext", "SaveFulltext" ] ],
		"articlesUpdated": [ [ "fulltext", "SaveFulltext" ] ]
	}
}
```

To add a special kind of field "search" to filters:

```js
{
	"Hooks": {
		"ProcessMap": [ "fulltext" ]
	}
}
```