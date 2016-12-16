variable "this_is_a_var" {
  default = "Default value"
}

variable "this_is_a_number" {
  default = 2
}

variable "this_is_a_bool" {
  default = "false"
}

variable "this_is_a_list" {
  type    = "list"
  default = ["value1", "value1"]
}

variable "this_is_a_map" {
  type = "map"
  default = {
    key1 = "val1"
    key2 = "val2"
  }
}
