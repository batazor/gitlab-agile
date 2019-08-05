workflow "New workflow" {
  resolves = ["Golang Action"]
  on = "push"
}

action "Golang Action" {
  uses = "cedrickring/golang-action@1.3.0"
  secrets = ["GITHUB_TOKEN"]
}
