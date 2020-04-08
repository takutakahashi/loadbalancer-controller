resource "aws_lb" "loadbalancer" {
  name               = var.lb_name
  internal           = var.lb_internal
  load_balancer_type = var.lb_type
  subnets            = var.lb_subnets

  enable_deletion_protection = var.lb_protection

  tags = var.lb_tags
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_lb.loadbalancer.arn
  port              = var.listener_port
  protocol          = var.listener_protocol
  default_action {
    type             = var.listener_action_type
    target_group_arn = aws_lb_target_group.tg.arn
  }
}

resource "aws_lb_target_group" "tg" {
  name        = var.lb_name
  port        = var.target_group_port
  protocol    = var.target_group_protocol
  target_type = var.target_group_type
  vpc_id      = var.vpc_id
}

resource "aws_lb_target_group_attachment" "tga" {
  count            = length(var.targets)
  target_group_arn = aws_lb_target_group.tg.arn
  target_id        = element(var.targets, count.index).destination
  availability_zone = "ap-northeast-1a"
  port             = element(var.targets, count.index).port
}
