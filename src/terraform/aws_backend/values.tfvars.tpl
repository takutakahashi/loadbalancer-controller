vpc_id = "vpc-000000000000"

lb_name = "test"
lb_internal = true
lb_type = "network"
lb_protection = false
lb_subnets = ["subnet-00000000000"]
lb_tags = {aaa = "iii"}

listener_port = 443
listener_protocol = "TCP"
listener_action_type = "forward"

target_group_port = 443
target_group_protocol = "TCP"
target_group_type = "ip"

target_port = 443
targets = ["10.0.0.10"]
