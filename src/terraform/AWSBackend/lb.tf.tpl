resource "aws_lb" "loadbalancer" {
  name               = var.lb_name
  internal           = var.lb_internal
  load_balancer_type = var.lb_type
  subnets            = var.lb_subnets

  enable_deletion_protection = var.lb_protection

  tags = var.lb_tags
}

{{- range $i, $l := .B.Spec.Listeners }}
{{- $name := printf "%s_%d" $l.Protocol $l.Port }}
{{- $tg := $l.DefaultAction.TargetGroup }}

resource "aws_lb_listener" "{{ $name }}_listener" {
  load_balancer_arn = aws_lb.loadbalancer.arn
  port              = var.{{ $name }}_listener.port
  protocol          = var.{{ $name }}_listener.protocol
  default_action {
    type             = var.{{ $name }}_listener.action_type
    target_group_arn = aws_lb_target_group.{{ $name }}_tg.arn
  }
}

resource "aws_lb_target_group" "{{ $name }}_tg" {
  port        = var.{{ $name }}_target_group.port
  protocol    = var.{{ $name }}_target_group.protocol
  target_type = var.{{ $name }}_target_group.type
  vpc_id      = var.vpc_id
}

resource "aws_lb_target_group_attachment" "{{ $name }}_tga" {
  count            = length(var.{{ $name }}_targets)
  target_group_arn = aws_lb_target_group.{{ $name }}_tg.arn
  target_id        = var.{{ $name }}_targets[count.index].destination
  availability_zone = "ap-northeast-1a"
  port             = var.{{ $name }}_targets[count.index].port
}

{{- end }}