package messages

import (
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/gobwas/vk"
)

var t *template.Template

func init() {
	t = template.New("index.html")
	t.Funcs(template.FuncMap(map[string]interface{}{
		"toDate": func(unix int64) string {
			t := time.Unix(unix, 0)
			return strings.Replace(t.Format(time.RFC3339), "T", " ", 1)
		},
		"homePage": func(user vk.User) string {
			s := user.Domain
			if s == "" {
				s = strconv.Itoa(user.ID)
			}
			return "https://vk.com/" + s
		},
	}))
	t = template.Must(t.Parse(index))
}

const index = `
<!DOCTYPE html>
<html>
<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.2/css/bootstrap.min.css" integrity="sha384-PsH8R72JQ3SOdhVi3uxftmaW6Vc51MKb0q5P2rRUpPvrszuE4W1povHYgTpBfshb" crossorigin="anonymous">
	<style>
	body {
		font-size: 13px;
	}
	.table td, .table th {
		vertical-align: middle;
	}
	.table-secondary>td {
		background-color: #fbfbfb;
	}
	.table .thead-dark th {
		background-color: #4a76a8;
		border-color: #4a76a8;
	}
	.table .vk-date {
		white-space: nowrap;
		color: #939393;
	}
	.table td.vk-user {
		color: #42648b;
	}
	.table .vk-user a {
		color: #42648b;
	}
	.table td.vk-user_me {
	}
	</style>
</head>
<body>
	<table class="table">
		<thead class="thead-dark">
			<tr>
				<th scope="col">Date</th>
				<th scope="col">From</th>
				<th scope="col">Body</th>
			</tr>
		</thead>
		<tbody>
		{{ range .Messages }}
			<tr class="{{ if ne .FromID $.User.ID }}table-secondary{{ end }}" id="{{ .ID }}">
				<td class="vk-date">{{ toDate .Date }}</td>
			{{ if eq .FromID $.User.ID }}
				<td class="vk-user font-weight-bold">
					<a target="_blank" title="{{ $.User.FirstName }} {{ $.User.LastName }}" href="{{ homePage $.User }}">{{ $.User.FirstName }}</a>
				</td>
				<td>{{ .Body }}</td>
			{{ else }}
				<td class="vk-user vk-user_me font-weight-bold">
					<a target="_blank" title="{{ $.Me.FirstName }} {{ $.Me.LastName }}" href="{{ homePage $.Me }}">{{ $.Me.FirstName }}</a>
				</td>
				<td class="vk-me">{{ .Body }}</td>
			{{ end }}
			</tr>
		{{ end }}
		</tbody>
	</table>
</body>
</html>
`
