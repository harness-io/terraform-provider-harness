resource "harness_encrypted_text" "test" {
	name = "git_secret"
	value = "foo"
}

resource "harness_git_connector" "test" {
	name = "test-git-connector"
	url = "https://github.com/micahlmartin/harness-demo"
	branch = "main"
	generate_webhook_url = true
	password_secret_id = harness_encrypted_text.test.id
	url_type = "REPO"
	username = "someuser"

	commit_details {
		author_email_id = "user@example.com"
		author_name = "test user"
		message = "custom commit message here"
	}

	delegate_selectors = ["primary"]
}	
