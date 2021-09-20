# variable definitions for all builds
variable "home"                    { default = "{{ env `HOME` }}" }
variable "user"                    { default = "{{ env `USER` }}" }
