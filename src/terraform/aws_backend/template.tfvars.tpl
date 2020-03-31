{{- $l := index .B.Spec.Listeners 0 }}
{{- $tg := $l.DefaultAction.TargetGroup }}
region = "{{ .B.Spec.Region }}"
vpc_id = "{{ .B.Spec.VPC.ID }}"

lb_name = "{{ .B.Name }}"
lb_internal = {{ .B.Spec.Internal }}
lb_type = "{{ .B.Spec.Type }}"
lb_protection = {{ .ServiceIn }}
lb_subnets = [
{{- range $i, $s := .B.Spec.Subnets }}
    "{{ $s.ID }}"
{{- end }}
    ]
lb_tags ={
{{- range $k, $v := .B.Labels }}
    {{ $k }} = "{{ $v }}"
{{- end }}
}
listener_port = {{ $l.Port }}
listener_protocol = "{{ $l.Protocol }}"
listener_action_type = "{{ $l.DefaultAction.Type }}"

target_group_port = {{ $tg.Port }}
target_group_protocol = "{{ $tg.Protocol }}"
target_group_type = "{{ $tg.TargetType }}"

target_port = {{ $tg.Port }}
targets = [
{{- range $i, $t := $tg.Targets }}
  "{{ $t.Destination }}"
{{- end }}
]