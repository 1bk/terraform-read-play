variable "image_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive   = false
#   nesting = {
#     val = 12334
#   }
}

variable "other_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive     = true
  sensitiveTwo   = "Testin123"
  nesting = {
    val = 123
  }
}
