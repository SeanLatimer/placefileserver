package sn

const PlaceFileTemplate = `
Refresh: 1
Threshold: 999
Title: {{.Title}}
Font: 1, 11, 0, "Courier New"
IconFile: 1, 22, 22, 11, 11, "https://www.spotternetwork.org/icon/spotternet_new.png"

{{range .Spotters -}}
{{if .Unix | isReporting -}}
Object: {{.Lat}},{{.Lon}}
Icon: 0,0,000,1,{{if .Unix | isStationary}}6{{else}}2{{end}},"{{.First}} {{.Last}}\n{{.LastReport}} UTC{{if .Unix | isStationary}}\nSTATIONARY{{end}}{{if .Phone}}\nPhone {{.Phone}}{{end}}{{if .Email}}\nEmail: {{.Email}}{{end}}{{if .IM}}\nIM: {{.IM}}{{end}}{{if .Twitter}}\nTwitter: {{.Twitter}}{{end}}{{if .Web}}\nWeb: {{.Web}}{{end}}{{if .Note}}\nNote: {{.Note}}{{end}} "
Text: 15, 10, 1, "{{.First}} {{.Last}}"
End:
{{end}}
{{end}}
`
