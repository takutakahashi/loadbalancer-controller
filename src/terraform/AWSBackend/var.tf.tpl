variable region {}
variable vpc_id {}

variable lb_name {}
variable lb_internal {}
variable lb_type {}
variable lb_protection {}
variable lb_subnets {
  type = list(string)
}
variable lb_tags {
  type = map
}

{{- range $i, $l := .B.Spec.Listeners }}
{{- $name := printf "%s_%d" $l.Protocol $l.Port }}
{{- $tg := $l.DefaultAction.TargetGroup }}

variable {{ $name }}_listener {
  type = map(string)
}

variable {{ $name }}_target_group {
  type = map(string)
}

variable {{ $name }}_targets {
  type = list(map(string))
}
{{- end }}
