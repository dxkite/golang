# HelloWorld

用来测试go-get的可行性

go-get协议返回的网页规则如下：


```html
<!DOCTYPE html>
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="go-import" content="{{ $:host }} git {{ $:repo }}">
    <meta name="go-source"
        content="{{ $:host }} {{ $:repo }} {{ $:repo }}/tree/master{/dir} {{ $:repo }}/blob/master{/dir}/{file}#L{line}">
    <meta http-equiv="refresh" content="0; url={{ $:doc }}">
</head>

<body>
    Nothing to see here; <a href="{{ $:doc }}">move along</a>.
</body>

</html>
```

**注意** 

- `{{ $:host }}` -> 当前运行的域名
- `{{ $:repo }}` -> Github Repo
- `{{ $:doc  }}` -> 文档地址

---

主要生效的数据是 meta 数据，其他的内容都是仿制 golang 官方。
