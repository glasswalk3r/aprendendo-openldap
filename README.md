# Aprendendo OpenLDAP

Este é um repositório foi criado para manter a estrutura de laboratório ensinada
no treinamento online "Aprenda a Gerenciar Diretórios com openLDAP em Linux" do
Marcos Pitanga.

O conteúdo aqui disponível é **não-oficial**, e limita-se a automatizar todos os
passos de criação do ambiente de laboratório utilizando-se do Vagrant e Ansible.
Se você se limitar a copiar e colar, sem entender como as coisas funcionam, o
problema é exclusivamente seu. ☺

## O que é gerenciado

- Configuração de NTP (usando servidores para o Brasil) e rede para todas as VM's
- vinfra: servidor DNS (via Bind)
- master: servidor mestre do OpenLDAP, base local de usuários Linux
- slave: servidor escravo do OpenLDAP
- client: cliente OpenLDAP, usado para autenticar logins na mesma via PAM.

## Pré-requisitos

- Virtualbox versão 6.1 ou maior
- Ansible versão 6.0.0 ou maior
- Vagrant versão 2.2.19 ou maior

O uso de um *virtual environment* do Python é recomendado para o Ansible.

## Como usar

Para evitar problemas devido a resolução de nomes é necessário iniciar primeiro
a VM `vinfra`, para que seja configurada primeiro e as demais VM's possam fazer
uso da mesma ao tentar baixar pacotes RPM:

```
$ vagrant up vinfra
```

Após tudo funcionar como o esperado, o restante poderá ser criado:

```
$ vagrant up
```

Se quiser repetir as configurações para alguma VM em específico, execute
`vagrant provision`. Vide a ajuda *online* para esta opção para maiores
detalhes.

## Referências sobre OpenLDAP

- O livro [OpenLDAP Ultimate](http://www.anahuac.eu/livros-em-cc-by/) do
Anahuac, em português.
- [Documentação oficial](https://openldap.org/doc/), em inglês.
- [LDAP for Rocket Scientists](https://www.zytrax.com/books/ldap/), em inglês.
