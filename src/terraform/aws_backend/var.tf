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

variable listener_port {}
variable listener_protocol {}
variable listener_action_type {}

variable target_group_port {}
variable target_group_protocol {}
variable target_group_type {}

variable target_port {}
variable targets {
  type = list(string)
}
