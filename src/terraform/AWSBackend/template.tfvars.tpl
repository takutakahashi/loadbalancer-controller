{{ $name := .Name }}
region = "{{ .B.Spec.Region }}"
vpc_id = "{{ .B.Spec.VPC.ID }}"

lb_name = "{{ $name }}"
lb_internal = {{ .B.Spec.Internal }}
lb_type = "{{ .B.Spec.Type }}"
lb_protection = {{ .ServiceIn }}
lb_subnets = [
{{- range $i, $s := .B.Spec.Subnets }}
    "{{ $s.ID }}",
{{- end }}
    ]
lb_tags ={
{{- range $k, $v := .B.Labels }}
    {{ $k }} = "{{ $v }}"
{{- end }}
}


{{- range $i, $l := .B.Spec.Listeners }}
{{- $name := printf "%s_%d" $l.Protocol $l.Port }}
{{- $tg := $l.DefaultAction.TargetGroup }}

{{ $name }}_listener = {
  port = {{ $l.Port }}
  protocol = "{{ $l.Protocol }}"
  action_type = "{{ $l.DefaultAction.Type }}"
}

{{ $name }}_target_group = {
  port = {{ $tg.Port }}
  protocol = "{{ $tg.Protocol }}"
  type = "{{ $tg.TargetType }}"

}

{{$name}}_targets = [
{{- range $i, $t := $tg.Targets }}
  {
{{- if eq $tg.TargetType "ip" }}
    destination = "{{ $t.Destination.IP }}"
{{- else }}
    destination = "{{ $t.Destination.InstanceID }}"
{{- end }}
    port = "{{ $t.Port }}"
  },
{{- end }}
]
{{- end }}
