<html>
<head>
       <title>Index / Upload file</title>
</head>
<body>
Вы можете использовать загруженный в систему файл или загрузить свой.</br>
Текущий файл:
<p>
{{range .}}
{{.}}</br>
{{end}}
</p>
<form enctype="multipart/form-data" action="http://127.0.0.1:9090/upload" method="post">
    <input type="file" name="uploadfile" />
    <!-- <input type="hidden" name="token" value="{{.}}"/> -->
    <input type="submit" value="upload" />
</form>
</body>
</html>
