resource "aws_route53_zone" "user_root" {
  name    = "user.takutakahashi.dev"
}

resource "aws_route53_record" "nlb_record" {
  zone_id = aws_route53_zone.user_root.zone_id
  name    = "nlb1.user.takutakahashi.dev"
  type    = "A"
  ttl     = "30"

  records = [
     "10.51.105.4"
  ]
}
