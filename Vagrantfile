# frozen_string_literal: true

base_dn = 'dc=local,dc=br'
ldap_server = 'master.local.br'
admin_pass = '123456'
sync_password = '654321'
admin_dn = "cn=Manager,#{base_dn}"

Vagrant.configure('2') do |config|
  config.vm.box = 'roboxes/centos7'

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  config.vm.box_check_update = true

  config.vm.provider 'virtualbox' do |vb|
    vb.memory = 1024
    vb.name = 'ldap-base'
    vb.linked_clone = true
    # to avoid warnings from Virtualbox
    vb.customize ['modifyvm', :id, '--vrde', 'off']
    vb.customize ['modifyvm', :id, '--graphicscontroller', 'vmsvga']
  end

  config.vm.provision 'shell', inline: 'yum makecache fast && yum upgrade -y && yum install python3 tree mailx -y'

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'network.yaml'
    ansible.config_file = 'ansible.cfg'
  end

  config.vm.define 'vinfra' do |i|
    i.vm.provider 'virtualbox' do |vb|
      vb.name = 'vinfra'
    end
    i.vm.hostname = 'vinfra.local.br'
    i.vm.network 'private_network', ip: '192.168.56.85'

    i.vm.provision :ansible do |ansible|
      ansible.playbook = 'vinfra/main.yaml'
      ansible.config_file = 'ansible.cfg'
    end
  end

  # enable debugging Ansible configuring tasks
  # ENV['ANSIBLE_VERBOSITY'] = '3'

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'openldap.yaml'
    ansible.config_file = 'ansible.cfg'
    ansible.host_vars = {
      'master' => {
        'admin_pass' => admin_pass,
        'admin_dn' => admin_dn,
        'base_dn' => base_dn
      }
    }
  end

  config.vm.define 'master' do |m|
    m.vm.provider 'virtualbox' do |vb|
      vb.name = 'master'
    end
    m.vm.hostname = ldap_server
    m.vm.network 'private_network', ip: '192.168.56.80'
    m.vm.network 'forwarded_port', guest: 389, host: 3389, host_ip: '127.0.0.1', id: 'ldap'

    m.vm.provision :ansible do |ansible|
      ansible.playbook = 'master/main.yaml'
      ansible.config_file = 'ansible.cfg'
      ansible.host_vars = {
        'master' => {
          'admin_pass' => admin_pass,
          'sync_password' => sync_password,
          'admin_dn' => admin_dn,
          'base_dn' => base_dn
        }
      }
    end
  end

  config.vm.define 'client' do |c|
    c.vm.provider 'virtualbox' do |vb|
      vb.name = 'client'
    end
    c.vm.hostname = 'client.local.br'
    c.vm.network 'private_network', ip: '192.168.56.82'

    c.vm.provision :ansible do |ansible|
      ansible.playbook = 'client/main.yaml'
      ansible.config_file = 'ansible.cfg'
      ansible.host_vars = {
        'client' => {
          'base_dn' => base_dn,
          'ldap_server' => ldap_server
        }
      }
    end
  end
end

# -*- mode: ruby -*-
# vi: set ft=ruby :
