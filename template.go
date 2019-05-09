package main

const infoPageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Information</title>
</head>
<body>
<h1>Available Links</h1>
<ul>
<li><a href="/gr">Gibson Ridge</a></li>
</ul>
<br/>
<h1>Currently tracked</h1>
<ul>
{{ range .Spotters }}
<li>{{.First}} {{.Last}}</li>
{{ end }}
</ul>
</body>
</html>
`
