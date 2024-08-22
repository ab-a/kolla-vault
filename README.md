# kolla-vault

Push kolla-ansible passwords in Hashicorp Vault and replace plaintext password by Vault lookups.

```bash
kolla-genpwd
export VAULT_TOKEN=$(vault print token)
go mod init kolla-vault
go get github.com/hashicorp/vault/api
go get gopkg.in/yaml.v2
# export kolla passwords.yml into Hashicorp Vault
go run store_kolla_passwords.go
# replace plaintext password in passwords.yml by Vault lookups
go run replace_passwords.go
```

Snippet of `passwords.yml` lookups:
```yml
nova_database_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_database_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
nova_keystone_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_keystone_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
nova_ssh_key: 'map[private_key:{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_ssh_key/private_key', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }} public_key:{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_ssh_key/public_key', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}]'
nova_api_database_password: '{{ lookup('community.general.hashi_vault', 'secret/data/kolla/default/nova_api_database_password', 'url={{ vault_url }}', token=lookup('env', 'VAULT_TOKEN')) }}'
```


Blog posts:

[Hashicorp Vault and Kolla Ansible, Part I: Integrate Vault secrets in your playbook](https://abayard.com/hashicorp-vault-and-kolla-ansible-part-i-integrate-secrets-in-playbook/)

[Hashicorp Vault and Kolla Ansible, Part II: integration with Gitlab CI](https://abayard.com/hashicorp-vault-and-kolla-ansible-part-ii-integration-with-gitlab-ci/)
