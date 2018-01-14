# 管理画面及びPHPバッチのCRUD({{.Datetime}} 時点)

#### ※ツール（ https://github.com/sky0621/go-crudmaker-by-regexp ）による自動生成

#### ・「controller」層、「batch」層から直接「service」層を呼んでいるケースのみ想定

#### ・CRUDの判定については「service」層のメソッドが「get〜〜」なら「READのR」、「insert」を含むなら「CREATEのC」といった恣意的なレベル

{{range .TargetGroups}}## {{.Name}}

{{range .Services}}

{{.Name}}

{{range .Tables}}

{{.Name}}

{{range .CRUDs}}

{{.Name}}

{{end}}
{{end}}
{{end}}
{{end}}