package templates

// Index template for the Vue SPA shell
const Index = `[[define "index"]]<!doctype html>
<html lang="en">
[[template "head" . ]]
<body>
<div id="app"></div>
</body>
</html>[[end]]`

// Head items
const Head = `[[define "head"]]<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="description" content="[[ .desc ]]" />[[if .nsfw]]
<meta name="rating" content="adult" />
<meta name="rating" content="RTA-5042-1996-1400-1577-RTA" />
[[end]]
<title>[[ .title ]]</title>
<link rel="stylesheet" href="/assets/prim/[[ .primcss ]]" />
<link rel="stylesheet" href="/assets/styles/[[ .style ]]" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
<script src="/assets/prim/[[ .primjs ]]"></script>
[[template "primconfig" . ]][[template "headinclude" . ]]
</head>[[end]]`

// PrimConfig injects the Vue app configuration into window.primConfig.
// The value is a JSON-encoded template.JS produced by the controller,
// which avoids XSS from DB fields containing quotes or special characters.
const PrimConfig = `[[define "primconfig"]]<script>
window.primConfig=[[ .config ]];
</script>[[end]]`

// Empty include for deployment-specific head overrides
const HeadInclude = `[[define "headinclude"]][[end]]`
