package static

import _ "embed"

//go:embed index.html
var FileIndex string

//go:embed home.html
var FileHome string

//go:embed state.html
var FileState string

// goFembed index_without_header.html
var FileIndexWithoutHeader string

//go:embed view_table.html
var FileViewTable string

//go:embed view_chart.html
var FileViewChart string

//go:embed style.css
var FileStyleCss string

//go:embed private_policy.html
var FilePrivatePolicy string

//go:embed user_agreement.html
var FileUserAgreement string

//go:embed single_chart.js
var FileSingleChart string
