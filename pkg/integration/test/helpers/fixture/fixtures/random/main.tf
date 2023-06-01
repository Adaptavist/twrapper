
resource "random_string" "this" {
  length = 10
}
output "this" {
  value = random_string.this.id
}
